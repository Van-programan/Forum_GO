package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Van-programan/Forum_GO/internal/client"
	"github.com/Van-programan/Forum_GO/internal/controller/middleware"
	"github.com/Van-programan/Forum_GO/internal/controller/request"
	"github.com/Van-programan/Forum_GO/internal/entity"
	"github.com/Van-programan/Forum_GO/internal/usecase"
	"github.com/Van-programan/Forum_GO/internal/ws"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		return origin == "http://localhost:5173"
	},
}

type CategoryHandler struct {
	Usecase usecase.CategoryUsecase
	Log     *zerolog.Logger
}

type PostHandler struct {
	Usecase usecase.PostUsecase
}

type TopicHandler struct {
	Usecase usecase.TopicUsecase
	Log     *zerolog.Logger
}

type ChatHandler struct {
	hub         *ws.Hub
	chatUsecase usecase.ChatUsecase
	userClient  client.UserClient
	log         *zerolog.Logger
}

const (
	createOp   = "CategoryHandler.Create"
	getTitleOp = "CategoryHandler.GetTitle"
	getAllOp   = "CategoryHandler.GetAll"
	deleteOp   = "CategoryHandler.Delete"
	updateOp   = "CategoryHandler.Update"
)

const (
	createTopicOp   = "TopicHandler.Create"
	getByCategoryOp = "TopicHandler.GetByCategory"
	deleteTopicOp   = "TopicHandler.Delete"
	updateTopicOp   = "TopicHandler.Update"
	getByIDTopicOP  = "TopicHandler.GetByID"
)

func NewChatHandler(hub *ws.Hub, chatUsecase usecase.ChatUsecase, userClient client.UserClient, log *zerolog.Logger) *ChatHandler {
	return &ChatHandler{hub: hub, chatUsecase: chatUsecase, userClient: userClient, log: log}
}

// Create godoc
// @Summary Create a new category
// @Description Creates a new category. Requires admin role.
// @Tags categories
// @Accept json
// @Produce json
// @Param category body entity.Category true "Category data to create. ID, CreatedAt, UpdatedAt will be ignored."
// @Success 201 {object} response.IDResponse "Category created successfully"
// @Failure 400 {object} response.ErrorResponseForum "Invalid request payload"
// @Failure 401 {object} response.ErrorResponseForum "Unauthorized (token is missing or invalid)"
// @Failure 403 {object} response.ErrorResponseForum "Forbidden (user is not an admin)"
// @Failure 500 {object} response.ErrorResponseForum "Internal server error"
// @Security ApiKeyAuth
// @Router /categories [post]
func (h *CategoryHandler) Create(c *gin.Context) {
	log := h.getRequestLogger(c).With().Str("op", createOp).Logger()

	var category entity.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		log.Warn().Err(err).Msg("Failed to bind request")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.Usecase.Create(c.Request.Context(), category)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create category")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// GetByID godoc
// @Summary Get a category by ID
// @Description Retrieves a specific category by its ID.
// @Tags categories
// @Produce json
// @Param id path int true "Category ID" Format(int64)
// @Success 200 {object} response.CategoryResponse "Successfully retrieved category"
// @Failure 400 {object} response.ErrorResponseForum "Invalid category ID"
// @Failure 500 {object} response.ErrorResponseForum "Failed to get category"
// @Router /categories/{id} [get]
func (h *CategoryHandler) GetByID(c *gin.Context) {
	log := h.getRequestLogger(c).With().Str("op", getTitleOp).Logger()

	categoryID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to parse category id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category id"})
		return
	}

	category, err := h.Usecase.GetByID(c.Request.Context(), categoryID)
	if err != nil {
		log.Error().Err(err).Int64("category_id", categoryID).Msg("Failed to get category")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"category": category})

}

// GetAll godoc
// @Summary Get all categories
// @Description Retrieves a list of all categories.
// @Tags categories
// @Produce json
// @Success 200 {object} response.CategoriesResponse "Successfully retrieved all categories"
// @Failure 500 {object} response.ErrorResponseForum "Internal server error"
// @Router /categories [get]
func (h *CategoryHandler) GetAll(c *gin.Context) {
	log := h.getRequestLogger(c).With().Str("op", getAllOp).Logger()

	posts, err := h.Usecase.GetAll(c.Request.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to get all categories")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"categories": posts})
}

// Delete godoc
// @Summary Delete a category
// @Description Deletes a category by its ID. Requires admin privileges.
// @Tags categories
// @Param id path int true "Category ID" Format(int64)
// @Success 200 "Category deleted successfully"
// @Failure 400 {object} response.ErrorResponseForum "Invalid category ID"
// @Failure 401 {object} response.ErrorResponseForum "Unauthorized (token is missing or invalid)"
// @Failure 403 {object} response.ErrorResponseForum "Forbidden (user is not an admin)"
// @Failure 500 {object} response.ErrorResponseForum "Failed to delete category"
// @Security ApiKeyAuth
// @Router /categories/{id} [delete]
func (h *CategoryHandler) Delete(c *gin.Context) {
	log := h.getRequestLogger(c).With().Str("op", deleteOp).Logger()

	categoryID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to parse category id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category id"})
		return
	}

	if err := h.Usecase.Delete(c.Request.Context(), categoryID); err != nil {
		log.Error().Err(err).Msg("Failed to delete category")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete category"})
		return
	}

	c.Status(http.StatusOK)
}

