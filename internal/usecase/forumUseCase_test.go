package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Van-programan/Forum_GO/internal/entity"
	mocksf "github.com/Van-programan/Forum_GO/mocks/forum"
	mocks "github.com/Van-programan/Forum_GO/mocks/forum/repository"
	"github.com/jackc/pgx/v5"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type CategoryUsecaseSuite struct {
	suite.Suite
	usecase  CategoryUsecase
	repoMock *mocks.CategoryRepository
	log      *zerolog.Logger
}

type ChatUsecaseSuite struct {
	suite.Suite
	usecase      ChatUsecase
	chatRepoMock *mocks.ChatRepository
	log          *zerolog.Logger
}

type PostUsecaseSuite struct {
	suite.Suite
	usecase         PostUsecase
	postRepoMock    *mocks.PostRepository
	topicRepoMock   *mocks.TopicRepository
	userClientMock  *mocksf.UserClient
	log             *zerolog.Logger
	defaultAuthorID int64
}

type TopicUsecaseSuite struct {
	suite.Suite
	usecase           TopicUsecase
	topicRepoMock     *mocks.TopicRepository
	categoryRepoMock  *mocks.CategoryRepository
	userClientMock    *mocksf.UserClient
	log               *zerolog.Logger
	defaultAuthorID   int64
	defaultCategoryID int64
}

func (s *CategoryUsecaseSuite) SetupTest() {
	s.repoMock = mocks.NewCategoryRepository(s.T())
	logger := zerolog.Nop()
	s.log = &logger
	s.usecase = NewCategoryUsecase(s.repoMock, s.log)
}

func TestCategoryUsecaseSuite(t *testing.T) {
	suite.Run(t, new(CategoryUsecaseSuite))
}

// Create
func (s *CategoryUsecaseSuite) TestCreateCategory_Success() {
	ctx := context.Background()
	category := entity.Category{Title: "New Category", Description: "Description"}
	expectedID := int64(1)

	s.repoMock.On("Create", ctx, category).Return(expectedID, nil).Once()

	id, err := s.usecase.Create(ctx, category)

	s.NoError(err)
	s.Equal(expectedID, id)
	s.repoMock.AssertExpectations(s.T())
}

func (s *CategoryUsecaseSuite) TestCreateCategory_RepoError() {
	ctx := context.Background()
	category := entity.Category{Title: "New Category", Description: "Description"}
	expectedError := errors.New("repository error")

	s.repoMock.On("Create", ctx, category).Return(int64(0), expectedError).Once()

	id, err := s.usecase.Create(ctx, category)

	s.Error(err)
	s.Equal(int64(0), id)
	s.Contains(err.Error(), "ForumService - CategoryUsecase - Create - repo.Create()")
	s.ErrorIs(err, expectedError)
	s.repoMock.AssertExpectations(s.T())
}

// GetByID
func (s *CategoryUsecaseSuite) TestGetByIDCategory_Success() {
	ctx := context.Background()
	categoryID := int64(1)
	expectedCategory := &entity.Category{ID: categoryID, Title: "Test Category", Description: "Test Description", CreatedAt: time.Now()}

	s.repoMock.On("GetByID", ctx, categoryID).Return(expectedCategory, nil).Once()

	category, err := s.usecase.GetByID(ctx, categoryID)

	s.NoError(err)
	s.NotNil(category)
	s.Equal(expectedCategory, category)
	s.repoMock.AssertExpectations(s.T())
}

func (s *CategoryUsecaseSuite) TestGetByIDCategory_RepoError() {
	ctx := context.Background()
	categoryID := int64(1)
	expectedError := errors.New("repository error")

	s.repoMock.On("GetByID", ctx, categoryID).Return(nil, expectedError).Once()

	category, err := s.usecase.GetByID(ctx, categoryID)

	s.Error(err)
	s.Nil(category)
	s.Contains(err.Error(), "ForumService - CategoryUsecase - GetByID - repo.GetByID()")
	s.ErrorIs(err, expectedError)
	s.repoMock.AssertExpectations(s.T())
}

// GetAll
func (s *CategoryUsecaseSuite) TestGetAllCategories_Success() {
	ctx := context.Background()
	expectedCategories := []entity.Category{
		{ID: 1, Title: "Category 1", Description: "Desc 1", CreatedAt: time.Now()},
		{ID: 2, Title: "Category 2", Description: "Desc 2", CreatedAt: time.Now()},
	}

	s.repoMock.On("GetAll", ctx).Return(expectedCategories, nil).Once()

	categories, err := s.usecase.GetAll(ctx)

	s.NoError(err)
	s.NotNil(categories)
	s.Equal(expectedCategories, categories)
	s.repoMock.AssertExpectations(s.T())
}

func (s *CategoryUsecaseSuite) TestGetAllCategories_RepoError() {
	ctx := context.Background()
	expectedError := errors.New("repository error")

	s.repoMock.On("GetAll", ctx).Return(nil, expectedError).Once()

	categories, err := s.usecase.GetAll(ctx)

	s.Error(err)
	s.Nil(categories)
	s.Contains(err.Error(), "ForumService - CategoryUsecase - GetAll - repo.GetAll()")
	s.ErrorIs(err, expectedError)
	s.repoMock.AssertExpectations(s.T())
}

// Update
func (s *CategoryUsecaseSuite) TestUpdateCategory_Success() {
	ctx := context.Background()
	categoryID := int64(1)
	title := "Updated Title"
	description := "Updated Description"

	s.repoMock.On("Update", ctx, categoryID, title, description).Return(nil).Once()

	err := s.usecase.Update(ctx, categoryID, title, description)

	s.NoError(err)
	s.repoMock.AssertExpectations(s.T())
}

func (s *CategoryUsecaseSuite) TestUpdateCategory_RepoError() {
	ctx := context.Background()
	categoryID := int64(1)
	title := "Updated Title"
	description := "Updated Description"
	expectedError := errors.New("repository error")

	s.repoMock.On("Update", ctx, categoryID, title, description).Return(expectedError).Once()

	err := s.usecase.Update(ctx, categoryID, title, description)

	s.Error(err)
	s.Contains(err.Error(), "ForumService - CategoryUsecase - Update - repo.Update()")
	s.ErrorIs(err, expectedError)
	s.repoMock.AssertExpectations(s.T())
}

