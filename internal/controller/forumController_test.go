package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Van-programan/Forum_GO/internal/controller/middleware"
	requests "github.com/Van-programan/Forum_GO/internal/controller/request"
	"github.com/Van-programan/Forum_GO/internal/entity"
	"github.com/Van-programan/Forum_GO/internal/usecase"
	"github.com/Van-programan/Forum_GO/internal/ws"
	mocksf "github.com/Van-programan/Forum_GO/mocks/forum"
	mocks "github.com/Van-programan/Forum_GO/mocks/forum/usecase"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	ContextUserIDKey = "user_id"
	ContextRoleKey   = "role"
)

func TestCategoryHandler_Create_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockUsecase := mocks.NewCategoryUsecase(t)
	logger := zerolog.Nop()
	handler := &CategoryHandler{
		Usecase: mockUsecase,
		Log:     &logger,
	}
	router.POST("/categories", handler.Create)

	reqBody := entity.Category{Title: "New Category", Description: "Desc"}
	expectedID := int64(1)

	mockUsecase.On("Create", mock.Anything, reqBody).Return(expectedID, nil).Once()

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/categories", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	var respBody map[string]int64
	err := json.Unmarshal(rr.Body.Bytes(), &respBody)
	assert.NoError(t, err)
	assert.Equal(t, expectedID, respBody["id"])
	mockUsecase.AssertExpectations(t)
}

func TestCategoryHandler_Create_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockUsecase := mocks.NewCategoryUsecase(t)
	logger := zerolog.Nop()
	handler := &CategoryHandler{
		Usecase: mockUsecase,
		Log:     &logger,
	}
	router.POST("/categories", handler.Create)

	req, _ := http.NewRequest(http.MethodPost, "/categories", bytes.NewBufferString("{invalid_json"))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	mockUsecase.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestCategoryHandler_Create_UsecaseError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockUsecase := mocks.NewCategoryUsecase(t)
	logger := zerolog.Nop()
	handler := &CategoryHandler{
		Usecase: mockUsecase,
		Log:     &logger,
	}
	router.POST("/categories", handler.Create)

	reqBody := entity.Category{Title: "New Category", Description: "Desc"}
	usecaseError := errors.New("usecase create error")

	mockUsecase.On("Create", mock.Anything, reqBody).Return(int64(0), usecaseError).Once()

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/categories", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	var respBody map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &respBody)
	assert.NoError(t, err)
	assert.Equal(t, usecaseError.Error(), respBody["error"])
	mockUsecase.AssertExpectations(t)
}

func TestCategoryHandler_GetByID_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockUsecase := mocks.NewCategoryUsecase(t)
	logger := zerolog.Nop()
	handler := &CategoryHandler{
		Usecase: mockUsecase,
		Log:     &logger,
	}
	categoryID := int64(1)
	router.GET("/categories/:id", handler.GetByID)

	expectedCategory := &entity.Category{ID: categoryID, Title: "Test", Description: "Test Desc", CreatedAt: time.Now()}
	mockUsecase.On("GetByID", mock.Anything, categoryID).Return(expectedCategory, nil).Once()

	req, _ := http.NewRequest(http.MethodGet, "/categories/"+strconv.FormatInt(categoryID, 10), nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var respBody map[string]entity.Category
	err := json.Unmarshal(rr.Body.Bytes(), &respBody)
	assert.NoError(t, err)
	assert.Equal(t, expectedCategory.ID, respBody["category"].ID)
	assert.Equal(t, expectedCategory.Title, respBody["category"].Title)
	mockUsecase.AssertExpectations(t)
}

func TestCategoryHandler_GetByID_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockUsecase := mocks.NewCategoryUsecase(t)
	logger := zerolog.Nop()
	handler := &CategoryHandler{
		Usecase: mockUsecase,
		Log:     &logger,
	}
	router.GET("/categories/:id", handler.GetByID)

	req, _ := http.NewRequest(http.MethodGet, "/categories/invalid", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	mockUsecase.AssertNotCalled(t, "GetByID", mock.Anything, mock.Anything)
}

func TestCategoryHandler_GetByID_UsecaseError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockUsecase := mocks.NewCategoryUsecase(t)
	logger := zerolog.Nop()
	handler := &CategoryHandler{
		Usecase: mockUsecase,
		Log:     &logger,
	}
	categoryID := int64(1)
	router.GET("/categories/:id", handler.GetByID)

	usecaseError := errors.New("usecase get by id error")
	mockUsecase.On("GetByID", mock.Anything, categoryID).Return(nil, usecaseError).Once()

	req, _ := http.NewRequest(http.MethodGet, "/categories/"+strconv.FormatInt(categoryID, 10), nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	var respBody map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &respBody)
	assert.NoError(t, err)
	assert.Equal(t, "failed to get category", respBody["error"])
	mockUsecase.AssertExpectations(t)
}

func TestCategoryHandler_GetAll_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockUsecase := mocks.NewCategoryUsecase(t)
	logger := zerolog.Nop()
	handler := &CategoryHandler{
		Usecase: mockUsecase,
		Log:     &logger,
	}
	router.GET("/categories", handler.GetAll)

	expectedCategories := []entity.Category{
		{ID: 1, Title: "Cat1", Description: "D1", CreatedAt: time.Now()},
		{ID: 2, Title: "Cat2", Description: "D2", CreatedAt: time.Now()},
	}
	mockUsecase.On("GetAll", mock.Anything).Return(expectedCategories, nil).Once()

	req, _ := http.NewRequest(http.MethodGet, "/categories", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var respBody map[string][]entity.Category
	err := json.Unmarshal(rr.Body.Bytes(), &respBody)
	assert.NoError(t, err)
	assert.Len(t, respBody["categories"], 2)
	assert.Equal(t, expectedCategories[0].ID, respBody["categories"][0].ID)
	mockUsecase.AssertExpectations(t)
}

func TestCategoryHandler_GetAll_UsecaseError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockUsecase := mocks.NewCategoryUsecase(t)
	logger := zerolog.Nop()
	handler := &CategoryHandler{
		Usecase: mockUsecase,
		Log:     &logger,
	}
	router.GET("/categories", handler.GetAll)

	usecaseError := errors.New("usecase get all error")
	mockUsecase.On("GetAll", mock.Anything).Return(nil, usecaseError).Once()

	req, _ := http.NewRequest(http.MethodGet, "/categories", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	var respBody map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &respBody)
	assert.NoError(t, err)
	assert.Equal(t, usecaseError.Error(), respBody["error"])
	mockUsecase.AssertExpectations(t)
}