// Update godoc
// @Summary Update a category
// @Description Updates a category's title and/or description by its ID. Requires admin privileges.
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID" Format(int64)
// @Param category_update body request.UpdateRequestCategory true "Category update data"
// @Success 200 "Category updated successfully"
// @Failure 400 {object} response.ErrorResponseForum "Invalid category ID or request payload"
// @Failure 401 {object} response.ErrorResponseForum "Unauthorized (token is missing or invalid)"
// @Failure 403 {object} response.ErrorResponseForum "Forbidden (user is not an admin)"
// @Failure 500 {object} response.ErrorResponseForum "Failed to update category"
// @Security ApiKeyAuth
// @Router /categories/{id} [patch]
func (h *CategoryHandler) Update(c *gin.Context) {
	log := h.getRequestLogger(c).With().Str("op", updateOp).Logger()

	categoryID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to parse category id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category id"})
		return
	}

	var req request.UpdateRequestCategory
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("Failed to bind request")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Usecase.Update(c.Request.Context(), categoryID, req.Title, req.Description); err != nil {
		log.Error().Err(err).Msg("Failed to update category")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update category"})
		return
	}

	c.Status(http.StatusOK)
}

func (h *CategoryHandler) getRequestLogger(c *gin.Context) *zerolog.Logger {
	reqLog := h.Log.With().
		Str("method", c.Request.Method).
		Str("path", c.Request.URL.Path).
		Str("remote_addr", c.ClientIP())

	logger := reqLog.Logger()
	return &logger
}

func (h *ChatHandler) ServeWs(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.log.Error().Err(err).Str("op", "ChatHandler.ServeWs").Msg("Failed to upgrade connection")
		return
	}

	if !exists {
		client := ws.NewUnauthorizedClient(h.hub, conn, h.chatUsecase)
		h.hub.Register <- client
		go client.WritePump()
		go client.ReadPump()
		return
	}

	username, err := h.userClient.GetUsername(c.Request.Context(), userID)
	if err != nil {
		h.log.Error().Err(err).Str("op", "ChatHandler.ServeWs").Msg("Failed to get username")
		conn.Close()
		return
	}

	client := ws.NewAuthorizedClient(h.hub, conn, userID, username, h.chatUsecase)
	h.hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
	c.JSON(http.StatusOK, gin.H{"message": "Connected to chat"})
}

func (h *PostHandler) Create(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
		return
	}

	topicID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category id"})
		return
	}

	var post entity.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post.TopicID = topicID
	post.AuthorID = &userID

	id, err := h.Usecase.Create(c.Request.Context(), post)
	if err != nil {
		if errors.Is(err, usecase.ErrTopicNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

// GetByTopic godoc
// @Summary Get posts by topic ID
// @Description Retrieves a list of posts for a topic ID.
// @Tags posts
// @Produce json
// @Param id path int true "Topic ID" Format(int64)
// @Success 200 {object} response.PostsResponse "Successfully retrieved posts"
// @Failure 400 {object} response.ErrorResponseForum "Invalid topic ID"
// @Failure 404 {object} response.ErrorResponseForum "Topic not found"
// @Failure 500 {object} response.ErrorResponseForum "Internal server error"
// @Router /topics/{id}/posts [get]
func (h *PostHandler) GetByTopic(c *gin.Context) {
	topicID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid topic id"})
		return
	}

	posts, err := h.Usecase.GetByTopic(c.Request.Context(), topicID)
	if err != nil {
		if errors.Is(err, usecase.ErrTopicNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"posts": posts})
}

// Update godoc
// @Summary Update a post
// @Description Updates a post. Requires authentication and ownership or admin role.
// @Tags posts
// @Accept json
// @Produce json
// @Param id path int true "Post ID" Format(int64)
// @Param post_update body request.UpdateRequestPost true "Post update data (only content)"
// @Success 200 {object} response.SuccessMessageResponse "Post updated successfully"
// @Failure 400 {object} response.ErrorResponseForum "Invalid post ID or request payload"
// @Failure 401 {object} response.ErrorResponseForum "Unauthorized (token is missing or invalid)"
// @Failure 403 {object} response.ErrorResponseForum "Forbidden (user is not an owner or admin)"
// @Failure 404 {object} response.ErrorResponseForum "Post not found"
// @Failure 500 {object} response.ErrorResponseForum "Internal server error"
// @Security ApiKeyAuth
// @Router /posts/{id} [patch]
func (h *PostHandler) Update(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
		return
	}
	role, _ := middleware.GetRoleFromContext(c)

	postID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	var req request.UpdateRequestPost
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.Usecase.Update(c.Request.Context(), postID, userID, role, req.Content)
	if err != nil {
		if errors.Is(err, usecase.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			return
		}
		if errors.Is(err, usecase.ErrPostNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "post updated"})
}