// Delete
func (s *CategoryUsecaseSuite) TestDeleteCategory_Success() {
	ctx := context.Background()
	categoryID := int64(1)

	s.repoMock.On("Delete", ctx, categoryID).Return(nil).Once()

	err := s.usecase.Delete(ctx, categoryID)

	s.NoError(err)
	s.repoMock.AssertExpectations(s.T())
}

func (s *CategoryUsecaseSuite) TestDeleteCategory_RepoError() {
	ctx := context.Background()
	categoryID := int64(1)
	expectedError := errors.New("repository error")

	s.repoMock.On("Delete", ctx, categoryID).Return(expectedError).Once()

	err := s.usecase.Delete(ctx, categoryID)

	s.Error(err)
	s.Contains(err.Error(), "ForumService - CategoryUsecase - Delete - repo.Delete()")
	s.ErrorIs(err, expectedError)
	s.repoMock.AssertExpectations(s.T())
}

func (s *ChatUsecaseSuite) SetupTest() {
	s.chatRepoMock = mocks.NewChatRepository(s.T())
	logger := zerolog.Nop()
	s.log = &logger
	s.usecase = NewChatUsecase(s.chatRepoMock, s.log)
}

func TestChatUsecaseSuite(t *testing.T) {
	suite.Run(t, new(ChatUsecaseSuite))
}

// GetMessageHistory
func (s *ChatUsecaseSuite) TestGetMessageHistory_Success() {
	ctx := context.Background()
	limit := int64(50)
	expectedMessages := []entity.ChatMessage{
		{ID: 1, UserID: 1, Username: "user1", Content: "msg1", CreatedAt: time.Now().Add(-time.Minute)},
		{ID: 2, UserID: 2, Username: "user2", Content: "msg2", CreatedAt: time.Now()},
	}

	s.chatRepoMock.On("GetMessages", ctx, limit).Return(expectedMessages, nil).Once()

	messages, err := s.usecase.GetMessageHistory(ctx, limit)

	s.NoError(err)
	s.NotNil(messages)
	s.Equal(expectedMessages, messages)
	s.chatRepoMock.AssertExpectations(s.T())
}

func (s *ChatUsecaseSuite) TestGetMessageHistory_RepoError() {
	ctx := context.Background()
	limit := int64(50)
	expectedError := errors.New("repository error")

	s.chatRepoMock.On("GetMessages", ctx, limit).Return(nil, expectedError).Once()

	messages, err := s.usecase.GetMessageHistory(ctx, limit)

	s.Error(err)
	s.Nil(messages)
	s.Contains(err.Error(), "ChatUsecase - GetMessageHistory - u.chatRepo.GetMessages()")
	s.ErrorIs(err, expectedError)
	s.chatRepoMock.AssertExpectations(s.T())
}

// SaveMessage
func (s *ChatUsecaseSuite) TestSaveMessage_Success() {
	ctx := context.Background()
	userID := int64(1)
	username := "TestUser"
	content := "This is a test message."
	expectedMessageID := int64(123)

	var capturedMessage *entity.ChatMessage
	s.chatRepoMock.On("SaveMessage", ctx, mock.MatchedBy(func(msg *entity.ChatMessage) bool {
		capturedMessage = msg
		return msg.UserID == userID && msg.Username == username && msg.Content == content
	})).Return(expectedMessageID, nil).Once()

	savedMessage, err := s.usecase.SaveMessage(ctx, userID, username, content)

	s.NoError(err)
	s.NotNil(savedMessage)
	s.Equal(userID, savedMessage.UserID)
	s.Equal(username, savedMessage.Username)
	s.Equal(content, savedMessage.Content)
	s.WithinDuration(time.Now(), savedMessage.CreatedAt, 2*time.Second)
	s.chatRepoMock.AssertExpectations(s.T())

	assert.NotNil(s.T(), capturedMessage)
	if capturedMessage != nil {
		assert.Equal(s.T(), userID, capturedMessage.UserID)
		assert.Equal(s.T(), username, capturedMessage.Username)
		assert.Equal(s.T(), content, capturedMessage.Content)
	}
}

func (s *ChatUsecaseSuite) TestSaveMessage_RepoError() {
	ctx := context.Background()
	userID := int64(1)
	username := "test-user"
	content := "test-message"
	expectedError := errors.New("repository error saving message")

	s.chatRepoMock.On("SaveMessage", ctx, mock.MatchedBy(func(msg *entity.ChatMessage) bool {
		return msg.UserID == userID && msg.Username == username && msg.Content == content
	})).Return(int64(0), expectedError).Once()

	savedMessage, err := s.usecase.SaveMessage(ctx, userID, username, content)

	s.Error(err)
	s.Nil(savedMessage)
	s.Contains(err.Error(), "ChatUsecase - SaveMessage - u.chatRepo.SaveMessage()")
	s.ErrorIs(err, expectedError)
	s.chatRepoMock.AssertExpectations(s.T())
}

func (s *PostUsecaseSuite) SetupTest() {
	s.postRepoMock = mocks.NewPostRepository(s.T())
	s.topicRepoMock = mocks.NewTopicRepository(s.T())
	s.userClientMock = mocksf.NewUserClient(s.T())
	logger := zerolog.Nop()
	s.log = &logger
	s.defaultAuthorID = int64(1)
	s.usecase = NewPostUsecase(s.postRepoMock, s.topicRepoMock, s.userClientMock, s.log)
}

func TestPostUsecaseSuite(t *testing.T) {
	suite.Run(t, new(PostUsecaseSuite))
}

// Create
func (s *PostUsecaseSuite) TestCreatePost_Success() {
	ctx := context.Background()
	post := entity.Post{TopicID: 1, AuthorID: &s.defaultAuthorID, Content: "content"}
	expectedPostID := int64(1)
	topic := &entity.Topic{ID: post.TopicID, Title: "Existing Topic"}

	s.topicRepoMock.On("GetByID", ctx, post.TopicID).Return(topic, nil).Once()
	s.postRepoMock.On("Create", ctx, post).Return(expectedPostID, nil).Once()

	id, err := s.usecase.Create(ctx, post)

	s.NoError(err)
	s.Equal(expectedPostID, id)
	s.topicRepoMock.AssertExpectations(s.T())
	s.postRepoMock.AssertExpectations(s.T())
}