func TestCategoryHandler_Delete_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockUsecase := mocks.NewCategoryUsecase(t)
	logger := zerolog.Nop()
	handler := &CategoryHandler{
		Usecase: mockUsecase,
		Log:     &logger,
	}
	categoryID := int64(1)
	router.DELETE("/categories/:id", handler.Delete)

	mockUsecase.On("Delete", mock.Anything, categoryID).Return(nil).Once()

	req, _ := http.NewRequest(http.MethodDelete, "/categories/"+strconv.FormatInt(categoryID, 10), nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockUsecase.AssertExpectations(t)
}

func TestCategoryHandler_Delete_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockUsecase := mocks.NewCategoryUsecase(t)
	logger := zerolog.Nop()
	handler := &CategoryHandler{
		Usecase: mockUsecase,
		Log:     &logger,
	}
	router.DELETE("/categories/:id", handler.Delete)

	req, _ := http.NewRequest(http.MethodDelete, "/categories/invalid", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	mockUsecase.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything)
}

func TestCategoryHandler_Delete_UsecaseError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockUsecase := mocks.NewCategoryUsecase(t)
	logger := zerolog.Nop()
	handler := &CategoryHandler{
		Usecase: mockUsecase,
		Log:     &logger,
	}
	categoryID := int64(1)
	router.DELETE("/categories/:id", handler.Delete)

	usecaseError := errors.New("usecase delete error")
	mockUsecase.On("Delete", mock.Anything, categoryID).Return(usecaseError).Once()

	req, _ := http.NewRequest(http.MethodDelete, "/categories/"+strconv.FormatInt(categoryID, 10), nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	var respBody map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &respBody)
	assert.NoError(t, err)
	assert.Equal(t, "failed to delete category", respBody["error"])
	mockUsecase.AssertExpectations(t)
}

func TestCategoryHandler_Update_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockUsecase := mocks.NewCategoryUsecase(t)
	logger := zerolog.Nop()
	handler := &CategoryHandler{
		Usecase: mockUsecase,
		Log:     &logger,
	}
	categoryID := int64(1)
	router.PUT("/categories/:id", handler.Update)

	reqBody := requests.UpdateRequestCategory{Title: "updated title", Description: "updated desc"}
	mockUsecase.On("Update", mock.Anything, categoryID, reqBody.Title, reqBody.Description).Return(nil).Once()

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPut, "/categories/"+strconv.FormatInt(categoryID, 10), bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockUsecase.AssertExpectations(t)
}

func TestCategoryHandler_Update_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockUsecase := mocks.NewCategoryUsecase(t)
	logger := zerolog.Nop()
	handler := &CategoryHandler{
		Usecase: mockUsecase,
		Log:     &logger,
	}
	router.PUT("/categories/:id", handler.Update)

	reqBody := requests.UpdateRequestCategory{Title: "updated title", Description: "updated desc"}
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPut, "/categories/invalid", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	mockUsecase.AssertNotCalled(t, "Update", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestCategoryHandler_Update_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockUsecase := mocks.NewCategoryUsecase(t)
	logger := zerolog.Nop()
	handler := &CategoryHandler{
		Usecase: mockUsecase,
		Log:     &logger,
	}
	categoryID := int64(1)
	router.PUT("/categories/:id", handler.Update)

	req, _ := http.NewRequest(http.MethodPut, "/categories/"+strconv.FormatInt(categoryID, 10), bytes.NewBufferString("{invalid_json"))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	mockUsecase.AssertNotCalled(t, "Update", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestCategoryHandler_Update_UsecaseError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockUsecase := mocks.NewCategoryUsecase(t)
	logger := zerolog.Nop()
	handler := &CategoryHandler{
		Usecase: mockUsecase,
		Log:     &logger,
	}
	categoryID := int64(1)
	router.PUT("/categories/:id", handler.Update)

	reqBody := requests.UpdateRequestCategory{Title: "updated title", Description: "updated desc"}
	usecaseError := errors.New("usecase update error")
	mockUsecase.On("Update", mock.Anything, categoryID, reqBody.Title, reqBody.Description).Return(usecaseError).Once()

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPut, "/categories/"+strconv.FormatInt(categoryID, 10), bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	var respBody map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &respBody)
	assert.NoError(t, err)
	assert.Equal(t, "failed to update category", respBody["error"])
	mockUsecase.AssertExpectations(t)
}

func setupTestServerForChatOnlyUpgrade(t *testing.T, handler *ChatHandler) (*httptest.Server, string) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/ws", handler.ServeWs)
	server := httptest.NewServer(router)
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"
	t.Cleanup(server.Close)
	return server, wsURL
}

func TestChatHandler_ServeWs_UpgradesConnection_UnauthorizedPath(t *testing.T) {
	logger := zerolog.Nop()
	emptyMockChatUsecase := new(mocks.ChatUsecase)
	emptyMockUserClient := new(mocksf.UserClient)
	dummyHub := ws.NewHub(&logger)

	chatHandler := NewChatHandler(dummyHub, emptyMockChatUsecase, emptyMockUserClient, &logger)
	_, wsURL := setupTestServerForChatOnlyUpgrade(t, chatHandler)

	dialer := websocket.Dialer{HandshakeTimeout: 2 * time.Second}
	conn, resp, errDial := dialer.Dial(wsURL, http.Header{"Origin": []string{"http://localhost:5173"}})

	if conn != nil {
		defer conn.Close()
	}

	assert.NoError(t, errDial, "Upgrade to WebSocket should succeed")
	if errDial != nil && resp != nil {
		t.Logf("Response status: %s", resp.Status)
	}
	assert.NotNil(t, conn, "Connection should be established")

	if resp != nil && errDial == nil {
		assert.Equal(t, http.StatusSwitchingProtocols, resp.StatusCode, "HTTP status should be 101 Switching Protocols")
	}
}

func TestChatHandler_ServeWs_UpgradesConnection_AuthorizedPath(t *testing.T) {
	logger := zerolog.Nop()
	emptyMockChatUsecase := new(mocks.ChatUsecase)
	mockUserClientActual := new(mocksf.UserClient)
	dummyHub := ws.NewHub(&logger)

	expectedUserID := int64(123)
	expectedUsername := "testuser"

	mockUserClientActual.On("GetUsername", mock.Anything, expectedUserID).Return(expectedUsername, nil).Maybe()

	chatHandler := NewChatHandler(dummyHub, emptyMockChatUsecase, mockUserClientActual, &logger)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(middleware.ContextUserIDKey, expectedUserID)
		c.Next()
	})
	router.GET("/ws", chatHandler.ServeWs)
	server := httptest.NewServer(router)
	defer server.Close()
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"

	dialer := websocket.Dialer{HandshakeTimeout: 2 * time.Second}
	conn, resp, errDial := dialer.Dial(wsURL, http.Header{"Origin": []string{"http://localhost:5173"}})

	if conn != nil {
		defer conn.Close()
	}

	assert.NoError(t, errDial, "Upgrade to WebSocket should succeed for authorized path")
	if errDial != nil && resp != nil {
		t.Logf("Response status: %s", resp.Status)
	}
	assert.NotNil(t, conn, "Connection should be established for authorized path")

	if resp != nil && errDial == nil {
		assert.Equal(t, http.StatusSwitchingProtocols, resp.StatusCode, "HTTP status should be 101 Switching Protocols")
	}
	// mockUserClientActual.AssertExpectations(t) // Опционально
}

