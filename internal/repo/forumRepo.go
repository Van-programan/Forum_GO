package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/Van-programan/Forum_GO/internal/entity"
	"github.com/Van-programan/Forum_GO/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

type MessageRepository interface {
	CreateMessage(ctx context.Context, message *entity.Message) error
	GetMessages(ctx context.Context, topic_id int64) ([]entity.Message, error)
	DeleteMessage(ctx context.Context, id int64) error
}
type TopicRepository interface {
	CreateTopic(ctx context.Context, topic *entity.Topic) error
	GetTopicByID(ctx context.Context, id int64) (*entity.Topic, error)
	GetTopics(ctx context.Context) ([]entity.Topic, error)
	UpdateTopic(ctx context.Context, id int64, newTitle string) error
	DeleteTopic(ctx context.Context, id int64) error
}

type topicRepo struct {
	pg *postgres.Postgres
}

type messageRepo struct {
	pg *postgres.Postgres
}

func NewMessageRepository(pg *postgres.Postgres) MessageRepository {
	return &messageRepo{pg: pg}
}

func NewTopicRepository(pg *postgres.Postgres) TopicRepository {
	return &topicRepo{pg: pg}
}

func (r *messageRepo) CreateMessage(ctx context.Context, message *entity.Message) error {
	query := `
		INSERT INTO messages (topic_id, user_id, content, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id`

	err := r.pg.Pool.QueryRow(ctx, query,
		message.TopicID,
		message.UserID,
		message.Content,
		message.CreatedAt,
	).Scan(&message.ID)

	if err != nil {
		return fmt.Errorf("messageRepo.CreateMessage: %w", err)
	}
	return nil
}

func (r *messageRepo) GetMessages(ctx context.Context, topic_id int64) ([]entity.Message, error) {
	query := `
		SELECT id, topic_id, user_id, content, created_at
		FROM messages
		WHERE topic_id = $1
		ORDER BY created_at DESC`

	rows, err := r.pg.Pool.Query(ctx, query, topic_id)
	if err != nil {
		return nil, fmt.Errorf("messageRepo.GetMessages: %w", err)
	}
	defer rows.Close()

	var messages []entity.Message
	for rows.Next() {
		var msg entity.Message
		if err := rows.Scan(
			&msg.ID,
			&msg.TopicID,
			&msg.UserID,
			&msg.Content,
			&msg.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("messageRepo.GetMessages: %w", err)
		}
		messages = append(messages, msg)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("messageRepo.GetMessages: %w", err)
	}

	return messages, nil
}

func (r *messageRepo) DeleteMessage(ctx context.Context, id int64) error {
	query := `DELETE FROM messages WHERE id = $1`
	_, err := r.pg.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("messageRepo.DeleteMessage: %w", err)
	}
	return nil
}

func (r *topicRepo) CreateTopic(ctx context.Context, topic *entity.Topic) error {
	query := `
		INSERT INTO topics (title, author_id)
		VALUES ($1, $2)
		RETURNING id, created_at`

	err := r.pg.Pool.QueryRow(ctx, query,
		topic.Title,
		topic.AuthorID,
	).Scan(
		&topic.ID,
		&topic.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("topicRepo.CreateTopic: %w", err)
	}
	return nil
}

func (r *topicRepo) GetTopicByID(ctx context.Context, id int64) (*entity.Topic, error) {
	query := `
		SELECT id, title, author_id, created_at
		FROM topics
		WHERE id = $1`

	var topic entity.Topic
	err := r.pg.Pool.QueryRow(ctx, query, id).Scan(
		&topic.ID,
		&topic.Title,
		&topic.AuthorID,
		&topic.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("topicRepo.GetTopicByID: %w", err)
	}
	return &topic, nil
}
func (r *topicRepo) GetTopics(ctx context.Context) ([]entity.Topic, error) {
	query := `
		SELECT id, title, author_id, created_at
		FROM topics
		ORDER BY created_at DESC`

	rows, err := r.pg.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("topicRepo.GetTopics: %w", err)
	}
	defer rows.Close()

	var topics []entity.Topic
	for rows.Next() {
		var topic entity.Topic
		if err := rows.Scan(
			&topic.ID,
			&topic.Title,
			&topic.AuthorID,
			&topic.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("topicRepo.GetTopics: %w", err)
		}
		topics = append(topics, topic)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("topicRepo.GetTopics: %w", err)
	}

	return topics, nil
}

func (r *topicRepo) UpdateTopic(ctx context.Context, id int64, newTitle string) error {
	query := `
		UPDATE topics
		SET title = $1
		WHERE id = $2`

	_, err := r.pg.Pool.Exec(ctx, query, newTitle, id)
	if err != nil {
		return fmt.Errorf("topicRepo.UpdateTopic: %w", err)
	}
	return nil
}

func (r *topicRepo) DeleteTopic(ctx context.Context, id int64) error {
	query := `DELETE FROM topics WHERE id = $1`
	_, err := r.pg.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("topicRepo.DeleteTopic: %w", err)
	}
	return nil
}