func (s *PostUsecaseSuite) TestCreatePost_TopicNotFound() {
	ctx := context.Background()
	post := entity.Post{TopicID: 1, AuthorID: &s.defaultAuthorID, Content: "content"}
	expectedError := ErrTopicNotFound

	s.topicRepoMock.On("GetByID", ctx, post.TopicID).Return(nil, pgx.ErrNoRows).Once()

	id, err := s.usecase.Create(ctx, post)

	s.Error(err)
	s.Equal(int64(0), id)
	s.ErrorIs(err, expectedError)
	s.topicRepoMock.AssertExpectations(s.T())
	s.postRepoMock.AssertNotCalled(s.T(), "Create", mock.Anything, mock.Anything)
}

/*
func (s *PostUsecaseSuite) TestCreatePost_TopicRepoError_OtherThanNotFound() {
	ctx := context.Background()
	post := entity.Post{TopicID: 1, AuthorID: &s.defaultAuthorID, Content: "content"}
	repoError := errors.New("some topic repo error")

	s.topicRepoMock.On("GetByID", ctx, post.TopicID).Return(nil, repoError).Once()

	id, err := s.usecase.Create(ctx, post)

	s.Error(err)
	s.Equal(int64(0), id)
	s.ErrorIs(err, repoError)
	s.topicRepoMock.AssertExpectations(s.T())
	s.postRepoMock.AssertNotCalled(s.T(), "Create", mock.Anything, mock.Anything)
}*/

func (s *PostUsecaseSuite) TestCreatePost_RepoError() {
	ctx := context.Background()
	post := entity.Post{TopicID: 1, AuthorID: &s.defaultAuthorID, Content: "content"}
	expectedError := errors.New("post repository create error")
	topic := &entity.Topic{ID: post.TopicID, Title: "Existing Topic"}

	s.topicRepoMock.On("GetByID", ctx, post.TopicID).Return(topic, nil).Once()
	s.postRepoMock.On("Create", ctx, post).Return(int64(0), expectedError).Once()

	id, err := s.usecase.Create(ctx, post)

	s.Error(err)
	s.Equal(int64(0), id)
	s.Contains(err.Error(), "ForumService - PostUsecase - Create - postRepo.Create()")
	s.ErrorIs(err, expectedError)
	s.topicRepoMock.AssertExpectations(s.T())
	s.postRepoMock.AssertExpectations(s.T())
}

// GetByTopic
func (s *PostUsecaseSuite) TestGetByTopic_Success() {
	ctx := context.Background()
	topicID := int64(1)
	authorID1 := int64(10)
	authorID2 := int64(20)
	postsFromRepo := []entity.Post{
		{ID: 1, TopicID: topicID, AuthorID: &authorID1, Content: "Post 1", CreatedAt: time.Now()},
		{ID: 2, TopicID: topicID, AuthorID: &authorID2, Content: "Post 2", CreatedAt: time.Now()},
		{ID: 3, TopicID: topicID, AuthorID: nil, Content: "Post 3 - Deleted User", CreatedAt: time.Now()},
	}
	usernamesFromClient := map[int64]string{
		authorID1: "UserOne",
		authorID2: "UserTwo",
	}
	expectedPosts := []entity.Post{
		{ID: 1, TopicID: topicID, AuthorID: &authorID1, Username: "UserOne", Content: "Post 1", CreatedAt: postsFromRepo[0].CreatedAt},
		{ID: 2, TopicID: topicID, AuthorID: &authorID2, Username: "UserTwo", Content: "Post 2", CreatedAt: postsFromRepo[1].CreatedAt},
		{ID: 3, TopicID: topicID, AuthorID: nil, Username: "Удаленный пользователь", Content: "Post 3 - Deleted User", CreatedAt: postsFromRepo[2].CreatedAt},
	}
	topic := &entity.Topic{ID: topicID, Title: "Existing Topic"}

	s.topicRepoMock.On("GetByID", ctx, topicID).Return(topic, nil).Once()
	s.postRepoMock.On("GetByTopic", ctx, topicID).Return(postsFromRepo, nil).Once()
	s.userClientMock.On("GetUsernames", ctx, mock.MatchedBy(func(ids []int64) bool {
		return len(ids) == 2 && ((ids[0] == authorID1 && ids[1] == authorID2) || (ids[0] == authorID2 && ids[1] == authorID1))
	})).Return(usernamesFromClient, nil).Once()

	posts, err := s.usecase.GetByTopic(ctx, topicID)

	s.NoError(err)
	s.NotNil(posts)
	s.ElementsMatch(expectedPosts, posts)
	s.topicRepoMock.AssertExpectations(s.T())
	s.postRepoMock.AssertExpectations(s.T())
	s.userClientMock.AssertExpectations(s.T())
}

func (s *PostUsecaseSuite) TestGetByTopic_RepoError() {
	ctx := context.Background()
	topicID := int64(1)
	expectedError := errors.New("post repository GetByTopic error")
	topic := &entity.Topic{ID: topicID, Title: "Existing Topic"}

	s.topicRepoMock.On("GetByID", ctx, topicID).Return(topic, nil).Once()
	s.postRepoMock.On("GetByTopic", ctx, topicID).Return(nil, expectedError).Once()

	posts, err := s.usecase.GetByTopic(ctx, topicID)

	s.Error(err)
	s.Nil(posts)
	s.Contains(err.Error(), "ForumService - PostUsecase - GetByTopic - postRepo.GetByTopic()")
	s.ErrorIs(err, expectedError)
	s.topicRepoMock.AssertExpectations(s.T())
	s.postRepoMock.AssertExpectations(s.T())
	s.userClientMock.AssertNotCalled(s.T(), "GetUsernames", mock.Anything, mock.Anything)
}