func TestChatHandler_ServeWs_UpgradeFail_BadOrigin_NoHubLogic(t *testing.T) {
	logger := zerolog.Nop()
	emptyMockChatUsecase := new(mocks.ChatUsecase)
	emptyMockUserClient := new(mocksf.UserClient)
	dummyHub := ws.NewHub(&logger)

	chatHandler := NewChatHandler(dummyHub, emptyMockChatUsecase, emptyMockUserClient, &logger)
	_, wsURL := setupTestServerForChatOnlyUpgrade(t, chatHandler)

	dialer := websocket.Dialer{HandshakeTimeout: 1 * time.Second}
	conn, resp, err := dialer.Dial(wsURL, http.Header{"Origin": []string{"http://bad-origin.com"}})
	if conn != nil {
		defer conn.Close()
	}

	assert.Error(t, err, "Expected error when dialing with bad origin")
	if assert.NotNil(t, resp, "Response should not be nil on handshake failure") {
		assert.Equal(t, http.StatusForbidden, resp.StatusCode, "Response status should be 403 Forbidden for bad origin")
	}
	if err != nil {
		assert.Contains(t, err.Error(), "bad handshake", "Error message should indicate bad handshake")
	}
}

func TestPostHandler_Create_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockUsecase := mocks.NewPostUsecase(t)
	handler := &PostHandler{
		Usecase: mockUsecase,
	}
	topicID := int64(1)
	userID := int64(10)

	router.POST("/topics/:id/posts", func(c *gin.Context) {
		c.Set(ContextUserIDKey, userID)
		c.Set(ContextRoleKey, "user")
		handler.Create(c)
	})

	reqBody := entity.Post{Content: "post content"}
	expectedPostID := int64(5)

	expectedEntityPost := entity.Post{TopicID: topicID, AuthorID: &userID, Content: reqBody.Content}
	mockUsecase.On("Create", mock.Anything, expectedEntityPost).Return(expectedPostID, nil).Once()

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/topics/"+strconv.FormatInt(topicID, 10)+"/posts", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var respBody map[string]int64
	err := json.Unmarshal(rr.Body.Bytes(), &respBody)
	assert.NoError(t, err)
	assert.Equal(t, expectedPostID, respBody["id"])
	mockUsecase.AssertExpectations(t)
}

func TestPostHandler_Create_NoUserIDInContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockUsecase := mocks.NewPostUsecase(t)
	handler := &PostHandler{
		Usecase: mockUsecase,
	}
	topicID := int64(1)
	router.POST("/topics/:id/posts", handler.Create)

	reqBody := entity.Post{Content: "Test Content"}
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/topics/"+strconv.FormatInt(topicID, 10)+"/posts", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusForbidden, rr.Code)
	mockUsecase.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestPostHandler_Create_InvalidTopicID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockUsecase := mocks.NewPostUsecase(t)
	handler := &PostHandler{
		Usecase: mockUsecase,
	}
	userID := int64(10)
	router.POST("/topics/:id/posts", func(c *gin.Context) {
		c.Set(ContextUserIDKey, userID)
		c.Set(ContextRoleKey, "user")
		handler.Create(c)
	})

	reqBody := entity.Post{Content: "Test Content"}
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/topics/invalid/posts", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	var respBody map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &respBody)
	assert.NoError(t, err)
	assert.Equal(t, "invalid category id", respBody["error"])
	mockUsecase.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestPostHandler_Create_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockUsecase := mocks.NewPostUsecase(t)
	handler := &PostHandler{
		Usecase: mockUsecase,
	}
	topicID := int64(1)
	userID := int64(10)
	router.POST("/topics/:id/posts", func(c *gin.Context) {
		c.Set(ContextUserIDKey, userID)
		c.Set(ContextRoleKey, "user")
		handler.Create(c)
	})

	req, _ := http.NewRequest(http.MethodPost, "/topics/"+strconv.FormatInt(topicID, 10)+"/posts", bytes.NewBufferString("{invalid_json"))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	mockUsecase.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestPostHandler_Create_TopicNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockUsecase := mocks.NewPostUsecase(t)
	handler := &PostHandler{
		Usecase: mockUsecase,
	}
	topicID := int64(1)
	userID := int64(10)
	router.POST("/topics/:id/posts", func(c *gin.Context) {
		c.Set(ContextUserIDKey, userID)
		c.Set(ContextRoleKey, "user")
		handler.Create(c)
	})

	reqBody := entity.Post{Content: "Test Content"}
	usecaseError := usecase.ErrTopicNotFound

	expectedEntityPost := entity.Post{TopicID: topicID, AuthorID: &userID, Content: reqBody.Content}
	mockUsecase.On("Create", mock.Anything, expectedEntityPost).Return(int64(0), usecaseError).Once()

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/topics/"+strconv.FormatInt(topicID, 10)+"/posts", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	var respBody map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &respBody)
	assert.NoError(t, err)
	assert.Equal(t, usecaseError.Error(), respBody["error"])
	mockUsecase.AssertExpectations(t)
}

func TestPostHandler_Create_UsecaseError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockUsecase := mocks.NewPostUsecase(t)
	handler := &PostHandler{
		Usecase: mockUsecase,
	}
	topicID := int64(1)
	userID := int64(10)
	router.POST("/topics/:id/posts", func(c *gin.Context) {
		c.Set(ContextUserIDKey, userID)
		c.Set(ContextRoleKey, "user")
		handler.Create(c)
	})

	reqBody := entity.Post{Content: "Test Content"}
	usecaseError := errors.New("some other usecase error")

	expectedEntityPost := entity.Post{TopicID: topicID, AuthorID: &userID, Content: reqBody.Content}
	mockUsecase.On("Create", mock.Anything, expectedEntityPost).Return(int64(0), usecaseError).Once()

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/topics/"+strconv.FormatInt(topicID, 10)+"/posts", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	var respBody map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &respBody)
	assert.NoError(t, err)
	assert.Equal(t, usecaseError.Error(), respBody["error"])
	mockUsecase.AssertExpectations(t)
}

