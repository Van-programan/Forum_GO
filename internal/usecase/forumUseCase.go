package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Van-programan/Forum_GO/internal/client"
	"github.com/Van-programan/Forum_GO/internal/entity"
	"github.com/Van-programan/Forum_GO/internal/repo"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
)

type (
	CategoryUsecase interface {
		Create(context.Context, entity.Category) (int64, error)
		GetByID(ctx context.Context, id int64) (*entity.Category, error)
		GetAll(context.Context) ([]entity.Category, error)
		Update(ctx context.Context, id int64, title, description string) error
		Delete(ctx context.Context, id int64) error
	}

	PostUsecase interface {
		Create(context.Context, entity.Post) (int64, error)
		GetByTopic(ctx context.Context, topicID int64) ([]entity.Post, error)
		Update(ctx context.Context, postID int64, userID int64, role string, content string) error
		Delete(ctx context.Context, postID int64, userID int64, role string) error
	}

	TopicUsecase interface {
		Create(context.Context, entity.Topic) (int64, error)
		GetByID(ctx context.Context, id int64) (*entity.Topic, error)
		GetByCategory(ct context.Context, categoryID int64) ([]entity.Topic, error)
		Update(ctx context.Context, topicID int64, userID int64, role string, title string) error
		Delete(ctx context.Context, topicID int64, userID int64, role string) error
	}

	ChatUsecase interface {
		GetMessageHistory(ctx context.Context, limit int64) ([]entity.ChatMessage, error)
		SaveMessage(ctx context.Context, userID int64, username string, content string) (*entity.ChatMessage, error)
	}
)

var (
	ErrCategoryNotFound = errors.New("category not found")
	ErrTopicNotFound    = errors.New("topic not found")
	ErrPostNotFound     = errors.New("post not found")
	ErrForbidden        = errors.New("forbidden")
)

const (
	createOp  = "CategoryUsecase.Create"
	getByIdOp = "CategoryUsecase.GetByID"
	getAllOp  = "CategoryUsecase.GetAll"
	deleteOp  = "CategoryUsecase.Delete"
	updateOp  = "CategoryUsecase.Update"
)

const (
	createPostOp = "PostUsecase.Create"
	getByTopicOp = "PostUsecase.GetByTopic"
	deletePostOp = "PostUsecase.Delete"
	updatePostOp = "PostUsecase.Update"
)

const (
	createTopicOp   = "TopicUsecase.Create"
	getByCategoryOp = "TopicUsecase.GetAll"
	deleteTopicOp   = "TopicUsecase.Delete"
	updateTopicOp   = "TopicUsecase.Update"
	getByIdTopicOp  = "TopicUsecase.GetByID"
)

type postUsecase struct {
	postRepo   repo.PostRepository
	topicRepo  repo.TopicRepository
	userClient client.UserClient
	log        *zerolog.Logger
}

type categoryUsecase struct {
	repo repo.CategoryRepository
	log  *zerolog.Logger
}

type chatUsecase struct {
	chatRepo repo.ChatRepository
	log      *zerolog.Logger
}

type topicUsecase struct {
	topicRepo    repo.TopicRepository
	categoryRepo repo.CategoryRepository
	userClient   client.UserClient
	log          *zerolog.Logger
}

func NewTopicUsecase(topicRepo repo.TopicRepository, categoryRepo repo.CategoryRepository, userClient client.UserClient, log *zerolog.Logger) TopicUsecase {
	return &topicUsecase{
		topicRepo:    topicRepo,
		categoryRepo: categoryRepo,
		userClient:   userClient,
		log:          log,
	}
}

func NewChatUsecase(chatRepo repo.ChatRepository, log *zerolog.Logger) ChatUsecase {
	return &chatUsecase{
		chatRepo: chatRepo,
		log:      log,
	}
}

func NewCategoryUsecase(repo repo.CategoryRepository, log *zerolog.Logger) CategoryUsecase {
	return &categoryUsecase{repo, log}
}

