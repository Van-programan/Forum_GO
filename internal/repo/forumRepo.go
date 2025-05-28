package repo

import (
	"context"
	"fmt"

	"github.com/Van-programan/Forum_GO/internal/entity"
	"github.com/Van-programan/Forum_GO/pkg/postgres"
	"github.com/rs/zerolog"
)

type (
	CategoryRepository interface {
		Create(context.Context, entity.Category) (int64, error)
		GetByID(context.Context, int64) (*entity.Category, error)
		GetAll(context.Context) ([]entity.Category, error)
		Update(ctx context.Context, id int64, title, description string) error
		Delete(ctx context.Context, id int64) error
	}

	TopicRepository interface {
		Create(context.Context, entity.Topic) (int64, error)
		GetByID(context.Context, int64) (*entity.Topic, error)
		GetByCategory(ct context.Context, categoryID int64) ([]entity.Topic, error)
		Update(ctx context.Context, id int64, title string) error
		Delete(ctx context.Context, id int64) error
	}

	PostRepository interface {
		Create(context.Context, entity.Post) (int64, error)
		GetByID(context.Context, int64) (*entity.Post, error)
		GetByTopic(ctx context.Context, topicID int64) ([]entity.Post, error)
		Update(ctx context.Context, id int64, content string) error
		Delete(ctx context.Context, id int64) error
	}

	ChatRepository interface {
		SaveMessage(ctx context.Context, message *entity.ChatMessage) (int64, error)
		GetMessages(ctx context.Context, limit int64) ([]entity.ChatMessage, error)
	}
)

const (
	createOpCategory  = "CategoryRepository.Create"
	getByIdOpCategory = "CategoryRepository.GetById"
	getAllOpCategory  = "CategoryRepository.GetAll"
	deleteOpCategory  = "CategoryRepository.Delete"
	updateOpCategory  = "CategoryRepository.Update"
)

const (
	createPostOp  = "PoptRepository.Create"
	getByIdPostOp = "PostRepository.GetById"
	getByTopicOp  = "PostRepository.GetAll"
	deletePostOp  = "PostRepository.Delete"
	updatePostOp  = "PostRepository.Update"
)

const (
	createTopicOp   = "TopicRepository.Create"
	getByIdTopicOp  = "TopicRepository.GetById"
	getByCategoryOp = "TopicRepository.GetAll"
	deleteTopicOp   = "TopicRepository.Delete"
	updateTopicOp   = "TopicRepository.Update"
	countTopicOp    = "TopicRepository.CountByCategory"
)

type topicRepository struct {
	pg  *postgres.Postgres
	log *zerolog.Logger
}

type categoryRepository struct {
	pg  *postgres.Postgres
	log *zerolog.Logger
}

type chatRepository struct {
	pg  *postgres.Postgres
	log *zerolog.Logger
}

type postRepository struct {
	pg  *postgres.Postgres
	log *zerolog.Logger
}

func NewChatRepository(pg *postgres.Postgres, log *zerolog.Logger) ChatRepository {
	return &chatRepository{pg, log}
}

func NewCategoryRepository(pg *postgres.Postgres, log *zerolog.Logger) CategoryRepository {
	return &categoryRepository{pg, log}
}

func NewPostRepository(pg *postgres.Postgres, log *zerolog.Logger) PostRepository {
	return &postRepository{pg, log}
}

func NewTopicRepository(pg *postgres.Postgres, log *zerolog.Logger) TopicRepository {
	return &topicRepository{pg, log}
}

func (r *categoryRepository) Create(ctx context.Context, category entity.Category) (int64, error) {
	row := r.pg.Pool.QueryRow(ctx, "INSERT INTO categories (title, description) VALUES($1, $2) RETURNING id", category.Title, category.Description)

	var id int64
	if err := row.Scan(&id); err != nil {
		r.log.Error().Err(err).Str("op", createOpCategory).Any("category", category).Msg("Failed to insert category")
		return 0, fmt.Errorf("CategoryRepository - Create - row.Scan(): %w", err)
	}

	return id, nil
}

func (r *categoryRepository) GetByID(ctx context.Context, id int64) (*entity.Category, error) {
	row := r.pg.Pool.QueryRow(ctx, "SELECT id, title, description, created_at, updated_at FROM categories WHERE id = $1", id)

	var c entity.Category
	if err := row.Scan(&c.ID, &c.Title, &c.Description, &c.CreatedAt, &c.UpdatedAt); err != nil {
		r.log.Error().Err(err).Str("op", getByIdOpCategory).Int64("id", id).Msg("Failed to get category")
		return nil, fmt.Errorf("CategoryRepository - GetByID - row.Scan(): %w", err)
	}

	return &c, nil
}

