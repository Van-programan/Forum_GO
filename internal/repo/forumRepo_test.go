package repo

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Van-programan/Forum_GO/internal/entity"
	"github.com/Van-programan/Forum_GO/pkg/postgres"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCategoryRepository_Create(t *testing.T) {
	ctx := context.Background()
	logger := zerolog.Nop()

	mockPool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockPool.Close()

	pg := postgres.NewWithPool(mockPool)
	repo := NewCategoryRepository(pg, &logger)

	testCategory := entity.Category{Title: "test", Description: "test"}
	expectedID := int64(1)

	t.Run("Success", func(t *testing.T) {
		row := pgxmock.NewRows([]string{"id"}).AddRow(expectedID)
		mockPool.ExpectQuery("INSERT INTO categories").WithArgs(testCategory.Title, testCategory.Description).WillReturnRows(row)

		id, err := repo.Create(ctx, testCategory)
		assert.NoError(t, err)
		assert.Equal(t, expectedID, id)
		assert.NoError(t, mockPool.ExpectationsWereMet())
	})

	t.Run("DB error", func(t *testing.T) {
		dbErr := errors.New("some db error")
		mockPool.ExpectQuery("INSERT INTO categories").WithArgs(testCategory.Title, testCategory.Description).WillReturnError(dbErr)
		_, err := repo.Create(ctx, testCategory)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "CategoryRepository - Create - row.Scan()")
		assert.ErrorIs(t, err, dbErr)
		assert.NoError(t, mockPool.ExpectationsWereMet())
	})
}

func TestCategoryRepository_GetByID(t *testing.T) {
	ctx := context.Background()
	logger := zerolog.Nop()

	mockPool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockPool.Close()

	pg := postgres.NewWithPool(mockPool)
	repo := NewCategoryRepository(pg, &logger)

	id := int64(1)
	expectedCategory := &entity.Category{ID: id, Title: "test", Description: "test", CreatedAt: time.Now(), UpdatedAt: time.Now()}

	t.Run("Success", func(t *testing.T) {
		row := pgxmock.NewRows([]string{"id", "title", "description", "created_at", "updated_at"}).AddRow(expectedCategory.ID, expectedCategory.Title, expectedCategory.Description, expectedCategory.CreatedAt, expectedCategory.UpdatedAt)
		mockPool.ExpectQuery("SELECT id, title, description, created_at, updated_at FROM categories WHERE id").WithArgs(id).WillReturnRows(row)

		category, err := repo.GetByID(ctx, id)
		assert.NoError(t, err)
		assert.Equal(t, expectedCategory, category)
		assert.NoError(t, mockPool.ExpectationsWereMet())
	})

	t.Run("DB error", func(t *testing.T) {
		dbErr := errors.New("some db error")
		mockPool.ExpectQuery("SELECT id, title, description, created_at, updated_at FROM categories WHERE id").WithArgs(id).WillReturnError(dbErr)

		_, err := repo.GetByID(ctx, id)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "CategoryRepository - GetByID - row.Scan()")
		assert.ErrorIs(t, err, dbErr)
		assert.NoError(t, mockPool.ExpectationsWereMet())
	})
}