func NewPostUsecase(postRepo repo.PostRepository, topicRepo repo.TopicRepository, userClient client.UserClient, log *zerolog.Logger) PostUsecase {
	return &postUsecase{postRepo: postRepo, topicRepo: topicRepo, userClient: userClient, log: log}
}

func (u *categoryUsecase) Create(ctx context.Context, category entity.Category) (int64, error) {
	id, err := u.repo.Create(ctx, category)
	if err != nil {
		u.log.Error().Err(err).Str("op", createOp).Any("category", category).Msg("Failed to create category in repository")
		return 0, fmt.Errorf("ForumService - CategoryUsecase - Create - repo.Create(): %w", err)
	}
	u.log.Info().Str("op", createOp).Any("category", category).Msg("Category created successfully")
	return id, nil
}

func (u *categoryUsecase) GetByID(ctx context.Context, id int64) (*entity.Category, error) {
	category, err := u.repo.GetByID(ctx, id)
	if err != nil {
		u.log.Error().Err(err).Str("op", getByIdOp).Int64("id", id).Msg("Failed to get category in repository")
		return nil, fmt.Errorf("ForumService - CategoryUsecase - GetByID - repo.GetByID(): %w", err)
	}

	u.log.Info().Str("op", getByIdOp).Int64("id", id).Msg("Category taken successfully")
	return category, nil
}

func (u *categoryUsecase) GetAll(ctx context.Context) ([]entity.Category, error) {
	categories, err := u.repo.GetAll(ctx)
	if err != nil {
		u.log.Error().Err(err).Str("op", getAllOp).Msg("Failed to get categories in repository")
		return nil, fmt.Errorf("ForumService - CategoryUsecase - GetAll - repo.GetAll(): %w", err)
	}
	u.log.Info().Str("op", getAllOp).Msg("All categories succesfully taken")
	return categories, nil
}

func (u *categoryUsecase) Update(ctx context.Context, id int64, title, description string) error {
	if err := u.repo.Update(ctx, id, title, description); err != nil {
		u.log.Error().Err(err).Str("op", updateOp).Int64("id", id).Msg("Failed to update category in repository")
		return fmt.Errorf("ForumService - CategoryUsecase - Update - repo.Update(): %w", err)
	}
	u.log.Info().Str("op", updateOp).Int64("id", id).Msg("Category updated successfully")
	return nil
}

func (u *categoryUsecase) Delete(ctx context.Context, id int64) error {
	if err := u.repo.Delete(ctx, id); err != nil {
		u.log.Error().Err(err).Str("op", deleteOp).Int64("id", id).Msg("Failed to delete category in repository")
		return fmt.Errorf("ForumService - CategoryUsecase - Delete - repo.Delete(): %w", err)
	}
	u.log.Info().Str("op", deleteOp).Int64("id", id).Msg("Category deleted successfully")
	return nil
}

func (u *chatUsecase) GetMessageHistory(ctx context.Context, limit int64) ([]entity.ChatMessage, error) {
	messages, err := u.chatRepo.GetMessages(ctx, limit)
	if err != nil {
		u.log.Error().Err(err).Str("op", "ChatUsecase.GetMessageHistory").Msg("Failed to get message history")
		return nil, fmt.Errorf("ChatUsecase - GetMessageHistory - u.chatRepo.GetMessages(): %w", err)
	}
	u.log.Info().Int64("limit", limit).Int64("total_messages", int64(len(messages))).Msg("Message history retrieved")
	return messages, nil
}

func (u *chatUsecase) SaveMessage(ctx context.Context, userID int64, username string, content string) (*entity.ChatMessage, error) {
	message := &entity.ChatMessage{
		UserID:    userID,
		Username:  username,
		Content:   content,
		CreatedAt: time.Now(),
	}

	if _, err := u.chatRepo.SaveMessage(ctx, message); err != nil {
		u.log.Error().Err(err).Str("op", "ChatUsecase.SaveMessage").Msg("Failed to save message")
		return nil, fmt.Errorf("ChatUsecase - SaveMessage - u.chatRepo.SaveMessage(): %w", err)
	}

	u.log.Info().Int64("user_id", message.UserID).Str("username", message.Username).Msg("Message saved successfully")
	return message, nil
}