func (s *PostUsecaseSuite) TestGetByTopic_UserClientError() {
	ctx := context.Background()
	topicID := int64(1)
	authorID1 := int64(10)
	postsFromRepo := []entity.Post{
		{ID: 1, TopicID: topicID, AuthorID: &authorID1, Content: "Post 1"},
	}
	expectedError := errors.New("user client error")
	topic := &entity.Topic{ID: topicID, Title: "Existing Topic"}

	s.topicRepoMock.On("GetByID", ctx, topicID).Return(topic, nil).Once()
	s.postRepoMock.On("GetByTopic", ctx, topicID).Return(postsFromRepo, nil).Once()
	s.userClientMock.On("GetUsernames", ctx, []int64{authorID1}).Return(nil, expectedError).Once()

	posts, err := s.usecase.GetByTopic(ctx, topicID)

	s.Error(err)
	s.Nil(posts)                                                                                        // В текущей реализации возвращается nil при ошибке клиента
	s.Contains(err.Error(), "ForumService - TopicUsecase  - GetByCategory - userClient.GetUsernames()") // Ошибка из TopicUsecase, т.к. логика похожа
	s.ErrorIs(err, expectedError)
	s.topicRepoMock.AssertExpectations(s.T())
	s.postRepoMock.AssertExpectations(s.T())
	s.userClientMock.AssertExpectations(s.T())
}

func (s *PostUsecaseSuite) TestGetByTopic_CheckTopicError() {
	ctx := context.Background()
	topicID := int64(1)
	expectedError := ErrTopicNotFound

	s.topicRepoMock.On("GetByID", ctx, topicID).Return(nil, pgx.ErrNoRows).Once() // Ошибка в checkTopic

	posts, err := s.usecase.GetByTopic(ctx, topicID)

	s.Error(err)
	s.Nil(posts)
	s.ErrorIs(err, expectedError)
	s.topicRepoMock.AssertExpectations(s.T())
	s.postRepoMock.AssertNotCalled(s.T(), "GetByTopic", mock.Anything, mock.Anything)
	s.userClientMock.AssertNotCalled(s.T(), "GetUsernames", mock.Anything, mock.Anything)
}

// Update
func (s *PostUsecaseSuite) TestUpdatePost_Success_Author() {
	ctx := context.Background()
	postID := int64(1)
	userID := s.defaultAuthorID
	role := "user"
	content := "updated content"
	postFromRepo := &entity.Post{ID: postID, AuthorID: &s.defaultAuthorID, Content: "old content"}

	s.postRepoMock.On("GetByID", ctx, postID).Return(postFromRepo, nil).Once()
	s.postRepoMock.On("Update", ctx, postID, content).Return(nil).Once()

	err := s.usecase.Update(ctx, postID, userID, role, content)

	s.NoError(err)
	s.postRepoMock.AssertExpectations(s.T())
}

func (s *PostUsecaseSuite) TestUpdatePost_Success_Admin() {
	ctx := context.Background()
	postID := int64(1)
	adminID := int64(999)
	otherUserID := s.defaultAuthorID
	role := "admin"
	content := "updated content by admin"
	postFromRepo := &entity.Post{ID: postID, AuthorID: &otherUserID, Content: "old content"}

	s.postRepoMock.On("GetByID", ctx, postID).Return(postFromRepo, nil).Once()
	s.postRepoMock.On("Update", ctx, postID, content).Return(nil).Once()

	err := s.usecase.Update(ctx, postID, adminID, role, content)

	s.NoError(err)
	s.postRepoMock.AssertExpectations(s.T())
}

func (s *PostUsecaseSuite) TestUpdatePost_AccessDenied_NotAuthorNotAdmin() {
	ctx := context.Background()
	postID := int64(1)
	anotherUserID := int64(555)
	authorID := s.defaultAuthorID
	role := "user"
	content := "updated content"
	postFromRepo := &entity.Post{ID: postID, AuthorID: &authorID, Content: "old content"}
	expectedError := ErrForbidden

	s.postRepoMock.On("GetByID", ctx, postID).Return(postFromRepo, nil).Once()

	err := s.usecase.Update(ctx, postID, anotherUserID, role, content)

	s.Error(err)
	s.ErrorIs(err, expectedError)
	s.postRepoMock.AssertExpectations(s.T())
	s.postRepoMock.AssertNotCalled(s.T(), "Update", mock.Anything, mock.Anything, mock.Anything)
}

func (s *PostUsecaseSuite) TestUpdatePost_PostNotFound_OnCheckAccess() {
	ctx := context.Background()
	postID := int64(1)
	userID := s.defaultAuthorID
	role := "user"
	content := "updated content"
	expectedError := ErrPostNotFound

	s.postRepoMock.On("GetByID", ctx, postID).Return(nil, pgx.ErrNoRows).Once()

	err := s.usecase.Update(ctx, postID, userID, role, content)

	s.Error(err)
	s.ErrorIs(err, expectedError)
	s.postRepoMock.AssertExpectations(s.T())
	s.postRepoMock.AssertNotCalled(s.T(), "Update", mock.Anything, mock.Anything, mock.Anything)
}

func (s *PostUsecaseSuite) TestUpdatePost_RepoUpdateError() {
	ctx := context.Background()
	postID := int64(1)
	userID := s.defaultAuthorID
	role := "user"
	content := "updated content"
	postFromRepo := &entity.Post{ID: postID, AuthorID: &s.defaultAuthorID, Content: "old content"}
	repoError := errors.New("repo update error")

	s.postRepoMock.On("GetByID", ctx, postID).Return(postFromRepo, nil).Once()
	s.postRepoMock.On("Update", ctx, postID, content).Return(repoError).Once()

	err := s.usecase.Update(ctx, postID, userID, role, content)

	s.Error(err)
	s.ErrorIs(err, repoError)
	s.Contains(err.Error(), "ForumService - PostUsecase - Update - postRepo.Update()")
	s.postRepoMock.AssertExpectations(s.T())
}