func TestCategoryRepository_GetAll(t *testing.T) {
	ctx := context.Background()
	logger := zerolog.Nop()

	mockPool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockPool.Close()

	pg := postgres.NewWithPool(mockPool)
	repo := NewCategoryRepository(pg, &logger)

	expectedCategories := []entity.Category{
		{ID: 1, Title: "test1", Description: "test1", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: 2, Title: "test2", Description: "test2", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	t.Run("Success", func(t *testing.T) {
		rows := pgxmock.NewRows([]string{"id", "title", "description", "created_at", "updated_at"}).AddRow(expectedCategories[0].ID, expectedCategories[0].Title, expectedCategories[0].Description, expectedCategories[0].CreatedAt, expectedCategories[0].UpdatedAt).
			AddRow(expectedCategories[1].ID, expectedCategories[1].Title, expectedCategories[1].Description, expectedCategories[1].CreatedAt, expectedCategories[1].UpdatedAt)
		mockPool.ExpectQuery("SELECT id, title, description, created_at, updated_at FROM categories ORDER BY id").WillReturnRows(rows)

		categories, err := repo.GetAll(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expectedCategories, categories)
		assert.NoError(t, mockPool.ExpectationsWereMet())
	})

	t.Run("Query error", func(t *testing.T) {
		dbErr := errors.New("query db error")
		mockPool.ExpectQuery("SELECT id, title, description, created_at, updated_at FROM categories ORDER BY id").WillReturnError(dbErr)

		_, err := repo.GetAll(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "CategoryRepository - GetCategories - pg.Pool.Query")
		assert.ErrorIs(t, err, dbErr)
		assert.NoError(t, mockPool.ExpectationsWereMet())
	})

	t.Run("Scan error", func(t *testing.T) {
		dbErr := errors.New("scan error")
		rows := pgxmock.NewRows([]string{"id", "title", "description", "created_at", "updated_at"}).AddRow(1, "test1", "test1", time.Now(), time.Now()).
			RowError(0, dbErr)

		mockPool.ExpectQuery("SELECT id, title, description, created_at, updated_at FROM categories ORDER BY id").WillReturnRows(rows)

		_, err := repo.GetAll(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "CategoryRepository - GetCategories - rows.Next() - rows.Scan()")
		assert.ErrorIs(t, err, dbErr)
		assert.NoError(t, mockPool.ExpectationsWereMet())
	})
}

func TestCategoryRepository_Update(t *testing.T) {
	ctx := context.Background()
	logger := zerolog.Nop()

	mockPool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockPool.Close()

	pg := postgres.NewWithPool(mockPool)
	repo := NewCategoryRepository(pg, &logger)

	expectedSql := "UPDATE categories SET title = COALESCE\\(\\$1, title\\), description = COALESCE\\(\\$2, description\\), updated_at = now\\(\\) WHERE id = \\$3"

	id := int64(1)
	title := "updated title"
	description := "updated description"

	t.Run("Success", func(t *testing.T) {
		mockPool.ExpectExec(expectedSql).WithArgs(title, description, id).WillReturnResult(pgxmock.NewResult("UPDATE", 1))

		err := repo.Update(ctx, id, title, description)
		assert.NoError(t, err)
		assert.NoError(t, mockPool.ExpectationsWereMet())
	})

	t.Run("DB error", func(t *testing.T) {
		dbErr := errors.New("some db error")
		mockPool.ExpectExec(expectedSql).WithArgs(title, description, id).WillReturnError(dbErr)

		err := repo.Update(ctx, id, title, description)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "CategoryRepository - Update - Exec")
		assert.ErrorIs(t, err, dbErr)
		assert.NoError(t, mockPool.ExpectationsWereMet())
	})
}

func TestCategoryRepository_Delete(t *testing.T) {
	ctx := context.Background()
	logger := zerolog.Nop()

	mockPool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockPool.Close()

	pg := postgres.NewWithPool(mockPool)
	repo := NewCategoryRepository(pg, &logger)

	id := int64(1)

	t.Run("Success", func(t *testing.T) {
		mockPool.ExpectExec("DELETE FROM categories WHERE id").WithArgs(id).WillReturnResult(pgxmock.NewResult("DELETE", 1))

		err := repo.Delete(ctx, id)
		assert.NoError(t, err)
		assert.NoError(t, mockPool.ExpectationsWereMet())
	})

	t.Run("DB error", func(t *testing.T) {
		dbErr := errors.New("some db error")
		mockPool.ExpectExec("DELETE FROM categories WHERE id").WithArgs(id).WillReturnError(dbErr)

		err := repo.Delete(ctx, id)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "CategoryRepository - Delete - pg.Pool.Exec()")
		assert.ErrorIs(t, err, dbErr)
		assert.NoError(t, mockPool.ExpectationsWereMet())
	})
}