func (u *postUsecase) Create(ctx context.Context, post entity.Post) (int64, error) {
	if err := u.checkTopic(ctx, post.TopicID); err != nil {
		u.log.Error().Err(err).Str("op", createPostOp).Int64("topic_id", post.TopicID).Msg("Topic not found")
		return 0, err
	}

	id, err := u.postRepo.Create(ctx, post)
	if err != nil {
		u.log.Error().Err(err).Str("op", createPostOp).Any("post", post).Msg("Failed to create post in repository")
		return 0, fmt.Errorf("ForumService - PostUsecase - Create - postRepo.Create(): %w", err)
	}

	u.log.Info().Str("op", createPostOp).Any("post", post).Msg("Post successfully created")
	return id, nil
}

func (u *postUsecase) GetByTopic(ctx context.Context, topicID int64) ([]entity.Post, error) {
	if err := u.checkTopic(ctx, topicID); err != nil {
		u.log.Error().Err(err).Str("op", getByTopicOp).Int64("topic_id", topicID).Msg("Topic not found")
		return nil, err
	}
	posts, err := u.postRepo.GetByTopic(ctx, topicID)
	if err != nil {
		u.log.Error().Err(err).Str("op", getByTopicOp).Int64("topic_id", topicID).Msg("Failed to get posts")
		return nil, fmt.Errorf("ForumService - PostUsecase - GetByTopic - postRepo.GetByTopic(): %w", err)
	}

	var authorIDs []int64
	authorIDSet := make(map[int64]bool)
	for i := range posts {
		if posts[i].AuthorID != nil {
			if _, exists := authorIDSet[*posts[i].AuthorID]; !exists {
				authorIDs = append(authorIDs, *posts[i].AuthorID)
				authorIDSet[*posts[i].AuthorID] = true
			}
		}
	}

	usernames, err := u.userClient.GetUsernames(ctx, authorIDs)
	if err != nil {
		return nil, fmt.Errorf("ForumService - TopicUsecase  - GetByCategory - userClient.GetUsernames(): %w", err)
	}

	for i := range posts {
		if posts[i].AuthorID == nil {
			posts[i].Username = "Удаленный пользователь"
			continue
		}

		if username, exists := usernames[*posts[i].AuthorID]; exists {
			posts[i].Username = username
		} else {
			posts[i].Username = "Удаленный пользователь"
		}
	}

	u.log.Info().Str("op", getByTopicOp).Int64("topic_id", topicID).Msg("Posts by topic succesfully taken")
	return posts, nil
}

func (u *postUsecase) Update(ctx context.Context, postID int64, userID int64, role string, content string) error {
	if err := u.checkAccess(ctx, postID, userID, role); err != nil {
		u.log.Warn().Err(err).Str("op", updatePostOp).Int64("post_id", postID).Int64("user_id", userID).Msg("Access denied")
		return err
	}

	if err := u.postRepo.Update(ctx, postID, content); err != nil {
		u.log.Error().Err(err).Str("op", updatePostOp).Int64("post_id", postID).Int64("user_id", userID).Msg("Failed to update post in repository")
		return fmt.Errorf("ForumService - PostUsecase - Update - postRepo.Update(): %w", err)
	}

	u.log.Info().Str("op", updatePostOp).Int64("post_id", postID).Msg("Post updated successfully")
	return nil
}