func TestPostHandler_GetByTopic_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockUsecase := mocks.NewPostUsecase(t)
	handler := &PostHandler{
		Usecase: mockUsecase,
	}
	topicID := int64(1)
	router.GET("/topics/:id/posts", handler.GetByTopic)

	expectedPosts := []entity.Post{
		{ID: 1, TopicID: topicID, Content: "Post 1", Username: "User1"},
		{ID: 2, TopicID: topicID, Content: "Post 2", Username: "User2"},
	}
	mockUsecase.On("GetByTopic", mock.Anything, topicID).Return(expectedPosts, nil).Once()

	req, _ := http.NewRequest(http.MethodGet, "/topics/"+strconv.FormatInt(topicID, 10)+"/posts", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var respBody map[string][]entity.Post
	err := json.Unmarshal(rr.Body.Bytes(), &respBody)
	assert.NoError(t, err)
	assert.Len(t, respBody["posts"], 2)
	assert.Equal(t, expectedPosts[0].Content, respBody["posts"][0].Content)
	mockUsecase.AssertExpectations(t)
}

func TestPostHandler_GetByTopic_InvalidTopicID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockUsecase := mocks.NewPostUsecase(t)
	handler := &PostHandler{
		Usecase: mockUsecase,
	}
	router.GET("/topics/:id/posts", handler.GetByTopic)

	req, _ := http.NewRequest(http.MethodGet, "/topics/invalid/posts", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	var respBody map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &respBody)
	assert.NoError(t, err)
	assert.Equal(t, "invalid topic id", respBody["error"])
	mockUsecase.AssertNotCalled(t, "GetByTopic", mock.Anything, mock.Anything)
}

func TestPostHandler_GetByTopic_TopicNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockUsecase := mocks.NewPostUsecase(t)
	handler := &PostHandler{
		Usecase: mockUsecase,
	}
	topicID := int64(1)
	router.GET("/topics/:id/posts", handler.GetByTopic)

	usecaseError := usecase.ErrTopicNotFound
	mockUsecase.On("GetByTopic", mock.Anything, topicID).Return(nil, usecaseError).Once()

	req, _ := http.NewRequest(http.MethodGet, "/topics/"+strconv.FormatInt(topicID, 10)+"/posts", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	var respBody map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &respBody)
	assert.NoError(t, err)
	assert.Equal(t, usecaseError.Error(), respBody["error"])
	mockUsecase.AssertExpectations(t)
}

func TestPostHandler_GetByTopic_UsecaseError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockUsecase := mocks.NewPostUsecase(t)
	handler := &PostHandler{
		Usecase: mockUsecase,
	}
	topicID := int64(1)
	router.GET("/topics/:id/posts", handler.GetByTopic)

	usecaseError := errors.New("some other get by topic error")
	mockUsecase.On("GetByTopic", mock.Anything, topicID).Return(nil, usecaseError).Once()

	req, _ := http.NewRequest(http.MethodGet, "/topics/"+strconv.FormatInt(topicID, 10)+"/posts", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	var respBody map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &respBody)
	assert.NoError(t, err)
	assert.Equal(t, usecaseError.Error(), respBody["error"])
	mockUsecase.AssertExpectations(t)
}

func TestPostHandler_Update_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockUsecase := mocks.NewPostUsecase(t)
	handler := &PostHandler{
		Usecase: mockUsecase,
	}
	postID := int64(1)
	userID := int64(10)
	userRole := "user"

	router.PUT("/posts/:id", func(c *gin.Context) {
		c.Set(ContextUserIDKey, userID)
		c.Set(ContextRoleKey, userRole)
		handler.Update(c)
	})

	reqBody := requests.UpdateRequestPost{Content: "updated content"}
	mockUsecase.On("Update", mock.Anything, postID, userID, userRole, reqBody.Content).Return(nil).Once()

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPut, "/posts/"+strconv.FormatInt(postID, 10), bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var respBody map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &respBody)
	assert.NoError(t, err)
	assert.Equal(t, "post updated", respBody["message"])
	mockUsecase.AssertExpectations(t)
}

func TestPostHandler_Update_NoUserIDInContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockUsecase := mocks.NewPostUsecase(t)
	handler := &PostHandler{
		Usecase: mockUsecase,
	}
	postID := int64(1)
	router.PUT("/posts/:id", handler.Update)

	reqBody := requests.UpdateRequestPost{Content: "updated content"}
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPut, "/posts/"+strconv.FormatInt(postID, 10), bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusForbidden, rr.Code)
	mockUsecase.AssertNotCalled(t, "Update", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestPostHandler_Update_InvalidPostID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockUsecase := mocks.NewPostUsecase(t)
	handler := &PostHandler{
		Usecase: mockUsecase,
	}
	userID := int64(10)
	userRole := "user"
	router.PUT("/posts/:id", func(c *gin.Context) {
		c.Set(ContextUserIDKey, userID)
		c.Set(ContextRoleKey, userRole)
		handler.Update(c)
	})

	reqBody := requests.UpdateRequestPost{Content: "updated content"}
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPut, "/posts/invalid", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	var respBody map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &respBody)
	assert.NoError(t, err)
	assert.Equal(t, "invalid post id", respBody["error"])
	mockUsecase.AssertNotCalled(t, "Update", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestPostHandler_Update_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockUsecase := mocks.NewPostUsecase(t)
	handler := &PostHandler{
		Usecase: mockUsecase,
	}
	postID := int64(1)
	userID := int64(10)
	userRole := "user"
	router.PUT("/posts/:id", func(c *gin.Context) {
		c.Set(ContextUserIDKey, userID)
		c.Set(ContextRoleKey, userRole)
		handler.Update(c)
	})

	req, _ := http.NewRequest(http.MethodPut, "/posts/"+strconv.FormatInt(postID, 10), bytes.NewBufferString("{invalid_json"))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	mockUsecase.AssertNotCalled(t, "Update", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestPostHandler_Update_Forbidden(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockUsecase := mocks.NewPostUsecase(t)
	handler := &PostHandler{
		Usecase: mockUsecase,
	}
	postID := int64(1)
	userID := int64(10)
	userRole := "user"
	router.PUT("/posts/:id", func(c *gin.Context) {
		c.Set(ContextUserIDKey, userID)
		c.Set(ContextRoleKey, userRole)
		handler.Update(c)
	})

	reqBody := requests.UpdateRequestPost{Content: "updated content"}
	usecaseError := usecase.ErrForbidden
	mockUsecase.On("Update", mock.Anything, postID, userID, userRole, reqBody.Content).Return(usecaseError).Once()

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPut, "/posts/"+strconv.FormatInt(postID, 10), bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusForbidden, rr.Code)
	var respBody map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &respBody)
	assert.NoError(t, err)
	assert.Equal(t, "insufficient permissions", respBody["error"])
	mockUsecase.AssertExpectations(t)
}

func TestPostHandler_Update_PostNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockUsecase := mocks.NewPostUsecase(t)
	handler := &PostHandler{
		Usecase: mockUsecase,
	}
	postID := int64(1)
	userID := int64(10)
	userRole := "user"
	router.PUT("/posts/:id", func(c *gin.Context) {
		c.Set(ContextUserIDKey, userID)
		c.Set(ContextRoleKey, userRole)
		handler.Update(c)
	})

	reqBody := requests.UpdateRequestPost{Content: "updated content"}
	usecaseError := usecase.ErrPostNotFound
	mockUsecase.On("Update", mock.Anything, postID, userID, userRole, reqBody.Content).Return(usecaseError).Once()

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPut, "/posts/"+strconv.FormatInt(postID, 10), bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	var respBody map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &respBody)
	assert.NoError(t, err)
	assert.Equal(t, "post not found", respBody["error"])
	mockUsecase.AssertExpectations(t)
}