// Delete
func (s *PostUsecaseSuite) TestDeletePost_Success_Author() {
	ctx := context.Background()
	postID := int64(1)
	userID := s.defaultAuthorID
	role := "user"
	postFromRepo := &entity.Post{ID: postID, AuthorID: &s.defaultAuthorID}

	s.postRepoMock.On("GetByID", ctx, postID).Return(postFromRepo, nil).Once()
	s.postRepoMock.On("Delete", ctx, postID).Return(nil).Once()

	err := s.usecase.Delete(ctx, postID, userID, role)

	s.NoError(err)
	s.postRepoMock.AssertExpectations(s.T())
}

func (s *PostUsecaseSuite) TestDeletePost_Success_Admin() {
	ctx := context.Background()
	postID := int64(1)
	adminID := int64(999)
	otherUserID := s.defaultAuthorID
	role := "admin"
	postFromRepo := &entity.Post{ID: postID, AuthorID: &otherUserID}

	s.postRepoMock.On("GetByID", ctx, postID).Return(postFromRepo, nil).Once()
	s.postRepoMock.On("Delete", ctx, postID).Return(nil).Once()

	err := s.usecase.Delete(ctx, postID, adminID, role)

	s.NoError(err)
	s.postRepoMock.AssertExpectations(s.T())
}

func (s *PostUsecaseSuite) TestDeletePost_AccessDenied_NotAuthorNotAdmin() {
	ctx := context.Background()
	postID := int64(1)
	anotherUserID := int64(555)
	authorID := s.defaultAuthorID
	role := "user"
	postFromRepo := &entity.Post{ID: postID, AuthorID: &authorID}
	expectedError := ErrForbidden

	s.postRepoMock.On("GetByID", ctx, postID).Return(postFromRepo, nil).Once()

	err := s.usecase.Delete(ctx, postID, anotherUserID, role)

	s.Error(err)
	s.ErrorIs(err, expectedError)
	s.postRepoMock.AssertExpectations(s.T())
	s.postRepoMock.AssertNotCalled(s.T(), "Delete", mock.Anything, mock.Anything)
}

func (s *PostUsecaseSuite) TestDeletePost_PostNotFound_OnCheckAccess() {
	ctx := context.Background()
	postID := int64(1)
	userID := s.defaultAuthorID
	role := "user"
	expectedError := ErrPostNotFound

	s.postRepoMock.On("GetByID", ctx, postID).Return(nil, pgx.ErrNoRows).Once()

	err := s.usecase.Delete(ctx, postID, userID, role)

	s.Error(err)
	s.ErrorIs(err, expectedError)
	s.postRepoMock.AssertExpectations(s.T())
	s.postRepoMock.AssertNotCalled(s.T(), "Delete", mock.Anything, mock.Anything)
}

func (s *PostUsecaseSuite) TestDeletePost_RepoError() {
	ctx := context.Background()
	postID := int64(1)
	userID := s.defaultAuthorID
	role := "user"
	postFromRepo := &entity.Post{ID: postID, AuthorID: &s.defaultAuthorID}
	repoError := errors.New("repo delete error")

	s.postRepoMock.On("GetByID", ctx, postID).Return(postFromRepo, nil).Once()

	s.postRepoMock.On("Delete", ctx, postID).Return(repoError).Once()

	err := s.usecase.Delete(ctx, postID, userID, role)

	s.Error(err)
	s.ErrorIs(err, repoError)
	s.Contains(err.Error(), "ForumService - PostUsecase - Delete - postRepo.delete()")
	s.postRepoMock.AssertExpectations(s.T())
}

func (s *TopicUsecaseSuite) SetupTest() {
	s.topicRepoMock = mocks.NewTopicRepository(s.T())
	s.categoryRepoMock = mocks.NewCategoryRepository(s.T())
	s.userClientMock = mocksf.NewUserClient(s.T())
	logger := zerolog.Nop()
	s.log = &logger
	s.defaultAuthorID = int64(123)
	s.defaultCategoryID = int64(1)

	s.usecase = NewTopicUsecase(s.topicRepoMock, s.categoryRepoMock, s.userClientMock, s.log)
}

func TestTopicUsecaseSuite(t *testing.T) {
	suite.Run(t, new(TopicUsecaseSuite))
}

// Create
func (s *TopicUsecaseSuite) TestCreateTopic_Success() {
	ctx := context.Background()
	topic := entity.Topic{CategoryID: s.defaultCategoryID, AuthorID: &s.defaultAuthorID, Title: "topic title"}
	expectedTopicID := int64(1)
	category := &entity.Category{ID: s.defaultCategoryID, Title: "Existing category"}

	s.categoryRepoMock.On("GetByID", ctx, s.defaultCategoryID).Return(category, nil).Once()
	s.topicRepoMock.On("Create", ctx, topic).Return(expectedTopicID, nil).Once()

	id, err := s.usecase.Create(ctx, topic)

	s.NoError(err)
	s.Equal(expectedTopicID, id)
	s.categoryRepoMock.AssertExpectations(s.T())
	s.topicRepoMock.AssertExpectations(s.T())
}

func (s *TopicUsecaseSuite) TestCreateTopic_CategoryNotFound() {
	ctx := context.Background()
	topic := entity.Topic{CategoryID: s.defaultCategoryID, AuthorID: &s.defaultAuthorID, Title: "topic title"}
	expectedError := ErrCategoryNotFound

	s.categoryRepoMock.On("GetByID", ctx, s.defaultCategoryID).Return(nil, pgx.ErrNoRows).Once()

	id, err := s.usecase.Create(ctx, topic)

	s.Error(err)
	s.Equal(int64(0), id)
	s.ErrorIs(err, expectedError)
	s.categoryRepoMock.AssertExpectations(s.T())
	s.topicRepoMock.AssertNotCalled(s.T(), "Create", mock.Anything, mock.Anything)
}

/*
func (s *TopicUsecaseSuite) TestCreateTopic_CategoryRepoError_OtherThanNotFound() {
	ctx := context.Background()
	topic := entity.Topic{CategoryID: s.defaultCategoryID, AuthorID: &s.defaultAuthorID, Title: "topic title"}
	repoError := errors.New("some category repo error")

	s.categoryRepoMock.On("GetByID", ctx, s.defaultCategoryID).Return(nil, repoError).Once()

	id, err := s.usecase.Create(ctx, topic)

	s.Error(err)
	s.Equal(int64(0), id)
	s.ErrorIs(err, repoError)
	s.categoryRepoMock.AssertExpectations(s.T())
	s.topicRepoMock.AssertNotCalled(s.T(), "Create", mock.Anything, mock.Anything)
}*/