func TestChatRepository_SaveMessage(t *testing.T) {
	ctx := context.Background()
	logger := zerolog.Nop()

	mockPool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockPool.Close()

	pg := postgres.NewWithPool(mockPool)
	repo := NewChatRepository(pg, &logger)

	testMessage := &entity.ChatMessage{
		UserID:    1,
		Username:  "user",
		Content:   "test message",
		CreatedAt: time.Now(),
	}

	expectedID := int64(1)

	t.Run("Success", func(t *testing.T) {
		row := pgxmock.NewRows([]string{"id"}).AddRow(expectedID)
		mockPool.ExpectQuery("INSERT INTO messages \\(user_id, username, content, created_at\\) VALUES\\(\\$1, \\$2, \\$3, \\$4\\) RETURNING id").WithArgs(testMessage.UserID, testMessage.Username, testMessage.Content, testMessage.CreatedAt).WillReturnRows(row)

		id, err := repo.SaveMessage(ctx, testMessage)
		assert.NoError(t, err)
		assert.Equal(t, expectedID, id)
		assert.NoError(t, mockPool.ExpectationsWereMet())
	})

	t.Run("DB error", func(t *testing.T) {
		dbErr := errors.New("some db error")
		mockPool.ExpectQuery("INSERT INTO messages \\(user_id, username, content, created_at\\) VALUES\\(\\$1, \\$2, \\$3, \\$4\\) RETURNING id").WithArgs(testMessage.UserID, testMessage.Username, testMessage.Content, testMessage.CreatedAt).WillReturnError(dbErr)

		_, err := repo.SaveMessage(ctx, testMessage)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ChatRepository - SaveMessage - row.Scan()")
		assert.ErrorIs(t, err, dbErr)
		assert.NoError(t, mockPool.ExpectationsWereMet())
	})
}

func TestChatRepository_GetMessages(t *testing.T) {
	ctx := context.Background()
	logger := zerolog.Nop()

	mockPool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockPool.Close()

	pg := postgres.NewWithPool(mockPool)
	repo := NewChatRepository(pg, &logger)

	expectedMessages := []entity.ChatMessage{
		{UserID: 1, Username: "user1", Content: "message1", CreatedAt: time.Now()},
		{UserID: 2, Username: "user2", Content: "message2", CreatedAt: time.Now()},
	}

	expectedLimit := int64(2)

	t.Run("Success", func(t *testing.T) {
		rows := pgxmock.NewRows([]string{"id", "user_id", "username", "content", "created_at"}).
			AddRow(expectedMessages[0].ID, expectedMessages[0].UserID, expectedMessages[0].Username, expectedMessages[0].Content, expectedMessages[0].CreatedAt).
			AddRow(expectedMessages[1].ID, expectedMessages[1].UserID, expectedMessages[1].Username, expectedMessages[1].Content, expectedMessages[1].CreatedAt)
		mockPool.ExpectQuery("SELECT id, user_id, username, content, created_at FROM \\(SELECT id, user_id, username, content, created_at FROM messages ORDER BY created_at DESC LIMIT \\$1\\) AS recent_mesages ORDER BY created_at ASC").WithArgs(expectedLimit).WillReturnRows(rows)

		messages, err := repo.GetMessages(ctx, expectedLimit)
		assert.NoError(t, err)
		assert.Equal(t, expectedMessages, messages)
		assert.NoError(t, mockPool.ExpectationsWereMet())
	})

	t.Run("Query error", func(t *testing.T) {
		dbErr := errors.New("some db error")
		mockPool.ExpectQuery("SELECT id, user_id, username, content, created_at FROM \\(SELECT id, user_id, username, content, created_at FROM messages ORDER BY created_at DESC LIMIT \\$1\\) AS recent_mesages ORDER BY created_at ASC").WithArgs(expectedLimit).WillReturnError(dbErr)

		_, err := repo.GetMessages(ctx, expectedLimit)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ChatRepository - GetMessages - r.pg.Pool.Query()")
		assert.ErrorIs(t, err, dbErr)
		assert.NoError(t, mockPool.ExpectationsWereMet())
	})

	t.Run("Scan error	", func(t *testing.T) {
		dbErr := errors.New("some db error")
		rows := pgxmock.NewRows([]string{"id", "user_id", "username", "content", "created_at"}).
			AddRow(expectedMessages[0].ID, expectedMessages[0].UserID, expectedMessages[0].Username, expectedMessages[0].Content, expectedMessages[0].CreatedAt).
			AddRow(expectedMessages[1].ID, expectedMessages[1].UserID, expectedMessages[1].Username, expectedMessages[1].Content, expectedMessages[1].CreatedAt).
			RowError(1, dbErr)
		mockPool.ExpectQuery("SELECT id, user_id, username, content, created_at FROM \\(SELECT id, user_id, username, content, created_at FROM messages ORDER BY created_at DESC LIMIT \\$1\\) AS recent_mesages ORDER BY created_at ASC").WithArgs(expectedLimit).WillReturnRows(rows)

		_, err := repo.GetMessages(ctx, expectedLimit)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ChatRepository - GetMessages - rows.Next()")
		assert.ErrorIs(t, err, dbErr)
		assert.NoError(t, mockPool.ExpectationsWereMet())
	})
}