func TestPostHandler_Update_UsecaseError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockUsecase := mocks.NewPostUsecase(t)
	handler := &PostHandler{
		Usecase: mockUsecase,
	}
	postID := int64(1)
	userID := int64(10)
	userRole := "user"
	router.PUT("/posts/:id", func(c *gin.Context) {
		c.Set(ContextUserIDKey, userID)
		c.Set(ContextRoleKey, userRole)
		handler.Update(c)
	})

	reqBody := requests.UpdateRequestPost{Content: "updated content"}
	usecaseError := errors.New("some other update error")
	mockUsecase.On("Update", mock.Anything, postID, userID, userRole, reqBody.Content).Return(usecaseError).Once()

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPut, "/posts/"+strconv.FormatInt(postID, 10), bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	var respBody map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &respBody)
	assert.NoError(t, err)
	assert.Equal(t, "internal server error", respBody["error"])
	mockUsecase.AssertExpectations(t)
}

func TestPostHandler_Delete_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockUsecase := mocks.NewPostUsecase(t)
	handler := &PostHandler{
		Usecase: mockUsecase,
	}
	postID := int64(1)
	userID := int64(10)
	userRole := "user"

	router.DELETE("/posts/:id", func(c *gin.Context) {
		c.Set(ContextUserIDKey, userID)
		c.Set(ContextRoleKey, userRole)
		handler.Delete(c)
	})

	mockUsecase.On("Delete", mock.Anything, postID, userID, userRole).Return(nil).Once()

	req, _ := http.NewRequest(http.MethodDelete, "/posts/"+strconv.FormatInt(postID, 10), nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var respBody map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &respBody)
	assert.NoError(t, err)
	assert.Equal(t, "post deleted", respBody["message"])
	mockUsecase.AssertExpectations(t)
}

func TestPostHandler_Delete_NoUserIDInContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockUsecase := mocks.NewPostUsecase(t)
	handler := &PostHandler{
		Usecase: mockUsecase,
	}
	postID := int64(1)
	router.DELETE("/posts/:id", handler.Delete)

	req, _ := http.NewRequest(http.MethodDelete, "/posts/"+strconv.FormatInt(postID, 10), nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusForbidden, rr.Code)
	mockUsecase.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestPostHandler_Delete_InvalidPostID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockUsecase := mocks.NewPostUsecase(t)
	handler := &PostHandler{
		Usecase: mockUsecase,
	}
	userID := int64(10)
	userRole := "user"
	router.DELETE("/posts/:id", func(c *gin.Context) {
		c.Set(ContextUserIDKey, userID)
		c.Set(ContextRoleKey, userRole)
		handler.Delete(c)
	})

	req, _ := http.NewRequest(http.MethodDelete, "/posts/invalid", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	var respBody map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &respBody)
	assert.NoError(t, err)
	assert.Equal(t, "invalid post id", respBody["error"])
	mockUsecase.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestPostHandler_Delete_Forbidden(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockUsecase := mocks.NewPostUsecase(t)
	handler := &PostHandler{
		Usecase: mockUsecase,
	}
	postID := int64(1)
	userID := int64(10)
	userRole := "user"
	router.DELETE("/posts/:id", func(c *gin.Context) {
		c.Set(ContextUserIDKey, userID)
		c.Set(ContextRoleKey, userRole)
		handler.Delete(c)
	})

	usecaseError := usecase.ErrForbidden
	mockUsecase.On("Delete", mock.Anything, postID, userID, userRole).Return(usecaseError).Once()

	req, _ := http.NewRequest(http.MethodDelete, "/posts/"+strconv.FormatInt(postID, 10), nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusForbidden, rr.Code)
	var respBody map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &respBody)
	assert.NoError(t, err)
	assert.Equal(t, "insufficient permissions", respBody["error"])
	mockUsecase.AssertExpectations(t)
}

func TestPostHandler_Delete_PostNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockUsecase := mocks.NewPostUsecase(t)
	handler := &PostHandler{
		Usecase: mockUsecase,
	}
	postID := int64(1)
	userID := int64(10)
	userRole := "user"
	router.DELETE("/posts/:id", func(c *gin.Context) {
		c.Set(ContextUserIDKey, userID)
		c.Set(ContextRoleKey, userRole)
		handler.Delete(c)
	})

	usecaseError := usecase.ErrPostNotFound
	mockUsecase.On("Delete", mock.Anything, postID, userID, userRole).Return(usecaseError).Once()

	req, _ := http.NewRequest(http.MethodDelete, "/posts/"+strconv.FormatInt(postID, 10), nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	var respBody map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &respBody)
	assert.NoError(t, err)
	assert.Equal(t, "post not found", respBody["error"])
	mockUsecase.AssertExpectations(t)
}

func TestPostHandler_Delete_UsecaseError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockUsecase := mocks.NewPostUsecase(t)
	handler := &PostHandler{
		Usecase: mockUsecase,
	}
	postID := int64(1)
	userID := int64(10)
	userRole := "user"
	router.DELETE("/posts/:id", func(c *gin.Context) {
		c.Set(ContextUserIDKey, userID)
		c.Set(ContextRoleKey, userRole)
		handler.Delete(c)
	})

	usecaseError := errors.New("some other delete error")
	mockUsecase.On("Delete", mock.Anything, postID, userID, userRole).Return(usecaseError).Once()

	req, _ := http.NewRequest(http.MethodDelete, "/posts/"+strconv.FormatInt(postID, 10), nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	var respBody map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &respBody)
	assert.NoError(t, err)
	assert.Equal(t, "internal server error", respBody["error"])
	mockUsecase.AssertExpectations(t)
}

func TestTopicHandler_Create_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockUsecase := mocks.NewTopicUsecase(t)
	logger := zerolog.Nop()
	handler := &TopicHandler{
		Usecase: mockUsecase,
		Log:     &logger,
	}
	categoryID := int64(1)
	userID := int64(10)

	router.POST("/categories/:id/topics", func(c *gin.Context) {
		c.Set(ContextUserIDKey, userID)
		c.Set(ContextRoleKey, "user")
		handler.Create(c)
	})

	reqBody := entity.Topic{Title: "new topic"}
	expectedTopicID := int64(5)

	expectedEntityTopic := entity.Topic{CategoryID: categoryID, AuthorID: &userID, Title: reqBody.Title}
	mockUsecase.On("Create", mock.Anything, expectedEntityTopic).Return(expectedTopicID, nil).Once()

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/categories/"+strconv.FormatInt(categoryID, 10)+"/topics", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var respBody map[string]int64
	err := json.Unmarshal(rr.Body.Bytes(), &respBody)
	assert.NoError(t, err)
	assert.Equal(t, expectedTopicID, respBody["id"])
	mockUsecase.AssertExpectations(t)
}