func (u *postUsecase) Delete(ctx context.Context, postID int64, userID int64, role string) error {
	fmt.Printf("USER_ID: %d ,  POST_ID: %d , ROLE: %s", userID, postID, role)
	if err := u.checkAccess(ctx, postID, userID, role); err != nil {
		u.log.Warn().Err(err).Str("op", deletePostOp).Int64("post_id", postID).Int64("user_id", userID).Msg("Access denied")
		return err
	}

	if err := u.postRepo.Delete(ctx, postID); err != nil {
		u.log.Error().Err(err).Str("op", deletePostOp).Int64("post_id", postID).Int64("user_id", userID).Msg("Failed to delete post in repository")
		return fmt.Errorf("ForumService - PostUsecase - Delete - postRepo.delete(): %w", err)
	}

	u.log.Info().Str("op", updatePostOp).Int64("post_id", postID).Msg("Post deleted successfully")
	return nil
}

func (u *postUsecase) checkTopic(ctx context.Context, topicID int64) error {
	if _, err := u.topicRepo.GetByID(ctx, topicID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("ForumService - PostUsecase - checkTopic - topicRepo.GetByID()")
		}
		return fmt.Errorf("ForumService - PostUsecase - checkTopic - topicRepo.GetByID(): %w", err)
	}

	return nil
}

func (u *postUsecase) checkAccess(ctx context.Context, postID int64, userID int64, role string) error {
	post, err := u.postRepo.GetByID(ctx, postID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("ForumService - PostUsecase - checkAccess - postRepo.GetByID()")
		}
		return fmt.Errorf("ForumService - PostUsecase - checkAccess  - postRepo.GetByID(): %w", err)
	}

	if role == "admin" {
		return nil
	}

	if post.AuthorID == nil || (*post.AuthorID != userID) {
		return fmt.Errorf("ForumService - PostUsecase - checkAccess  - postRepo.Update()")
	}

	return nil
}

func (u *topicUsecase) Create(ctx context.Context, topic entity.Topic) (int64, error) {
	if err := u.checkCategory(ctx, topic.CategoryID); err != nil {
		u.log.Error().Err(err).Str("op", createTopicOp).Int64("category_id", topic.CategoryID).Msg("Category not found")
		return 0, err
	}

	id, err := u.topicRepo.Create(ctx, topic)
	if err != nil {
		u.log.Error().Err(err).Str("op", createTopicOp).Any("topic", topic).Msg("Failed to create topic in repository")
		return 0, fmt.Errorf("ForumService - TopicUsecase - Create - topicRepo.Create(): %w", err)
	}

	u.log.Info().Str("op", createTopicOp).Any("topic", topic).Msg("Topic created successfully")
	return id, nil
}

func (u *topicUsecase) GetByID(ctx context.Context, id int64) (*entity.Topic, error) {
	topic, err := u.topicRepo.GetByID(ctx, id)
	if err != nil {
		u.log.Error().Err(err).Str("op", getByIdTopicOp).Int64("id", id).Msg("Failed to get topic in repository")
		return nil, fmt.Errorf("ForumService - TopicUsecase - GetByID - repo.GetByID(): %w", err)
	}

	var username string

	if topic.AuthorID == nil {
		username = "Удаленный пользователь"
	} else {
		if uname, err := u.userClient.GetUsername(ctx, *topic.AuthorID); err != nil {
			return nil, fmt.Errorf("ForumService - TopicUsecase - GetById - userClient.GetUsername(): %w", err)
		} else {
			username = uname
		}
	}

	topic.Username = username

	u.log.Info().Str("op", getByIdTopicOp).Int64("id", id).Msg("Topic taken successfully")
	return topic, nil
}