func TestPostRepository_Create(t *testing.T) {
	ctx := context.Background()
	logger := zerolog.Nop()

	mockPool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockPool.Close()

	pg := postgres.NewWithPool(mockPool)
	repo := NewPostRepository(pg, &logger)
	authorID := int64(1)

	testPost := entity.Post{TopicID: 1, AuthorID: &authorID, Content: "test", ReplyTo: nil}
	expectedID := int64(1)

	t.Run("Success", func(t *testing.T) {
		row := pgxmock.NewRows([]string{"id"}).AddRow(expectedID)
		mockPool.ExpectQuery("INSERT INTO posts").WithArgs(testPost.TopicID, testPost.AuthorID, testPost.Content, testPost.ReplyTo).WillReturnRows(row)

		id, err := repo.Create(ctx, testPost)
		assert.NoError(t, err)
		assert.Equal(t, expectedID, id)
		assert.NoError(t, mockPool.ExpectationsWereMet())
	})

	t.Run("DB error", func(t *testing.T) {
		dbErr := errors.New("some db error")
		mockPool.ExpectQuery("INSERT INTO posts").WithArgs(testPost.TopicID, testPost.AuthorID, testPost.Content, testPost.ReplyTo).WillReturnError(dbErr)

		_, err := repo.Create(ctx, testPost)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "PostRepository - Create - row.Scan()")
		assert.ErrorIs(t, err, dbErr)
		assert.NoError(t, mockPool.ExpectationsWereMet())
	})
}

func TestPostRepository_GetByID(t *testing.T) {
	ctx := context.Background()
	logger := zerolog.Nop()

	mockPool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockPool.Close()

	pg := postgres.NewWithPool(mockPool)
	repo := NewPostRepository(pg, &logger)

	id := int64(1)
	authorID := int64(1)

	expectedPost := &entity.Post{ID: 1, AuthorID: &authorID, Content: "test", ReplyTo: nil, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	t.Run("Success", func(t *testing.T) {
		row := pgxmock.NewRows([]string{"id", "content", "author_id", "reply_to", "created_at", "updated_at"}).AddRow(expectedPost.ID, expectedPost.Content, expectedPost.AuthorID, expectedPost.ReplyTo, expectedPost.CreatedAt, expectedPost.UpdatedAt)
		mockPool.ExpectQuery("SELECT id, content, author_id, reply_to, created_at, updated_at FROM posts WHERE id").WithArgs(id).WillReturnRows(row)

		post, err := repo.GetByID(ctx, id)
		assert.NoError(t, err)
		assert.Equal(t, expectedPost, post)
		assert.NoError(t, mockPool.ExpectationsWereMet())
	})

	t.Run("DB error", func(t *testing.T) {
		dbErr := errors.New("some db error")
		mockPool.ExpectQuery("SELECT id, content, author_id, reply_to, created_at, updated_at FROM posts WHERE id").WithArgs(id).WillReturnError(dbErr)

		_, err := repo.GetByID(ctx, id)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "PostRepository - GetByID - row.Scan()")
		assert.ErrorIs(t, err, dbErr)
		assert.NoError(t, mockPool.ExpectationsWereMet())
	})
}