func TestTopicHandler_Create_NoUserIDInContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockUsecase := mocks.NewTopicUsecase(t)
	logger := zerolog.Nop()
	handler := &TopicHandler{
		Usecase: mockUsecase,
		Log:     &logger,
	}
	categoryID := int64(1)
	router.POST("/categories/:id/topics", handler.Create)

	reqBody := entity.Topic{Title: "new topic"}
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/categories/"+strconv.FormatInt(categoryID, 10)+"/topics", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusForbidden, rr.Code)
	mockUsecase.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestTopicHandler_Create_InvalidCategoryID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockUsecase := mocks.NewTopicUsecase(t)
	logger := zerolog.Nop()
	handler := &TopicHandler{
		Usecase: mockUsecase,
		Log:     &logger,
	}
	userID := int64(10)
	router.POST("/categories/:id/topics", func(c *gin.Context) {
		c.Set(ContextUserIDKey, userID)
		c.Set(ContextRoleKey, "user")
		handler.Create(c)
	})

	reqBody := entity.Topic{Title: "new topic"}
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/categories/invalid/topics", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	var respBody map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &respBody)
	assert.NoError(t, err)
	assert.Equal(t, "invalid category id", respBody["error"])
	mockUsecase.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestTopicHandler_Create_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockUsecase := mocks.NewTopicUsecase(t)
	logger := zerolog.Nop()
	handler := &TopicHandler{
		Usecase: mockUsecase,
		Log:     &logger,
	}
	categoryID := int64(1)
	userID := int64(10)
	router.POST("/categories/:id/topics", func(c *gin.Context) {
		c.Set(ContextUserIDKey, userID)
		c.Set(ContextRoleKey, "user")
		handler.Create(c)
	})

	req, _ := http.NewRequest(http.MethodPost, "/categories/"+strconv.FormatInt(categoryID, 10)+"/topics", bytes.NewBufferString("{invalid_json"))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	mockUsecase.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestTopicHandler_Create_CategoryNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockUsecase := mocks.NewTopicUsecase(t)
	logger := zerolog.Nop()
	handler := &TopicHandler{
		Usecase: mockUsecase,
		Log:     &logger,
	}
	categoryID := int64(1)
	userID := int64(10)
	router.POST("/categories/:id/topics", func(c *gin.Context) {
		c.Set(ContextUserIDKey, userID)
		c.Set(ContextRoleKey, "user")
		handler.Create(c)
	})

	reqBody := entity.Topic{Title: "new topic"}
	usecaseError := usecase.ErrCategoryNotFound

	expectedEntityTopic := entity.Topic{CategoryID: categoryID, AuthorID: &userID, Title: reqBody.Title}
	mockUsecase.On("Create", mock.Anything, expectedEntityTopic).Return(int64(0), usecaseError).Once()

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/categories/"+strconv.FormatInt(categoryID, 10)+"/topics", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	var respBody map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &respBody)
	assert.NoError(t, err)
	assert.Equal(t, usecaseError.Error(), respBody["error"])
	mockUsecase.AssertExpectations(t)
}

func TestTopicHandler_Create_UsecaseError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockUsecase := mocks.NewTopicUsecase(t)
	logger := zerolog.Nop()
	handler := &TopicHandler{
		Usecase: mockUsecase,
		Log:     &logger,
	}
	categoryID := int64(1)
	userID := int64(10)
	router.POST("/categories/:id/topics", func(c *gin.Context) {
		c.Set(ContextUserIDKey, userID)
		c.Set(ContextRoleKey, "user")
		handler.Create(c)
	})

	reqBody := entity.Topic{Title: "new topic"}
	usecaseError := errors.New("some other create error")

	expectedEntityTopic := entity.Topic{CategoryID: categoryID, AuthorID: &userID, Title: reqBody.Title}
	mockUsecase.On("Create", mock.Anything, expectedEntityTopic).Return(int64(0), usecaseError).Once()

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/categories/"+strconv.FormatInt(categoryID, 10)+"/topics", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	var respBody map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &respBody)
	assert.NoError(t, err)
	assert.Equal(t, usecaseError.Error(), respBody["error"])
	mockUsecase.AssertExpectations(t)
}

func TestTopicHandler_GetByID_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockUsecase := mocks.NewTopicUsecase(t)
	logger := zerolog.Nop()
	handler := &TopicHandler{
		Usecase: mockUsecase,
		Log:     &logger,
	}
	topicID := int64(1)
	router.GET("/topics/:id", handler.GetByID)

	expectedTopic := &entity.Topic{ID: topicID, Title: "Test Topic", Username: "Author"}
	mockUsecase.On("GetByID", mock.Anything, topicID).Return(expectedTopic, nil).Once()

	req, _ := http.NewRequest(http.MethodGet, "/topics/"+strconv.FormatInt(topicID, 10), nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var respBody map[string]entity.Topic
	err := json.Unmarshal(rr.Body.Bytes(), &respBody)
	assert.NoError(t, err)
	assert.Equal(t, expectedTopic.Title, respBody["topic"].Title)
	mockUsecase.AssertExpectations(t)
}

func TestTopicHandler_GetByID_InvalidTopicID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockUsecase := mocks.NewTopicUsecase(t)
	logger := zerolog.Nop()
	handler := &TopicHandler{
		Usecase: mockUsecase,
		Log:     &logger,
	}
	router.GET("/topics/:id", handler.GetByID)

	req, _ := http.NewRequest(http.MethodGet, "/topics/invalid", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	mockUsecase.AssertNotCalled(t, "GetByID", mock.Anything, mock.Anything)
}

func TestTopicHandler_GetByID_UsecaseError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockUsecase := mocks.NewTopicUsecase(t)
	logger := zerolog.Nop()
	handler := &TopicHandler{
		Usecase: mockUsecase,
		Log:     &logger,
	}
	topicID := int64(1)
	router.GET("/topics/:id", handler.GetByID)

	usecaseError := errors.New("usecase get by id error")
	mockUsecase.On("GetByID", mock.Anything, topicID).Return(nil, usecaseError).Once()

	req, _ := http.NewRequest(http.MethodGet, "/topics/"+strconv.FormatInt(topicID, 10), nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	var respBody map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &respBody)
	assert.NoError(t, err)
	assert.Equal(t, "failed to get topic", respBody["error"])
	mockUsecase.AssertExpectations(t)
}