func (s *TopicUsecaseSuite) TestCreateTopic_RepoError() {
	ctx := context.Background()
	topic := entity.Topic{CategoryID: s.defaultCategoryID, AuthorID: &s.defaultAuthorID, Title: "topic title"}
	expectedError := errors.New("topic repository create error")
	category := &entity.Category{ID: s.defaultCategoryID, Title: "Existing category"}

	s.categoryRepoMock.On("GetByID", ctx, s.defaultCategoryID).Return(category, nil).Once()
	s.topicRepoMock.On("Create", ctx, topic).Return(int64(0), expectedError).Once()

	id, err := s.usecase.Create(ctx, topic)

	s.Error(err)
	s.Equal(int64(0), id)
	s.Contains(err.Error(), "ForumService - TopicUsecase - Create - topicRepo.Create()")
	s.ErrorIs(err, expectedError)
	s.categoryRepoMock.AssertExpectations(s.T())
	s.topicRepoMock.AssertExpectations(s.T())
}

// GetByID
func (s *TopicUsecaseSuite) TestGetByIDTopic_Success_WithAuthor() {
	ctx := context.Background()
	topicID := int64(1)
	authorID := s.defaultAuthorID
	expectedUsername := "TestUser"
	topicFromRepo := &entity.Topic{ID: topicID, AuthorID: &authorID, Title: "Test Topic", CategoryID: s.defaultCategoryID, CreatedAt: time.Now()}
	expectedTopic := &entity.Topic{ID: topicID, AuthorID: &authorID, Title: "Test Topic", CategoryID: s.defaultCategoryID, Username: expectedUsername, CreatedAt: topicFromRepo.CreatedAt}

	s.topicRepoMock.On("GetByID", ctx, topicID).Return(topicFromRepo, nil).Once()
	s.userClientMock.On("GetUsername", ctx, authorID).Return(expectedUsername, nil).Once()

	topic, err := s.usecase.GetByID(ctx, topicID)

	s.NoError(err)
	s.NotNil(topic)
	s.Equal(expectedTopic, topic)
	s.topicRepoMock.AssertExpectations(s.T())
	s.userClientMock.AssertExpectations(s.T())
}

func (s *TopicUsecaseSuite) TestGetByIDTopic_Success_AuthorNil() {
	ctx := context.Background()
	topicID := int64(1)
	topicFromRepo := &entity.Topic{ID: topicID, AuthorID: nil, Title: "Test Topic", CategoryID: s.defaultCategoryID, CreatedAt: time.Now()}
	expectedTopic := &entity.Topic{ID: topicID, AuthorID: nil, Title: "Test Topic", CategoryID: s.defaultCategoryID, Username: "Удаленный пользователь", CreatedAt: topicFromRepo.CreatedAt}

	s.topicRepoMock.On("GetByID", ctx, topicID).Return(topicFromRepo, nil).Once()

	topic, err := s.usecase.GetByID(ctx, topicID)

	s.NoError(err)
	s.NotNil(topic)
	s.Equal(expectedTopic, topic)
	s.topicRepoMock.AssertExpectations(s.T())
	s.userClientMock.AssertNotCalled(s.T(), "GetUsername", mock.Anything, mock.Anything)
}

func (s *TopicUsecaseSuite) TestGetByIDTopic_RepoError() {
	ctx := context.Background()
	topicID := int64(1)
	expectedError := errors.New("repository get by id error")

	s.topicRepoMock.On("GetByID", ctx, topicID).Return(nil, expectedError).Once()

	topic, err := s.usecase.GetByID(ctx, topicID)

	s.Error(err)
	s.Nil(topic)
	s.Contains(err.Error(), "ForumService - TopicUsecase - GetByID - repo.GetByID()")
	s.ErrorIs(err, expectedError)
	s.topicRepoMock.AssertExpectations(s.T())
	s.userClientMock.AssertNotCalled(s.T(), "GetUsername", mock.Anything, mock.Anything)
}

func (s *TopicUsecaseSuite) TestGetByIDTopic_UserClientError() {
	ctx := context.Background()
	topicID := int64(1)
	authorID := s.defaultAuthorID
	topicFromRepo := &entity.Topic{ID: topicID, AuthorID: &authorID, Title: "Test Topic"}
	expectedError := errors.New("user client GetUsername error")

	s.topicRepoMock.On("GetByID", ctx, topicID).Return(topicFromRepo, nil).Once()
	s.userClientMock.On("GetUsername", ctx, authorID).Return("", expectedError).Once()

	topic, err := s.usecase.GetByID(ctx, topicID)

	s.Error(err)
	s.Nil(topic)
	s.Contains(err.Error(), "ForumService - TopicUsecase - GetById - userClient.GetUsername()")
	s.ErrorIs(err, expectedError)
	s.topicRepoMock.AssertExpectations(s.T())
	s.userClientMock.AssertExpectations(s.T())
}