// Delete godoc
// @Summary Delete a post
// @Description Deletes a post by its ID. Requires authentication and ownership or admin role.
// @Tags posts
// @Param id path int true "Post ID" Format(int64)
// @Success 200 {object} response.SuccessMessageResponse "Post deleted successfully"
// @Failure 400 {object} response.ErrorResponseForum "Invalid post ID"
// @Failure 401 {object} response.ErrorResponseForum "Unauthorized (token is missing or invalid)"
// @Failure 403 {object} response.ErrorResponseForum "Forbidden (user is not an owner or admin)"
// @Failure 404 {object} response.ErrorResponseForum "Post not found"
// @Failure 500 {object} response.ErrorResponseForum "Internal server error"
// @Security ApiKeyAuth
// @Router /posts/{id} [delete]
func (h *PostHandler) Delete(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
		return
	}
	role, _ := middleware.GetRoleFromContext(c)

	postID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	err = h.Usecase.Delete(c.Request.Context(), postID, userID, role)
	if err != nil {
		if errors.Is(err, usecase.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			return
		}
		if errors.Is(err, usecase.ErrPostNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "post deleted"})
}

// Create godoc
// @Summary Create a new topic
// @Description Creates a new topic in a category. Requires authentication.
// @Tags topics
// @Accept json
// @Produce json
// @Param id path int true "Category ID to create topic in" Format(int64)
// @Param topic body entity.Topic true "Topic data to create. ID, AuthorID, CategoryID, CreatedAt, UpdatedAt will be ignored or overridden."
// @Success 200 {object} response.IDResponse "Topic created successfully"
// @Failure 400 {object} response.ErrorResponseForum "Invalid category ID or request payload"
// @Failure 401 {object} response.ErrorResponseForum "Unauthorized (token is missing or invalid)"
// @Failure 403 {object} response.ErrorResponseForum "Forbidden (user is not authorized or trying to impersonate)"
// @Failure 500 {object} response.ErrorResponseForum "Internal server error"
// @Security ApiKeyAuth
// @Router /categories/{id}/topics [post]
func (h *TopicHandler) Create(c *gin.Context) {
	log := h.getRequestLogger(c).With().Str("op", createTopicOp).Logger()

	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		log.Warn().Msg("insufficient permissions")
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
		return
	}

	categoryID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Warn().Msg("invalid category id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category id"})
		return
	}

	var topic entity.Topic
	if err := c.ShouldBindJSON(&topic); err != nil {
		log.Warn().Err(err).Msg("failed to bind request")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	topic.AuthorID = &userID
	topic.CategoryID = categoryID

	id, err := h.Usecase.Create(c.Request.Context(), topic)
	if err != nil {
		if errors.Is(err, usecase.ErrCategoryNotFound) {
			log.Warn().Msg("category not found")
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		log.Error().Err(err).Msg("failed to create topic")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

// GetByID godoc
// @Summary Get a topic by ID
// @Description Retrieves a specific topic by its ID.
// @Tags topics
// @Produce json
// @Param id path int true "Topic ID" Format(int64)
// @Success 200 {object} response.TopicResponse "Successfully retrieved topic"
// @Failure 400 {object} response.ErrorResponseForum "Invalid topic ID"
// @Failure 500 {object} response.ErrorResponseForum "Failed to get topic"
// @Router /topics/{id} [get]
func (h *TopicHandler) GetByID(c *gin.Context) {
	log := h.getRequestLogger(c).With().Str("op", getByIDTopicOP).Logger()

	topicID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to parse topic id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid topic id"})
		return
	}

	topic, err := h.Usecase.GetByID(c.Request.Context(), topicID)
	if err != nil {
		log.Error().Err(err).Int64("topic_id", topicID).Msg("Failed to get topic")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get topic"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"topic": topic})

}

// GetByCategory godoc
// @Summary Get topics by category ID
// @Description Retrieves a list of topics for a category ID.
// @Tags topics
// @Produce json
// @Param id path int true "Category ID" Format(int64)
// @Success 200 {object} response.TopicsResponse "Successfully retrieved topics"
// @Failure 400 {object} response.ErrorResponseForum "Invalid category ID"
// @Failure 404 {object} response.ErrorResponseForum "Category not found"
// @Failure 500 {object} response.ErrorResponseForum "Internal server error"
// @Router /categories/{id}/topics [get]
func (h *TopicHandler) GetByCategory(c *gin.Context) {
	log := h.getRequestLogger(c).With().Str("op", getByCategoryOp).Logger()

	categoryID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Warn().Msg("invalid category id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category id"})
		return
	}

	topics, err := h.Usecase.GetByCategory(c.Request.Context(), categoryID)
	if err != nil {
		if errors.Is(err, usecase.ErrCategoryNotFound) {
			log.Warn().Msg("category not found")
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		log.Error().Err(err).Msg("failed to get topics by category")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"topics": topics})
}

// Update godoc
// @Summary Update a topic
// @Description Updates a topic. Requires authentication and ownership or admin role.
// @Tags topics
// @Accept json
// @Produce json
// @Param id path int true "Topic ID" Format(int64)
// @Param topic_update body request.UpdateRequestTopic true "Topic update data (only title)"
// @Success 200 {object} response.SuccessMessageResponse "Topic updated successfully"
// @Failure 400 {object} response.ErrorResponseForum "Invalid topic ID or request payload"
// @Failure 401 {object} response.ErrorResponseForum "Unauthorized (token is missing or invalid)"
// @Failure 403 {object} response.ErrorResponseForum "Forbidden (user is not an owner or admin)"
// @Failure 404 {object} response.ErrorResponseForum "Topic not found"
// @Failure 500 {object} response.ErrorResponseForum "Internal server error"
// @Security ApiKeyAuth
// @Router /topics/{id} [patch]
func (h *TopicHandler) Update(c *gin.Context) {
	log := h.getRequestLogger(c).With().Str("op", updateTopicOp).Logger()

	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		log.Warn().Msg("insufficient permissions")
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
		return
	}
	role, _ := middleware.GetRoleFromContext(c)

	topicID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid topic id"})
		return
	}

	var req request.UpdateRequestTopic
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("failed to bind request")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.Usecase.Update(c.Request.Context(), topicID, userID, role, req.Title)
	if err != nil {
		if errors.Is(err, usecase.ErrForbidden) {
			log.Warn().Msg("insufficient permissions")
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			return
		}
		if errors.Is(err, usecase.ErrPostNotFound) {
			log.Warn().Msg("post not found")
			c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "post updated"})
}