func TestPostRepository_GetByTopic(t *testing.T) {
	ctx := context.Background()
	logger := zerolog.Nop()

	mockPool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockPool.Close()

	pg := postgres.NewWithPool(mockPool)
	repo := NewPostRepository(pg, &logger)

	topicID := int64(1)
	authorID := int64(1)
	expectedPosts := []entity.Post{
		{ID: 1, TopicID: topicID, Content: "test", AuthorID: &authorID, ReplyTo: nil, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: 2, TopicID: topicID, Content: "test2", AuthorID: &authorID, ReplyTo: nil, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	t.Run("Success", func(t *testing.T) {
		rows := pgxmock.NewRows([]string{"id", "topic_id", "content", "author_id", "reply_to", "created_at", "updated_at"}).AddRow(expectedPosts[0].ID, expectedPosts[0].TopicID, expectedPosts[0].Content, expectedPosts[0].AuthorID, expectedPosts[0].ReplyTo, expectedPosts[0].CreatedAt, expectedPosts[0].UpdatedAt).
			AddRow(expectedPosts[1].ID, expectedPosts[1].TopicID, expectedPosts[1].Content, expectedPosts[1].AuthorID, expectedPosts[1].ReplyTo, expectedPosts[1].CreatedAt, expectedPosts[1].UpdatedAt)
		mockPool.ExpectQuery("SELECT id, topic_id, content, author_id, reply_to, created_at, updated_at FROM posts WHERE topic_id = \\$1 ORDER BY created_at").WithArgs(topicID).WillReturnRows(rows)

		posts, err := repo.GetByTopic(ctx, topicID)
		assert.NoError(t, err)
		assert.Equal(t, expectedPosts, posts)
		assert.NoError(t, mockPool.ExpectationsWereMet())
	})

	t.Run("Query error", func(t *testing.T) {
		dbErr := errors.New("some db error")
		mockPool.ExpectQuery("SELECT id, topic_id, content, author_id, reply_to, created_at, updated_at FROM posts WHERE topic_id = \\$1 ORDER BY created_at").WithArgs(topicID).WillReturnError(dbErr)

		_, err := repo.GetByTopic(ctx, topicID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "PostRepository - GetByTopic - pg.Pool.Query")
		assert.ErrorIs(t, err, dbErr)
		assert.NoError(t, mockPool.ExpectationsWereMet())
	})

	t.Run("Scan error", func(t *testing.T) {
		dbErr := errors.New("scan db error")
		rows := pgxmock.NewRows([]string{"id", "topic_id", "content", "author_id", "reply_to", "created_at", "updated_at"}).AddRow(expectedPosts[0].ID, expectedPosts[0].TopicID, expectedPosts[0].Content, expectedPosts[0].AuthorID, expectedPosts[0].ReplyTo, expectedPosts[0].CreatedAt, expectedPosts[0].UpdatedAt).
			RowError(0, dbErr)
		mockPool.ExpectQuery("SELECT id, topic_id, content, author_id, reply_to, created_at, updated_at FROM posts WHERE topic_id = \\$1 ORDER BY created_at").WithArgs(topicID).WillReturnRows(rows)

		_, err := repo.GetByTopic(ctx, topicID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "PostRepository - GetByTopic - rows.Next() - rows.Scan()")
		assert.ErrorIs(t, err, dbErr)
		assert.NoError(t, mockPool.ExpectationsWereMet())
	})
}

func TestPostRepository_Update(t *testing.T) {
	ctx := context.Background()
	logger := zerolog.Nop()

	mockPool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockPool.Close()

	pg := postgres.NewWithPool(mockPool)
	repo := NewPostRepository(pg, &logger)

	expectedSql := "UPDATE posts SET content = \\$1, updated_at = now\\(\\) WHERE id = \\$2"

	id := int64(1)
	content := "updated content"

	t.Run("Success", func(t *testing.T) {
		mockPool.ExpectExec(expectedSql).WithArgs(content, id).WillReturnResult(pgxmock.NewResult("UPDATE", 1))

		err := repo.Update(ctx, id, content)
		assert.NoError(t, err)
		assert.NoError(t, mockPool.ExpectationsWereMet())
	})

	t.Run("DB error", func(t *testing.T) {
		dbErr := errors.New("some db error")
		mockPool.ExpectExec(expectedSql).WithArgs(content, id).WillReturnError(dbErr)

		err := repo.Update(ctx, id, content)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "PostRepository - Update - Exec")
		assert.ErrorIs(t, err, dbErr)
		assert.NoError(t, mockPool.ExpectationsWereMet())
	})
}

func TestPostRepository_Delete(t *testing.T) {
	ctx := context.Background()
	logger := zerolog.Nop()

	mockPool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockPool.Close()

	pg := postgres.NewWithPool(mockPool)
	repo := NewPostRepository(pg, &logger)

	id := int64(1)

	t.Run("Success", func(t *testing.T) {
		mockPool.ExpectExec("DELETE FROM posts WHERE id").WithArgs(id).WillReturnResult(pgxmock.NewResult("DELETE", 1))

		err := repo.Delete(ctx, id)
		assert.NoError(t, err)
		assert.NoError(t, mockPool.ExpectationsWereMet())
	})

	t.Run("DB error", func(t *testing.T) {
		dbErr := errors.New("some db error")
		mockPool.ExpectExec("DELETE FROM posts WHERE id").WithArgs(id).WillReturnError(dbErr)

		err := repo.Delete(ctx, id)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "PostRepository - Delete - pg.Pool.Exec()")
		assert.ErrorIs(t, err, dbErr)
		assert.NoError(t, mockPool.ExpectationsWereMet())
	})
}