// GetByCategory
func (s *TopicUsecaseSuite) TestGetByCategory_Success() {
	ctx := context.Background()
	categoryID := s.defaultCategoryID
	authorID1 := int64(10)
	authorID2 := int64(20)
	topicsFromRepo := []entity.Topic{
		{ID: 1, CategoryID: categoryID, AuthorID: &authorID1, Title: "Topic 1", CreatedAt: time.Now()},
		{ID: 2, CategoryID: categoryID, AuthorID: &authorID2, Title: "Topic 2", CreatedAt: time.Now()},
		{ID: 3, CategoryID: categoryID, AuthorID: nil, Title: "Topic 3 - Deleted User", CreatedAt: time.Now()},
	}
	usernamesFromClient := map[int64]string{
		authorID1: "UserOne",
		authorID2: "UserTwo",
	}
	expectedTopics := []entity.Topic{
		{ID: 1, CategoryID: categoryID, AuthorID: &authorID1, Username: "UserOne", Title: "Topic 1", CreatedAt: topicsFromRepo[0].CreatedAt},
		{ID: 2, CategoryID: categoryID, AuthorID: &authorID2, Username: "UserTwo", Title: "Topic 2", CreatedAt: topicsFromRepo[1].CreatedAt},
		{ID: 3, CategoryID: categoryID, AuthorID: nil, Username: "Удаленный пользователь", Title: "Topic 3 - Deleted User", CreatedAt: topicsFromRepo[2].CreatedAt},
	}
	category := &entity.Category{ID: categoryID, Title: "Existing category"}

	s.categoryRepoMock.On("GetByID", ctx, categoryID).Return(category, nil).Once()
	s.topicRepoMock.On("GetByCategory", ctx, categoryID).Return(topicsFromRepo, nil).Once()
	s.userClientMock.On("GetUsernames", ctx, mock.MatchedBy(func(ids []int64) bool {
		s.ElementsMatch([]int64{authorID1, authorID2}, ids)
		return true
	})).Return(usernamesFromClient, nil).Once()

	topics, err := s.usecase.GetByCategory(ctx, categoryID)

	s.NoError(err)
	s.NotNil(topics)
	s.ElementsMatch(expectedTopics, topics)
	s.categoryRepoMock.AssertExpectations(s.T())
	s.topicRepoMock.AssertExpectations(s.T())
	s.userClientMock.AssertExpectations(s.T())
}

func (s *TopicUsecaseSuite) TestGetByCategory_CheckCategoryError_NotFound() {
	ctx := context.Background()
	categoryID := s.defaultCategoryID
	expectedError := ErrCategoryNotFound

	s.categoryRepoMock.On("GetByID", ctx, categoryID).Return(nil, pgx.ErrNoRows).Once()

	topics, err := s.usecase.GetByCategory(ctx, categoryID)

	s.Error(err)
	s.Nil(topics)
	s.ErrorIs(err, expectedError)
	s.categoryRepoMock.AssertExpectations(s.T())
	s.topicRepoMock.AssertNotCalled(s.T(), "GetByCategory", mock.Anything, mock.Anything)
	s.userClientMock.AssertNotCalled(s.T(), "GetUsernames", mock.Anything, mock.Anything)
}

func (s *TopicUsecaseSuite) TestGetByCategory_TopicRepoError() {
	ctx := context.Background()
	categoryID := s.defaultCategoryID
	expectedError := errors.New("topic repo GetByCategory error")
	category := &entity.Category{ID: categoryID, Title: "Existing category"}

	s.categoryRepoMock.On("GetByID", ctx, categoryID).Return(category, nil).Once()
	s.topicRepoMock.On("GetByCategory", ctx, categoryID).Return(nil, expectedError).Once()

	topics, err := s.usecase.GetByCategory(ctx, categoryID)

	s.Error(err)
	s.Nil(topics)
	s.Contains(err.Error(), "ForumService - TopicUsecase  - GetByCategory - topicRepo.GetByCategory()")
	s.ErrorIs(err, expectedError)
	s.categoryRepoMock.AssertExpectations(s.T())
	s.topicRepoMock.AssertExpectations(s.T())
	s.userClientMock.AssertNotCalled(s.T(), "GetUsernames", mock.Anything, mock.Anything)
}

func (s *TopicUsecaseSuite) TestGetByCategory_UserClientError() {
	ctx := context.Background()
	categoryID := s.defaultCategoryID
	authorID1 := int64(10)
	topicsFromRepo := []entity.Topic{
		{ID: 1, CategoryID: categoryID, AuthorID: &authorID1, Title: "Topic 1"},
	}
	expectedError := errors.New("user client GetUsernames error")
	category := &entity.Category{ID: categoryID, Title: "Existing category"}

	s.categoryRepoMock.On("GetByID", ctx, categoryID).Return(category, nil).Once()
	s.topicRepoMock.On("GetByCategory", ctx, categoryID).Return(topicsFromRepo, nil).Once()
	s.userClientMock.On("GetUsernames", ctx, []int64{authorID1}).Return(nil, expectedError).Once()

	topics, err := s.usecase.GetByCategory(ctx, categoryID)

	s.Error(err)
	s.Nil(topics)
	s.Contains(err.Error(), "ForumService - TopicUsecase  - GetByCategory - userClient.GetUsernames()")
	s.ErrorIs(err, expectedError)
	s.categoryRepoMock.AssertExpectations(s.T())
	s.topicRepoMock.AssertExpectations(s.T())
	s.userClientMock.AssertExpectations(s.T())
}

// Update
func (s *TopicUsecaseSuite) TestUpdateTopic_Success_Author() {
	ctx := context.Background()
	topicID := int64(1)
	userID := s.defaultAuthorID
	role := "user"
	title := "updated title"
	topicFromRepo := &entity.Topic{ID: topicID, AuthorID: &s.defaultAuthorID, Title: "Old title"}

	s.topicRepoMock.On("GetByID", ctx, topicID).Return(topicFromRepo, nil).Once()
	s.topicRepoMock.On("Update", ctx, topicID, title).Return(nil).Once()

	err := s.usecase.Update(ctx, topicID, userID, role, title)

	s.NoError(err)
	s.topicRepoMock.AssertExpectations(s.T())
}

func (s *TopicUsecaseSuite) TestUpdateTopic_Success_Admin() {
	ctx := context.Background()
	topicID := int64(1)
	adminID := int64(555)
	authorID := s.defaultAuthorID
	role := "admin"
	title := "Updated by Admin"
	topicFromRepo := &entity.Topic{ID: topicID, AuthorID: &authorID, Title: "Old title"}

	s.topicRepoMock.On("GetByID", ctx, topicID).Return(topicFromRepo, nil).Once()
	s.topicRepoMock.On("Update", ctx, topicID, title).Return(nil).Once()

	err := s.usecase.Update(ctx, topicID, adminID, role, title)

	s.NoError(err)
	s.topicRepoMock.AssertExpectations(s.T())
}