func TestTopicHandler_GetByCategory_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockUsecase := mocks.NewTopicUsecase(t)
	logger := zerolog.Nop()
	handler := &TopicHandler{
		Usecase: mockUsecase,
		Log:     &logger,
	}
	categoryID := int64(1)
	router.GET("/categories/:id/topics", handler.GetByCategory)

	expectedTopics := []entity.Topic{
		{ID: 1, CategoryID: categoryID, Title: "Topic 1", Username: "User1"},
		{ID: 2, CategoryID: categoryID, Title: "Topic 2", Username: "User2"},
	}
	mockUsecase.On("GetByCategory", mock.Anything, categoryID).Return(expectedTopics, nil).Once()

	req, _ := http.NewRequest(http.MethodGet, "/categories/"+strconv.FormatInt(categoryID, 10)+"/topics", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var respBody map[string][]entity.Topic
	err := json.Unmarshal(rr.Body.Bytes(), &respBody)
	assert.NoError(t, err)
	assert.Len(t, respBody["topics"], 2)
	assert.Equal(t, expectedTopics[0].Title, respBody["topics"][0].Title)
	mockUsecase.AssertExpectations(t)
}

func TestTopicHandler_GetByCategory_InvalidCategoryID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockUsecase := mocks.NewTopicUsecase(t)
	logger := zerolog.Nop()
	handler := &TopicHandler{
		Usecase: mockUsecase,
		Log:     &logger,
	}
	router.GET("/categories/:id/topics", handler.GetByCategory)

	req, _ := http.NewRequest(http.MethodGet, "/categories/invalid/topics", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	mockUsecase.AssertNotCalled(t, "GetByCategory", mock.Anything, mock.Anything)
}

func TestTopicHandler_GetByCategory_CategoryNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockUsecase := mocks.NewTopicUsecase(t)
	logger := zerolog.Nop()
	handler := &TopicHandler{
		Usecase: mockUsecase,
		Log:     &logger,
	}
	categoryID := int64(1)
	router.GET("/categories/:id/topics", handler.GetByCategory)

	usecaseError := usecase.ErrCategoryNotFound
	mockUsecase.On("GetByCategory", mock.Anything, categoryID).Return(nil, usecaseError).Once()

	req, _ := http.NewRequest(http.MethodGet, "/categories/"+strconv.FormatInt(categoryID, 10)+"/topics", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	var respBody map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &respBody)
	assert.NoError(t, err)
	assert.Equal(t, usecaseError.Error(), respBody["error"])
	mockUsecase.AssertExpectations(t)
}

func TestTopicHandler_GetByCategory_UsecaseError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockUsecase := mocks.NewTopicUsecase(t)
	logger := zerolog.Nop()
	handler := &TopicHandler{
		Usecase: mockUsecase,
		Log:     &logger,
	}
	categoryID := int64(1)
	router.GET("/categories/:id/topics", handler.GetByCategory)

	usecaseError := errors.New("some other get by category error")
	mockUsecase.On("GetByCategory", mock.Anything, categoryID).Return(nil, usecaseError).Once()

	req, _ := http.NewRequest(http.MethodGet, "/categories/"+strconv.FormatInt(categoryID, 10)+"/topics", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	var respBody map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &respBody)
	assert.NoError(t, err)
	assert.Equal(t, usecaseError.Error(), respBody["error"])
	mockUsecase.AssertExpectations(t)
}

func TestTopicHandler_Update_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockUsecase := mocks.NewTopicUsecase(t)
	logger := zerolog.Nop()
	handler := &TopicHandler{
		Usecase: mockUsecase,
		Log:     &logger,
	}
	topicID := int64(1)
	userID := int64(10)
	userRole := "user"

	router.PUT("/topics/:id", func(c *gin.Context) {
		c.Set(ContextUserIDKey, userID)
		c.Set(ContextRoleKey, userRole)
		handler.Update(c)
	})

	reqBody := requests.UpdateRequestTopic{Title: "Updated Topic Title"}
	mockUsecase.On("Update", mock.Anything, topicID, userID, userRole, reqBody.Title).Return(nil).Once()

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPut, "/topics/"+strconv.FormatInt(topicID, 10), bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var respBody map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &respBody)
	assert.NoError(t, err)
	assert.Equal(t, "post updated", respBody["message"])
	mockUsecase.AssertExpectations(t)
}

func TestTopicHandler_Update_NoUserIDInContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockUsecase := mocks.NewTopicUsecase(t)
	logger := zerolog.Nop()
	handler := &TopicHandler{
		Usecase: mockUsecase,
		Log:     &logger,
	}
	topicID := int64(1)
	router.PUT("/topics/:id", handler.Update)

	reqBody := requests.UpdateRequestTopic{Title: "Updated Topic Title"}
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPut, "/topics/"+strconv.FormatInt(topicID, 10), bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusForbidden, rr.Code)
	mockUsecase.AssertNotCalled(t, "Update", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestTopicHandler_Update_InvalidTopicID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockUsecase := mocks.NewTopicUsecase(t)
	logger := zerolog.Nop()
	handler := &TopicHandler{
		Usecase: mockUsecase,
		Log:     &logger,
	}
	userID := int64(10)
	userRole := "user"
	router.PUT("/topics/:id", func(c *gin.Context) {
		c.Set(ContextUserIDKey, userID)
		c.Set(ContextRoleKey, userRole)
		handler.Update(c)
	})

	reqBody := requests.UpdateRequestTopic{Title: "Updated Topic Title"}
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPut, "/topics/invalid", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	mockUsecase.AssertNotCalled(t, "Update", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestTopicHandler_Update_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockUsecase := mocks.NewTopicUsecase(t)
	logger := zerolog.Nop()
	handler := &TopicHandler{
		Usecase: mockUsecase,
		Log:     &logger,
	}
	topicID := int64(1)
	userID := int64(10)
	userRole := "user"
	router.PUT("/topics/:id", func(c *gin.Context) {
		c.Set(ContextUserIDKey, userID)
		c.Set(ContextRoleKey, userRole)
		handler.Update(c)
	})

	req, _ := http.NewRequest(http.MethodPut, "/topics/"+strconv.FormatInt(topicID, 10), bytes.NewBufferString("{invalid_json"))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	mockUsecase.AssertNotCalled(t, "Update", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestTopicHandler_Update_Forbidden(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockUsecase := mocks.NewTopicUsecase(t)
	logger := zerolog.Nop()
	handler := &TopicHandler{
		Usecase: mockUsecase,
		Log:     &logger,
	}
	topicID := int64(1)
	userID := int64(10)
	userRole := "user"
	router.PUT("/topics/:id", func(c *gin.Context) {
		c.Set(ContextUserIDKey, userID)
		c.Set(ContextRoleKey, userRole)
		handler.Update(c)
	})

	reqBody := requests.UpdateRequestTopic{Title: "Updated Topic Title"}
	usecaseError := usecase.ErrForbidden
	mockUsecase.On("Update", mock.Anything, topicID, userID, userRole, reqBody.Title).Return(usecaseError).Once()

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPut, "/topics/"+strconv.FormatInt(topicID, 10), bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusForbidden, rr.Code)
	var respBody map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &respBody)
	assert.NoError(t, err)
	assert.Equal(t, "insufficient permissions", respBody["error"])
	mockUsecase.AssertExpectations(t)
}

func TestTopicHandler_Update_TopicOrPostNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockUsecase := mocks.NewTopicUsecase(t)
	logger := zerolog.Nop()
	handler := &TopicHandler{
		Usecase: mockUsecase,
		Log:     &logger,
	}
	topicID := int64(1)
	userID := int64(10)
	userRole := "user"
	router.PUT("/topics/:id", func(c *gin.Context) {
		c.Set(ContextUserIDKey, userID)
		c.Set(ContextRoleKey, userRole)
		handler.Update(c)
	})

	reqBody := requests.UpdateRequestTopic{Title: "Updated Topic Title"}
	usecaseError := usecase.ErrTopicNotFound
	mockUsecase.On("Update", mock.Anything, topicID, userID, userRole, reqBody.Title).Return(usecaseError).Once()

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPut, "/topics/"+strconv.FormatInt(topicID, 10), bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "internal server error")

	mockUsecase.AssertExpectations(t)
}

