package usecase

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/Van-programan/Forum_GO/internal/entity"
	"github.com/Van-programan/Forum_GO/internal/repo"
	"github.com/Van-programan/Forum_GO/internal/ws"
	"github.com/Van-programan/Forum_GO/pkg/tokens"
)

type Forum interface {
	CreateTopic(ctx context.Context, title string, authorID int64) (*entity.Topic, error)
	GetTopicByID(ctx context.Context, id int64) (*entity.Topic, error)
	GetTopics(ctx context.Context) ([]*entity.Topic, error)
	UpdateTopic(ctx context.Context, id int64, newTitle string) error
	DeleteTopic(ctx context.Context, id int64) error

	CreateMessage(ctx context.Context, topicID, userID int64, content string) (*entity.Message, error)
	GetMessages(ctx context.Context, topicID int64) ([]*entity.Message, error)
	DeleteMessage(ctx context.Context, id int64) error

	RegisterWSClient(topicID string, client *ws.Client)
	BroadcastNewMessage(msg *entity.Message) error
	CanAccessTopic(userID int64, topicID string) bool
	ValidateToken(token string) (int64, error)
}

type ForumUseCase struct {
	topicRepo   repo.TopicRepository
	messageRepo repo.MessageRepository
	wsHub       *ws.Hub
	tm          tokens.TokenManager
}

func NewForumUseCase(topicRepo repo.TopicRepository, messageRepo repo.MessageRepository,
	wsHub *ws.Hub, tm tokens.TokenManager) Forum {
	return &ForumUseCase{
		topicRepo:   topicRepo,
		messageRepo: messageRepo,
		wsHub:       wsHub,
		tm:          tm,
	}
}

func (uc *ForumUseCase) CreateTopic(ctx context.Context, title string, authorID int64) (*entity.Topic, error) {
	if title == "" {
		return nil, errors.New("title cannot be empty")
	}

	topic := &entity.Topic{
		Title:    title,
		AuthorID: authorID,
	}

	err := uc.topicRepo.CreateTopic(ctx, topic)
	if err != nil {
		return nil, err
	}

	return topic, nil
}

func (uc *ForumUseCase) GetTopicByID(ctx context.Context, id int64) (*entity.Topic, error) {
	topic, err := uc.topicRepo.GetTopicByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if topic == nil {
		return nil, errors.New("topic not found")
	}

	return topic, nil
}

func (uc *ForumUseCase) GetTopics(ctx context.Context) ([]*entity.Topic, error) {
	topics, err := uc.topicRepo.GetTopics(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*entity.Topic, len(topics))
	for i := range topics {
		result[i] = &topics[i]
	}

	return result, nil
}

func (uc *ForumUseCase) UpdateTopic(ctx context.Context, id int64, newTitle string) error {
	if newTitle == "" {
		return errors.New("title cannot be empty")
	}

	return uc.topicRepo.UpdateTopic(ctx, id, newTitle)
}

func (uc *ForumUseCase) DeleteTopic(ctx context.Context, id int64) error {
	return uc.topicRepo.DeleteTopic(ctx, id)
}

func (uc *ForumUseCase) CreateMessage(ctx context.Context, topicID, userID int64, content string) (*entity.Message, error) {
	if content == "" {
		return nil, errors.New("message content cannot be empty")
	}

	message := &entity.Message{
		TopicID:   topicID,
		UserID:    userID,
		Content:   content,
		CreatedAt: time.Now(),
	}

	err := uc.messageRepo.CreateMessage(ctx, message)
	if err != nil {
		return nil, err
	}

	return message, nil
}

func (uc *ForumUseCase) GetMessages(ctx context.Context, topicID int64) ([]*entity.Message, error) {
	messages, err := uc.messageRepo.GetMessages(ctx, topicID)
	if err != nil {
		return nil, err
	}

	result := make([]*entity.Message, len(messages))
	for i := range messages {
		result[i] = &messages[i]
	}

	return result, nil
}

func (uc *ForumUseCase) DeleteMessage(ctx context.Context, id int64) error {
	return uc.messageRepo.DeleteMessage(ctx, id)
}

func (uc *ForumUseCase) RegisterWSClient(topicID string, client *ws.Client) {
	uc.wsHub.RegisterClient(client)

	uc.wsHub.BroadcastToTopic(topicID, ws.Message{
		Type: "user_joined",
		Payload: map[string]interface{}{
			"user_id": client.UserID,
		},
	})
}

func (uc *ForumUseCase) BroadcastNewMessage(msg *entity.Message) error {
	uc.wsHub.BroadcastToTopic(strconv.FormatInt(msg.TopicID, 10), ws.Message{
		Type: "new_message",
		Payload: map[string]interface{}{
			"id":         msg.ID,
			"user_id":    msg.UserID,
			"content":    msg.Content,
			"created_at": msg.CreatedAt,
		},
	})
	return nil
}

func (uc *ForumUseCase) CanAccessTopic(userID int64, topicID string) bool {
	id, err := strconv.ParseInt(topicID, 10, 64)
	if err != nil {
		return false
	}

	topic, err := uc.topicRepo.GetTopicByID(context.Background(), id)
	if err != nil {
		return false
	}

	return topic != nil
}

func (uc *ForumUseCase) ValidateToken(token string) (int64, error) {
	return uc.tm.ParseAccessToken(token)
}