func (r *categoryRepository) GetAll(ctx context.Context) ([]entity.Category, error) {
	rows, err := r.pg.Pool.Query(ctx, "SELECT id, title, description, created_at, updated_at FROM categories ORDER BY id")
	if err != nil {
		r.log.Error().Err(err).Str("op", getAllOpCategory).Msg("Failed to get categories")
		return nil, fmt.Errorf("CategoryRepository - GetCategories - pg.Pool.Query: %w", err)
	}
	defer rows.Close()

	var categories []entity.Category
	var c entity.Category
	for rows.Next() {
		err := rows.Scan(&c.ID, &c.Title, &c.Description, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			r.log.Error().Err(err).Str("op", getAllOpCategory).Msg("Failed to scan category")
			return nil, fmt.Errorf("CategoryRepository - GetCategories - rows.Next() - rows.Scan(): %w", err)
		}
		categories = append(categories, c)
	}

	return categories, nil
}

func (r *categoryRepository) Update(ctx context.Context, id int64, title, description string) error {
	_, err := r.pg.Pool.Exec(ctx, `
	UPDATE categories
	SET
		title = COALESCE($1, title),
		description = COALESCE($2, description),
		updated_at = now()
	WHERE id = $3
	`, title, description, id)

	if err != nil {
		r.log.Error().Err(err).Str("op", updateOpCategory).Msg("Failed to update category")
		return fmt.Errorf("CategoryRepository - Update - Exec: %w", err)
	}

	return nil
}

func (r *categoryRepository) Delete(ctx context.Context, id int64) error {
	if _, err := r.pg.Pool.Exec(ctx, `DELETE FROM categories WHERE id = $1`, id); err != nil {
		r.log.Error().Err(err).Str("op", deleteOpCategory).Msg("Failed to delete category")
		return fmt.Errorf("CategoryRepository - Delete - pg.Pool.Exec(): %w", err)
	}
	return nil
}

func (r *chatRepository) SaveMessage(ctx context.Context, message *entity.ChatMessage) (int64, error) {
	row := r.pg.Pool.QueryRow(ctx, "INSERT INTO messages (user_id, username, content, created_at) VALUES($1, $2, $3, $4) RETURNING id", message.UserID, message.Username, message.Content, message.CreatedAt)

	var id int64
	if err := row.Scan(&id); err != nil {
		r.log.Error().Err(err).Str("op", "ChatRepository.SaveMessage").Any("message", message).Msg("Failed to insert message")
		return 0, fmt.Errorf("ChatRepository - SaveMessage - row.Scan(): %w", err)
	}

	return id, nil
}

func (r *chatRepository) GetMessages(ctx context.Context, limit int64) ([]entity.ChatMessage, error) {
	rows, err := r.pg.Pool.Query(ctx, "SELECT id, user_id, username, content, created_at FROM (SELECT id, user_id, username, content, created_at FROM messages ORDER BY created_at DESC LIMIT $1) AS recent_mesages ORDER BY created_at ASC", limit)
	if err != nil {
		r.log.Error().Err(err).Str("op", "ChatRepository.GetMessages").Msg("Failed to get messages")
		return nil, fmt.Errorf("ChatRepository - GetMessages - r.pg.Pool.Query(): %w", err)
	}
	defer rows.Close()

	var messages []entity.ChatMessage
	for rows.Next() {
		var message entity.ChatMessage
		if err := rows.Scan(&message.ID, &message.UserID, &message.Username, &message.Content, &message.CreatedAt); err != nil {
			r.log.Error().Err(err).Str("op", "ChatRepository.GetMessages").Msg("Failed to scan message")
			return nil, fmt.Errorf("ChatRepository - GetMessages - rows.Next(): %w", err)
		}
		messages = append(messages, message)
	}

	return messages, nil
}

func (r *postRepository) Create(ctx context.Context, post entity.Post) (int64, error) {
	row := r.pg.Pool.QueryRow(ctx, "INSERT INTO posts (topic_id, author_id, content, reply_to) VALUES($1, $2, $3, $4) RETURNING id", post.TopicID, post.AuthorID, post.Content, post.ReplyTo)

	var id int64
	if err := row.Scan(&id); err != nil {
		r.log.Error().Err(err).Str("op", createPostOp).Any("post", post).Msg("Failed to insert post")
		return 0, fmt.Errorf("PostRepository - Create - row.Scan(): %w", err)
	}

	return id, nil
}

func (r *postRepository) GetByID(ctx context.Context, id int64) (*entity.Post, error) {
	row := r.pg.Pool.QueryRow(ctx, "SELECT id, content, author_id, reply_to, created_at, updated_at FROM posts WHERE id = $1", id)

	var p entity.Post
	if err := row.Scan(&p.ID, &p.Content, &p.AuthorID, &p.ReplyTo, &p.CreatedAt, &p.UpdatedAt); err != nil {
		r.log.Error().Err(err).Str("op", getByIdPostOp).Int64("id", id).Msg("Failed to get post")
		return nil, fmt.Errorf("PostRepository - GetByID - row.Scan(): %w", err)
	}

	return &p, nil
}