func TestTopicHandler_Update_UsecaseError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockUsecase := mocks.NewTopicUsecase(t)
	logger := zerolog.Nop()
	handler := &TopicHandler{
		Usecase: mockUsecase,
		Log:     &logger,
	}
	topicID := int64(1)
	userID := int64(10)
	userRole := "user"
	router.PUT("/topics/:id", func(c *gin.Context) {
		c.Set(ContextUserIDKey, userID)
		c.Set(ContextRoleKey, userRole)
		handler.Update(c)
	})

	reqBody := requests.UpdateRequestTopic{Title: "Updated Topic Title"}
	usecaseError := errors.New("some other update error")
	mockUsecase.On("Update", mock.Anything, topicID, userID, userRole, reqBody.Title).Return(usecaseError).Once()

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPut, "/topics/"+strconv.FormatInt(topicID, 10), bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	mockUsecase.AssertExpectations(t)
}

func TestTopicHandler_Delete_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockUsecase := mocks.NewTopicUsecase(t)
	logger := zerolog.Nop()
	handler := &TopicHandler{
		Usecase: mockUsecase,
		Log:     &logger,
	}
	topicID := int64(1)
	userID := int64(10)
	userRole := "user"

	router.DELETE("/topics/:id", func(c *gin.Context) {
		c.Set(ContextUserIDKey, userID)
		c.Set(ContextRoleKey, userRole)
		handler.Delete(c)
	})

	mockUsecase.On("Delete", mock.Anything, topicID, userID, userRole).Return(nil).Once()

	req, _ := http.NewRequest(http.MethodDelete, "/topics/"+strconv.FormatInt(topicID, 10), nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var respBody map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &respBody)
	assert.NoError(t, err)
	assert.Equal(t, "topic deleted", respBody["message"])
	mockUsecase.AssertExpectations(t)
}

func TestTopicHandler_Delete_NoUserIDInContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockUsecase := mocks.NewTopicUsecase(t)
	logger := zerolog.Nop()
	handler := &TopicHandler{
		Usecase: mockUsecase,
		Log:     &logger,
	}
	topicID := int64(1)
	router.DELETE("/topics/:id", handler.Delete) // userID не установлен

	req, _ := http.NewRequest(http.MethodDelete, "/topics/"+strconv.FormatInt(topicID, 10), nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusForbidden, rr.Code)
	mockUsecase.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestTopicHandler_Delete_InvalidTopicID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockUsecase := mocks.NewTopicUsecase(t)
	logger := zerolog.Nop()
	handler := &TopicHandler{
		Usecase: mockUsecase,
		Log:     &logger,
	}
	userID := int64(10)
	userRole := "user"
	router.DELETE("/topics/:id", func(c *gin.Context) {
		c.Set(ContextUserIDKey, userID)
		c.Set(ContextRoleKey, userRole)
		handler.Delete(c)
	})

	req, _ := http.NewRequest(http.MethodDelete, "/topics/invalid", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	mockUsecase.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestTopicHandler_Delete_Forbidden(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockUsecase := mocks.NewTopicUsecase(t)
	logger := zerolog.Nop()
	handler := &TopicHandler{
		Usecase: mockUsecase,
		Log:     &logger,
	}
	topicID := int64(1)
	userID := int64(10)
	userRole := "user"
	router.DELETE("/topics/:id", func(c *gin.Context) {
		c.Set(ContextUserIDKey, userID)
		c.Set(ContextRoleKey, userRole)
		handler.Delete(c)
	})

	usecaseError := usecase.ErrForbidden
	mockUsecase.On("Delete", mock.Anything, topicID, userID, userRole).Return(usecaseError).Once()

	req, _ := http.NewRequest(http.MethodDelete, "/topics/"+strconv.FormatInt(topicID, 10), nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusForbidden, rr.Code)
	mockUsecase.AssertExpectations(t)
}

func TestTopicHandler_Delete_TopicOrPostNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockUsecase := mocks.NewTopicUsecase(t)
	logger := zerolog.Nop()
	handler := &TopicHandler{
		Usecase: mockUsecase,
		Log:     &logger,
	}
	topicID := int64(1)
	userID := int64(10)
	userRole := "user"
	router.DELETE("/topics/:id", func(c *gin.Context) {
		c.Set(ContextUserIDKey, userID)
		c.Set(ContextRoleKey, userRole)
		handler.Delete(c)
	})

	usecaseError := usecase.ErrTopicNotFound
	mockUsecase.On("Delete", mock.Anything, topicID, userID, userRole).Return(usecaseError).Once()

	req, _ := http.NewRequest(http.MethodDelete, "/topics/"+strconv.FormatInt(topicID, 10), nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "internal server error")

	mockUsecase.AssertExpectations(t)
}

func TestTopicHandler_Delete_UsecaseError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockUsecase := mocks.NewTopicUsecase(t)
	logger := zerolog.Nop()
	handler := &TopicHandler{
		Usecase: mockUsecase,
		Log:     &logger,
	}
	topicID := int64(1)
	userID := int64(10)
	userRole := "user"
	router.DELETE("/topics/:id", func(c *gin.Context) {
		c.Set(ContextUserIDKey, userID)
		c.Set(ContextRoleKey, userRole)
		handler.Delete(c)
	})

	usecaseError := errors.New("some other delete error")
	mockUsecase.On("Delete", mock.Anything, topicID, userID, userRole).Return(usecaseError).Once()

	req, _ := http.NewRequest(http.MethodDelete, "/topics/"+strconv.FormatInt(topicID, 10), nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	mockUsecase.AssertExpectations(t)
}