// Delete godoc
// @Summary Delete a topic
// @Description Deletes a topic by its ID. Requires authentication and ownership or role.
// @Tags topics
// @Param id path int true "Topic ID" Format(int64)
// @Success 200 {object} response.SuccessMessageResponse "Topic deleted successfully"
// @Failure 400 {object} response.ErrorResponseForum "Invalid topic ID"
// @Failure 401 {object} response.ErrorResponseForum "Unauthorized (token is missing or invalid)"
// @Failure 403 {object} response.ErrorResponseForum "Forbidden (user is not an owner or admin)"
// @Failure 404 {object} response.ErrorResponseForum "Topic not found"
// @Failure 500 {object} response.ErrorResponseForum "Internal server error"
// @Security ApiKeyAuth
// @Router /topics/{id} [delete]
func (h *TopicHandler) Delete(c *gin.Context) {
	log := h.getRequestLogger(c).With().Str("op", deleteTopicOp).Logger()
	userID, exists := middleware.GetUserIDFromContext(c)

	if !exists {
		log.Warn().Msg("insufficient permissions")
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
		return
	}
	role, _ := middleware.GetRoleFromContext(c)

	topicID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Warn().Msg("invalid topic id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid topic id"})
		return
	}

	err = h.Usecase.Delete(c.Request.Context(), topicID, userID, role)
	if err != nil {
		if errors.Is(err, usecase.ErrForbidden) {
			log.Warn().Msg("insufficient permissions")
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			return
		}
		if errors.Is(err, usecase.ErrPostNotFound) {
			log.Warn().Msg("post not found")
			c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
			return
		}
		log.Error().Err(err).Msg("failed to delete topic")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "topic deleted"})
}

func (h *TopicHandler) getRequestLogger(c *gin.Context) *zerolog.Logger {
	reqLog := h.Log.With().
		Str("method", c.Request.Method).
		Str("path", c.Request.URL.Path).
		Str("remote_addr", c.ClientIP())

	logger := reqLog.Logger()
	return &logger
}