func (s *TopicUsecaseSuite) TestUpdateTopic_AccessDenied_NotAuthorNotAdmin() {
	ctx := context.Background()
	topicID := int64(1)
	nonAuthorID := int64(555)
	authorID := s.defaultAuthorID
	role := "user"
	title := "Attempted Update"
	topicFromRepo := &entity.Topic{ID: topicID, AuthorID: &authorID, Title: "Old title"}
	expectedError := ErrForbidden

	s.topicRepoMock.On("GetByID", ctx, topicID).Return(topicFromRepo, nil).Once()

	err := s.usecase.Update(ctx, topicID, nonAuthorID, role, title)

	s.Error(err)
	s.ErrorIs(err, expectedError)
	s.topicRepoMock.AssertCalled(s.T(), "GetByID", ctx, topicID)
	s.topicRepoMock.AssertNotCalled(s.T(), "Update", mock.Anything, mock.Anything, mock.Anything)
}

func (s *TopicUsecaseSuite) TestUpdateTopic_TopicNotFound_OnCheckAccess() {
	ctx := context.Background()
	topicID := int64(1)
	userID := s.defaultAuthorID
	role := "user"
	title := "updated title"
	expectedError := ErrTopicNotFound

	s.topicRepoMock.On("GetByID", ctx, topicID).Return(nil, pgx.ErrNoRows).Once()

	err := s.usecase.Update(ctx, topicID, userID, role, title)

	s.Error(err)
	s.ErrorIs(err, expectedError)
	s.topicRepoMock.AssertExpectations(s.T())
	s.topicRepoMock.AssertNotCalled(s.T(), "Update", mock.Anything, mock.Anything, mock.Anything)
}

func (s *TopicUsecaseSuite) TestUpdateTopic_RepoUpdateError() {
	ctx := context.Background()
	topicID := int64(1)
	userID := s.defaultAuthorID
	role := "user"
	title := "updated title"
	topicFromRepo := &entity.Topic{ID: topicID, AuthorID: &s.defaultAuthorID, Title: "Old title"}
	repoError := errors.New("repo update error")

	s.topicRepoMock.On("GetByID", ctx, topicID).Return(topicFromRepo, nil).Once()
	s.topicRepoMock.On("Update", ctx, topicID, title).Return(repoError).Once()

	err := s.usecase.Update(ctx, topicID, userID, role, title)

	s.Error(err)
	s.ErrorIs(err, repoError)
	s.Contains(err.Error(), "ForumService - TopicUsecase - Update - topicRepo.Update()")
	s.topicRepoMock.AssertExpectations(s.T())
}

// Delete
func (s *TopicUsecaseSuite) TestDeleteTopic_Success_Author() {
	ctx := context.Background()
	topicID := int64(1)
	userID := s.defaultAuthorID
	role := "user"
	topicFromRepo := &entity.Topic{ID: topicID, AuthorID: &s.defaultAuthorID}

	s.topicRepoMock.On("GetByID", ctx, topicID).Return(topicFromRepo, nil).Once()
	s.topicRepoMock.On("Delete", ctx, topicID).Return(nil).Once()

	err := s.usecase.Delete(ctx, topicID, userID, role)

	s.NoError(err)
	s.topicRepoMock.AssertExpectations(s.T())
}

func (s *TopicUsecaseSuite) TestDeleteTopic_Success_Admin() {
	ctx := context.Background()
	topicID := int64(1)
	adminID := int64(999)
	authorID := s.defaultAuthorID
	role := "admin"
	topicFromRepo := &entity.Topic{ID: topicID, AuthorID: &authorID}

	s.topicRepoMock.On("GetByID", ctx, topicID).Return(topicFromRepo, nil).Once()
	s.topicRepoMock.On("Delete", ctx, topicID).Return(nil).Once()

	err := s.usecase.Delete(ctx, topicID, adminID, role)

	s.NoError(err)
	s.topicRepoMock.AssertExpectations(s.T())
}

func (s *TopicUsecaseSuite) TestDeleteTopic_AccessDenied_NotAuthorNotAdmin() {
	ctx := context.Background()
	topicID := int64(1)
	nonAuthorID := int64(555)
	authorID := s.defaultAuthorID
	role := "user"
	topicFromRepo := &entity.Topic{ID: topicID, AuthorID: &authorID}
	expectedError := ErrForbidden

	s.topicRepoMock.On("GetByID", ctx, topicID).Return(topicFromRepo, nil).Once()

	err := s.usecase.Delete(ctx, topicID, nonAuthorID, role)

	s.Error(err)
	s.ErrorIs(err, expectedError)
	s.topicRepoMock.AssertCalled(s.T(), "GetByID", ctx, topicID)
	s.topicRepoMock.AssertNotCalled(s.T(), "Delete", mock.Anything, mock.Anything)
}

func (s *TopicUsecaseSuite) TestDeleteTopic_TopicNotFound_OnCheckAccess() {
	ctx := context.Background()
	topicID := int64(1)
	userID := s.defaultAuthorID
	role := "user"
	expectedError := ErrTopicNotFound

	s.topicRepoMock.On("GetByID", ctx, topicID).Return(nil, pgx.ErrNoRows).Once()

	err := s.usecase.Delete(ctx, topicID, userID, role)

	s.Error(err)
	s.ErrorIs(err, expectedError)
	s.topicRepoMock.AssertExpectations(s.T())
	s.topicRepoMock.AssertNotCalled(s.T(), "Delete", mock.Anything, mock.Anything)
}

func (s *TopicUsecaseSuite) TestDeleteTopic_RepoDeleteError() {
	ctx := context.Background()
	topicID := int64(1)
	userID := s.defaultAuthorID
	role := "user"
	topicFromRepo := &entity.Topic{ID: topicID, AuthorID: &s.defaultAuthorID}
	repoError := errors.New("repo delete error")

	s.topicRepoMock.On("GetByID", ctx, topicID).Return(topicFromRepo, nil).Once()
	s.topicRepoMock.On("Delete", ctx, topicID).Return(repoError).Once()

	err := s.usecase.Delete(ctx, topicID, userID, role)

	s.Error(err)
	s.ErrorIs(err, repoError)
	s.Contains(err.Error(), "ForumService - TopicUsecase - Delete - topicRepo.Delete()")
	s.topicRepoMock.AssertExpectations(s.T())
}