func TestTopicRepository_Create(t *testing.T) {
	ctx := context.Background()
	logger := zerolog.Nop()

	mockPool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockPool.Close()

	pg := postgres.NewWithPool(mockPool)
	repo := NewTopicRepository(pg, &logger)
	authorID := int64(1)

	testTopic := entity.Topic{CategoryID: 1, Title: "test", AuthorID: &authorID}
	expectedID := int64(1)

	t.Run("Success", func(t *testing.T) {
		row := pgxmock.NewRows([]string{"id"}).AddRow(expectedID)
		mockPool.ExpectQuery("INSERT INTO topics").WithArgs(testTopic.CategoryID, testTopic.Title, testTopic.AuthorID).WillReturnRows(row)

		id, err := repo.Create(ctx, testTopic)
		assert.NoError(t, err)
		assert.Equal(t, expectedID, id)
		assert.NoError(t, mockPool.ExpectationsWereMet())
	})

	t.Run("DB error", func(t *testing.T) {
		dbErr := errors.New("some db error")
		mockPool.ExpectQuery("INSERT INTO topics").WithArgs(testTopic.CategoryID, testTopic.Title, testTopic.AuthorID).WillReturnError(dbErr)

		_, err := repo.Create(ctx, testTopic)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "TopicRepository - Create - row.Scan()")
		assert.ErrorIs(t, err, dbErr)
		assert.NoError(t, mockPool.ExpectationsWereMet())
	})
}

func TestTopicRepository_GetByID(t *testing.T) {
	ctx := context.Background()
	logger := zerolog.Nop()

	mockPool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockPool.Close()

	pg := postgres.NewWithPool(mockPool)
	repo := NewTopicRepository(pg, &logger)

	id := int64(1)
	authorID := int64(1)

	expectedTopic := &entity.Topic{ID: id, CategoryID: 1, Title: "test", AuthorID: &authorID, CreatedAt: time.Now(), UpdatedAt: time.Now()}

	t.Run("Success", func(t *testing.T) {
		row := pgxmock.NewRows([]string{"id", "category_id", "title", "author_id", "created_at", "updated_at"}).AddRow(expectedTopic.ID, expectedTopic.CategoryID, expectedTopic.Title, expectedTopic.AuthorID, expectedTopic.CreatedAt, expectedTopic.UpdatedAt)
		mockPool.ExpectQuery("SELECT id, category_id, title, author_id, created_at, updated_at FROM topics WHERE id").WithArgs(id).WillReturnRows(row)

		topic, err := repo.GetByID(ctx, id)
		assert.NoError(t, err)
		assert.Equal(t, expectedTopic, topic)
		assert.NoError(t, mockPool.ExpectationsWereMet())
	})

	t.Run("DB error", func(t *testing.T) {
		dbErr := errors.New("some db error")
		mockPool.ExpectQuery("SELECT id, category_id, title, author_id, created_at, updated_at FROM topics WHERE id").WithArgs(id).WillReturnError(dbErr)

		_, err := repo.GetByID(ctx, id)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "TopicRepository - GetByID - row.Scan()")
		assert.ErrorIs(t, err, dbErr)
		assert.NoError(t, mockPool.ExpectationsWereMet())
	})
}