func (r *postRepository) GetByTopic(ctx context.Context, topicID int64) ([]entity.Post, error) {
	rows, err := r.pg.Pool.Query(ctx, "SELECT id, topic_id, content, author_id, reply_to, created_at, updated_at FROM posts WHERE topic_id = $1 ORDER BY created_at", topicID)
	if err != nil {
		r.log.Error().Err(err).Str("op", getByTopicOp).Int64("topic_id", topicID).Msg("Failed to get posts")
		return nil, fmt.Errorf("PostRepository - GetByTopic - pg.Pool.Query: %w", err)
	}
	defer rows.Close()

	var posts []entity.Post
	var p entity.Post
	for rows.Next() {
		err := rows.Scan(&p.ID, &p.TopicID, &p.Content, &p.AuthorID, &p.ReplyTo, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			r.log.Error().Err(err).Str("op", getByTopicOp).Int64("topic_id", topicID).Msg("Failed to scan post")
			return nil, fmt.Errorf("PostRepository - GetByTopic - rows.Next() - rows.Scan(): %w", err)
		}
		posts = append(posts, p)
	}

	return posts, nil
}

func (r *postRepository) Update(ctx context.Context, id int64, content string) error {
	if _, err := r.pg.Pool.Exec(ctx, "UPDATE posts SET content = $1, updated_at = now() WHERE id = $2", content, id); err != nil {
		r.log.Error().Err(err).Str("op", getByTopicOp).Int64("id", id).Msg("Failed to update post")
		return fmt.Errorf("PostRepository - Update - Exec: %w", err)
	}
	return nil
}

func (r *postRepository) Delete(ctx context.Context, id int64) error {
	if _, err := r.pg.Pool.Exec(ctx, `DELETE FROM posts WHERE id = $1`, id); err != nil {
		return fmt.Errorf("PostRepository - Delete - pg.Pool.Exec(): %w", err)
	}
	return nil
}

func (r *topicRepository) Create(ctx context.Context, topic entity.Topic) (int64, error) {
	row := r.pg.Pool.QueryRow(ctx, "INSERT INTO topics (category_id, title, author_id) VALUES($1, $2, $3) RETURNING id", topic.CategoryID, topic.Title, topic.AuthorID)
	var id int64
	if err := row.Scan(&id); err != nil {
		r.log.Error().Err(err).Str("op", createTopicOp).Any("topic", topic).Msg("Failed to insert topic")
		return 0, fmt.Errorf("TopicRepository - Create - row.Scan(): %w", err)
	}

	return id, nil
}

func (r *topicRepository) GetByID(ctx context.Context, id int64) (*entity.Topic, error) {
	row := r.pg.Pool.QueryRow(ctx, "SELECT id, category_id, title, author_id, created_at, updated_at FROM topics WHERE id = $1", id)

	var t entity.Topic
	if err := row.Scan(&t.ID, &t.CategoryID, &t.Title, &t.AuthorID, &t.CreatedAt, &t.UpdatedAt); err != nil {
		r.log.Error().Err(err).Str("op", getByIdTopicOp).Int64("id", id).Msg("Failed to get topic")
		return nil, fmt.Errorf("TopicRepository - GetByID - row.Scan(): %w", err)
	}

	return &t, nil
}

func (r *topicRepository) GetByCategory(ctx context.Context, categoryID int64) ([]entity.Topic, error) {
	rows, err := r.pg.Pool.Query(ctx, "SELECT id, category_id, title, author_id, created_at, updated_at FROM topics WHERE category_id = $1 ORDER BY created_at DESC", categoryID)
	if err != nil {
		r.log.Error().Err(err).Str("op", getByCategoryOp).Int64("category_id", categoryID).Msg("Failed to get topics")
		return nil, fmt.Errorf("TopicRepository - GetByCategory - pg.Pool.Query: %w", err)
	}
	defer rows.Close()

	var topics []entity.Topic
	var t entity.Topic
	for rows.Next() {
		err := rows.Scan(&t.ID, &t.CategoryID, &t.Title, &t.AuthorID, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			r.log.Error().Err(err).Str("op", getByCategoryOp).Int64("category_id", categoryID).Msg("Failed to scan topic")
			return nil, fmt.Errorf("TopicRepository - GetByCategory - rows.Next() - rows.Scan(): %w", err)
		}
		topics = append(topics, t)
	}

	return topics, nil
}

func (r *topicRepository) Update(ctx context.Context, id int64, title string) error {
	if _, err := r.pg.Pool.Exec(ctx, "UPDATE topics SET title = $1, updated_at = now() WHERE id = $2", title, id); err != nil {
		r.log.Error().Err(err).Str("op", updateTopicOp).Int64("id", id).Msg("Failed to update topic")
		return fmt.Errorf("TopicRepository - Update - Exec: %w", err)
	}
	return nil
}

func (r *topicRepository) Delete(ctx context.Context, id int64) error {
	if _, err := r.pg.Pool.Exec(ctx, `DELETE FROM topics WHERE id = $1`, id); err != nil {
		r.log.Error().Err(err).Str("op", deleteTopicOp).Int64("id", id).Msg("Failed to delete topic")
		return fmt.Errorf("TopicRepository - Delete - pg.Pool.Exec(): %w", err)
	}
	return nil
}