func (u *topicUsecase) GetByCategory(ctx context.Context, categoryID int64) ([]entity.Topic, error) {
	if err := u.checkCategory(ctx, categoryID); err != nil {
		u.log.Error().Err(err).Str("op", getByCategoryOp).Int64("category_id", categoryID).Msg("Category not found")
		return nil, err
	}

	topics, err := u.topicRepo.GetByCategory(ctx, categoryID)
	if err != nil {
		u.log.Error().Err(err).Str("op", getByCategoryOp).Int64("category_id", categoryID).Msg("Failed to get topics in repository")
		return nil, fmt.Errorf("ForumService - TopicUsecase  - GetByCategory - topicRepo.GetByCategory(): %w", err)
	}

	var authorIDs []int64
	authorIDSet := make(map[int64]bool)
	for i := range topics {
		if topics[i].AuthorID != nil {
			if _, exists := authorIDSet[*topics[i].AuthorID]; !exists {
				authorIDs = append(authorIDs, *topics[i].AuthorID)
				authorIDSet[*topics[i].AuthorID] = true
			}
		}
	}

	usernames, err := u.userClient.GetUsernames(ctx, authorIDs)
	if err != nil {
		return nil, fmt.Errorf("ForumService - TopicUsecase  - GetByCategory - userClient.GetUsernames(): %w", err)
	}

	for i := range topics {
		if topics[i].AuthorID == nil {
			topics[i].Username = "Удаленный пользователь"
			continue
		}

		if username, exists := usernames[*topics[i].AuthorID]; exists {
			topics[i].Username = username
		} else {
			topics[i].Username = "Удаленный пользователь"
		}
	}

	u.log.Info().Str("op", getByCategoryOp).Int64("category_id", categoryID).Msg("Topics by category succesfully taken")
	return topics, nil
}

func (u *topicUsecase) Update(ctx context.Context, topicID int64, userID int64, role string, title string) error {
	if err := u.checkAccess(ctx, topicID, userID, role); err != nil {
		u.log.Warn().Err(err).Str("op", updateTopicOp).Int64("topic_id", topicID).Int64("user_id", userID).Msg("Access denied")
		return err
	}

	if err := u.topicRepo.Update(ctx, topicID, title); err != nil {
		u.log.Error().Err(err).Str("op", updateTopicOp).Int64("topic_id", topicID).Int64("user_id", userID).Msg("Failed to update topic in repository")
		return fmt.Errorf("ForumService - TopicUsecase - Update - topicRepo.Update(): %w", err)
	}

	u.log.Info().Str("op", updateTopicOp).Int64("topic_id", topicID).Msg("Topic updated successfully")
	return nil
}

func (u *topicUsecase) Delete(ctx context.Context, topicID int64, userID int64, role string) error {
	if err := u.checkAccess(ctx, topicID, userID, role); err != nil {
		u.log.Warn().Err(err).Str("op", deleteTopicOp).Int64("topic_id", topicID).Int64("user_id", userID).Msg("Access denied")
		return err
	}

	if err := u.topicRepo.Delete(ctx, topicID); err != nil {
		u.log.Error().Err(err).Str("op", deleteTopicOp).Int64("topic_id", topicID).Int64("user_id", userID).Msg("Access denied")
		return fmt.Errorf("ForumService - TopicUsecase - Delete - topicRepo.Delete(): %w", err)
	}

	u.log.Info().Str("op", deleteTopicOp).Int64("topic_id", topicID).Msg("Topic deleted successfully")
	return nil
}

func (u *topicUsecase) checkAccess(ctx context.Context, topicID int64, userID int64, role string) error {
	post, err := u.topicRepo.GetByID(ctx, topicID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("ForumService - TopicUsecase - checkAccess - topicRepo.GetByID(): %w", ErrTopicNotFound)
		}
		return fmt.Errorf("ForumService - TopicUsecase - checkAccess  - topicRepo.GetByID(): %w", err)
	}

	if role == "admin" {
		return nil
	}

	if post.AuthorID == nil || (*post.AuthorID != userID) {
		return fmt.Errorf("ForumService - TopicUsecase - checkAccess  - topicRepo.Update(): %w", ErrForbidden)
	}

	return nil
}

func (u *topicUsecase) checkCategory(ctx context.Context, categoryID int64) error {
	fmt.Println("checkCategory", categoryID)
	if _, err := u.categoryRepo.GetByID(ctx, categoryID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("ForumService - TopicUsecase - checkCategory - categoryRepo.GetByID(): %w", ErrCategoryNotFound)
		}
		return fmt.Errorf("ForumService - TopicUsecase - checkCategory - categoryRepo.GetByID(): %w", err)
	}

	return nil
}