func TestTopicRepository_GetByCategory(t *testing.T) {
	ctx := context.Background()
	logger := zerolog.Nop()

	mockPool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockPool.Close()

	pg := postgres.NewWithPool(mockPool)
	repo := NewTopicRepository(pg, &logger)

	categoryID := int64(1)
	authorID := int64(1)
	expectedTopics := []entity.Topic{
		{ID: 1, CategoryID: categoryID, Title: "test", AuthorID: &authorID, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: 2, CategoryID: categoryID, Title: "test2", AuthorID: &authorID, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	t.Run("Success", func(t *testing.T) {
		rows := pgxmock.NewRows([]string{"id", "category_id", "title", "author_id", "created_at", "updated_at"}).AddRow(expectedTopics[0].ID, expectedTopics[0].CategoryID, expectedTopics[0].Title, expectedTopics[0].AuthorID, expectedTopics[0].CreatedAt, expectedTopics[0].UpdatedAt).
			AddRow(expectedTopics[1].ID, expectedTopics[1].CategoryID, expectedTopics[1].Title, expectedTopics[1].AuthorID, expectedTopics[1].CreatedAt, expectedTopics[1].UpdatedAt)
		mockPool.ExpectQuery("SELECT id, category_id, title, author_id, created_at, updated_at FROM topics WHERE category_id").WithArgs(categoryID).WillReturnRows(rows)

		topics, err := repo.GetByCategory(ctx, categoryID)
		assert.NoError(t, err)
		assert.Equal(t, expectedTopics, topics)
		assert.NoError(t, mockPool.ExpectationsWereMet())
	})

	t.Run("Query error", func(t *testing.T) {
		dbErr := errors.New("some db error")
		mockPool.ExpectQuery("SELECT id, category_id, title, author_id, created_at, updated_at FROM topics WHERE category_id").WithArgs(categoryID).WillReturnError(dbErr)

		_, err := repo.GetByCategory(ctx, categoryID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "TopicRepository - GetByCategory - pg.Pool.Query")
		assert.ErrorIs(t, err, dbErr)
		assert.NoError(t, mockPool.ExpectationsWereMet())

	})

	t.Run("Scan error", func(t *testing.T) {
		dbErr := errors.New("scan db error")
		rows := pgxmock.NewRows([]string{"id", "category_id", "title", "author_id", "created_at", "updated_at"}).AddRow(expectedTopics[0].ID, expectedTopics[0].CategoryID, expectedTopics[0].Title, expectedTopics[0].AuthorID, expectedTopics[0].CreatedAt, expectedTopics[0].UpdatedAt).
			RowError(0, dbErr)
		mockPool.ExpectQuery("SELECT id, category_id, title, author_id, created_at, updated_at FROM topics WHERE category_id").WithArgs(categoryID).WillReturnRows(rows)

		_, err := repo.GetByCategory(ctx, categoryID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "TopicRepository - GetByCategory - rows.Next() - rows.Scan()")
		assert.ErrorIs(t, err, dbErr)
		assert.NoError(t, mockPool.ExpectationsWereMet())
	})
}

func TestTopicRepository_Update(t *testing.T) {
	ctx := context.Background()
	logger := zerolog.Nop()

	mockPool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockPool.Close()

	pg := postgres.NewWithPool(mockPool)
	repo := NewTopicRepository(pg, &logger)

	expectedSql := "UPDATE topics SET title = \\$1, updated_at = now\\(\\) WHERE id = \\$2"

	id := int64(1)
	title := "updated title"

	t.Run("Success", func(t *testing.T) {
		mockPool.ExpectExec(expectedSql).WithArgs(title, id).WillReturnResult(pgxmock.NewResult("UPDATE", 1))

		err := repo.Update(ctx, id, title)
		assert.NoError(t, err)
		assert.NoError(t, mockPool.ExpectationsWereMet())
	})

	t.Run("DB error", func(t *testing.T) {
		dbErr := errors.New("some db error")
		mockPool.ExpectExec(expectedSql).WithArgs(title, id).WillReturnError(dbErr)

		err := repo.Update(ctx, id, title)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "TopicRepository - Update - Exec")
		assert.ErrorIs(t, err, dbErr)
		assert.NoError(t, mockPool.ExpectationsWereMet())
	})
}

func TestTopicRepository_Delete(t *testing.T) {
	ctx := context.Background()
	logger := zerolog.Nop()

	mockPool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockPool.Close()

	pg := postgres.NewWithPool(mockPool)
	repo := NewTopicRepository(pg, &logger)

	id := int64(1)

	t.Run("Success", func(t *testing.T) {
		mockPool.ExpectExec("DELETE FROM topics WHERE id").WithArgs(id).WillReturnResult(pgxmock.NewResult("DELETE", 1))

		err := repo.Delete(ctx, id)
		assert.NoError(t, err)
		assert.NoError(t, mockPool.ExpectationsWereMet())
	})

	t.Run("DB error", func(t *testing.T) {
		dbErr := errors.New("some db error")
		mockPool.ExpectExec("DELETE FROM topics WHERE id").WithArgs(id).WillReturnError(dbErr)

		err := repo.Delete(ctx, id)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "TopicRepository - Delete - pg.Pool.Exec()")
		assert.ErrorIs(t, err, dbErr)
		assert.NoError(t, mockPool.ExpectationsWereMet())
	})
}
