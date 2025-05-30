package repo

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Van-programan/Forum_GO/internal/entity"
	"github.com/Van-programan/Forum_GO/pkg/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRefreshTokenRepository_Save(t *testing.T) {
	ctx := context.Background()
	logger := zerolog.Nop()

	mockPool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockPool.Close()

	pg := postgres.NewWithPool(mockPool)

	repo := NewRefreshTokenRepository(pg, &logger)

	token := "test-token"
	userID := int64(1)

	t.Run("Success", func(t *testing.T) {
		mockPool.ExpectExec("INSERT INTO refresh_tokens").WithArgs(token, userID).WillReturnResult(pgxmock.NewResult("INSERT", 1))

		err := repo.Save(ctx, token, userID)
		assert.NoError(t, err)
		assert.NoError(t, mockPool.ExpectationsWereMet())
	})

	t.Run("DB Error", func(t *testing.T) {
		dbErr := errors.New("database insert error")
		mockPool.ExpectExec("INSERT INTO refresh_tokens").WithArgs(token, userID).WillReturnError(dbErr)

		err := repo.Save(ctx, token, userID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "RefreshTokenRepository - Save - pg.Pool.Exec()")
		assert.ErrorIs(t, err, dbErr)
		assert.NoError(t, mockPool.ExpectationsWereMet())
	})
}

func TestRefreshTokenRepository_Delete(t *testing.T) {
	ctx := context.Background()
	logger := zerolog.Nop()

	mockPool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockPool.Close()

	pg := postgres.NewWithPool(mockPool)

	repo := NewRefreshTokenRepository(pg, &logger)

	token := "test-token"

	t.Run("Success", func(t *testing.T) {
		mockPool.ExpectExec("DELETE FROM refresh_tokens WHERE token").WithArgs(token).WillReturnResult(pgxmock.NewResult("DELETE", 1))

		err := repo.Delete(ctx, token)
		assert.NoError(t, err)
		assert.NoError(t, mockPool.ExpectationsWereMet())
	})

	t.Run("DB Error", func(t *testing.T) {
		dbErr := errors.New("database delete error")
		mockPool.ExpectExec("DELETE FROM refresh_tokens WHERE token").WithArgs(token).WillReturnError(dbErr)

		err := repo.Delete(ctx, token)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "RefreshTokenRepository - Delete - pg.Pool.Exec()")
		assert.ErrorIs(t, err, dbErr)
		assert.NoError(t, mockPool.ExpectationsWereMet())
	})
}

func TestRefreshTokenRepository_GetUserID(t *testing.T) {
	ctx := context.Background()
	logger := zerolog.Nop()

	mockPool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockPool.Close()

	pg := postgres.NewWithPool(mockPool)

	repo := NewRefreshTokenRepository(pg, &logger)

	token := "test-token"
	expectedUserID := int64(1)

	t.Run("Success", func(t *testing.T) {
		rows := pgxmock.NewRows([]string{"user_id"}).AddRow(expectedUserID)
		mockPool.ExpectQuery("SELECT user_id FROM refresh_tokens WHERE token").WithArgs(token).WillReturnRows(rows)

		userID, err := repo.GetUserID(ctx, token)
		assert.NoError(t, err)
		assert.Equal(t, expectedUserID, userID)
		assert.NoError(t, mockPool.ExpectationsWereMet())
	})

	t.Run("pgx.ErrNoRows", func(t *testing.T) {
		mockPool.ExpectQuery("SELECT user_id FROM refresh_tokens WHERE token").WithArgs(token).WillReturnError(pgx.ErrNoRows)

		userID, err := repo.GetUserID(ctx, token)
		assert.Error(t, err)
		assert.Equal(t, int64(0), userID)
		assert.Contains(t, err.Error(), "RefreshTokenRepository - GetByUserID - row.Scan()")
		assert.ErrorIs(t, err, pgx.ErrNoRows)
		assert.NoError(t, mockPool.ExpectationsWereMet())
	})

	t.Run("DB Error", func(t *testing.T) {
		dbErr := errors.New("some db error")
		mockPool.ExpectQuery("SELECT user_id FROM refresh_tokens WHERE token = \\$1").
			WithArgs(token).
			WillReturnError(dbErr)

		userID, err := repo.GetUserID(ctx, token)
		assert.Error(t, err)
		assert.Equal(t, int64(0), userID)
		assert.Contains(t, err.Error(), "RefreshTokenRepository - GetByUserID - row.Scan")
		assert.ErrorIs(t, err, dbErr)
		assert.NoError(t, mockPool.ExpectationsWereMet())
	})

}

func TestUserPostgres_Create(t *testing.T) {
	ctx := context.Background()
	logger := zerolog.Nop()

	mockPool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockPool.Close()

	pg := postgres.NewWithPool(mockPool)
	repo := NewUserRepository(pg, &logger)

	testUser := &entity.User{Username: "test", PasswordHash: "testpasswordhash"}
	expectedId := int64(1)

	t.Run("Success", func(t *testing.T) {
		row := pgxmock.NewRows([]string{"id"}).AddRow(expectedId)
		mockPool.ExpectQuery("INSERT INTO users").WithArgs(testUser.Username, testUser.PasswordHash).WillReturnRows(row)

		id, err := repo.Create(ctx, testUser)
		assert.NoError(t, err)
		assert.Equal(t, expectedId, id)
		assert.NoError(t, mockPool.ExpectationsWereMet())
	})

	t.Run("DB error", func(t *testing.T) {
		dbErr := errors.New("some db error")
		mockPool.ExpectQuery("INSERT INTO users").WithArgs(testUser.Username, testUser.PasswordHash).WillReturnError(dbErr)

		_, err := repo.Create(ctx, testUser)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "UserRepository - Create - row.Scan()")
		assert.ErrorIs(t, err, dbErr)
		assert.NoError(t, mockPool.ExpectationsWereMet())
	})
}

func TestUserPostgres_Delete(t *testing.T) {
	ctx := context.Background()
	logger := zerolog.Nop()

	mockPool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockPool.Close()

	pg := postgres.NewWithPool(mockPool)
	repo := NewUserRepository(pg, &logger)

	testID := int64(1)

	t.Run("Success", func(t *testing.T) {
		mockPool.ExpectExec("DELETE FROM users WHERE id").WithArgs(testID).WillReturnResult(pgxmock.NewResult("DELETE", 1))

		err := repo.Delete(ctx, testID)
		assert.NoError(t, err)
		assert.NoError(t, mockPool.ExpectationsWereMet())
	})

	t.Run("DB error", func(t *testing.T) {
		dbErr := errors.New("database delete error")
		mockPool.ExpectExec("DELETE FROM users WHERE id").WithArgs(testID).WillReturnError(dbErr)

		err := repo.Delete(ctx, testID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "UserRepository - Delete - pg.Pool.Exec()")
		assert.ErrorIs(t, err, dbErr)
		assert.NoError(t, mockPool.ExpectationsWereMet())
	})
}

func TestUserPostgres_GetByUsername(t *testing.T) {
	ctx := context.Background()
	logger := zerolog.Nop()

	mockPool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockPool.Close()

	pg := postgres.NewWithPool(mockPool)
	repo := NewUserRepository(pg, &logger)

	testUsername := "username"
	expectedUser := &entity.User{ID: 1, Username: "username", Role: "user", PasswordHash: "hash", CreatedAt: time.Now()}

	t.Run("Succes", func(t *testing.T) {
		row := pgxmock.NewRows([]string{"id", "username", "role", "password_hash", "created_at"}).AddRow(expectedUser.ID, expectedUser.Username, expectedUser.Role, expectedUser.PasswordHash, expectedUser.CreatedAt)
		mockPool.ExpectQuery("SELECT id, username, role, password_hash, created_at FROM users WHERE username").WithArgs(testUsername).WillReturnRows(row)

		user, err := repo.GetByUsername(ctx, testUsername)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
		assert.NoError(t, mockPool.ExpectationsWereMet())
	})

	t.Run("DB error", func(t *testing.T) {
		dbErr := errors.New("some database error")
		mockPool.ExpectQuery("SELECT id, username, role, password_hash, created_at FROM users WHERE username").WithArgs(testUsername).WillReturnError(dbErr)

		_, err := repo.GetByUsername(ctx, testUsername)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "UserRepository - GetByUsername - row.Scan()")
		assert.ErrorIs(t, err, dbErr)
		assert.NoError(t, mockPool.ExpectationsWereMet())
	})
}

func TestUserPostgres_GetByID(t *testing.T) {
	ctx := context.Background()
	logger := zerolog.Nop()

	mockPool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockPool.Close()

	pg := postgres.NewWithPool(mockPool)
	repo := NewUserRepository(pg, &logger)

	testID := int64(1)
	expectedUser := &entity.User{ID: 1, Username: "username", Role: "user", CreatedAt: time.Now()}

	t.Run("Succes", func(t *testing.T) {
		row := pgxmock.NewRows([]string{"id", "username", "role", "created_at"}).AddRow(expectedUser.ID, expectedUser.Username, expectedUser.Role, expectedUser.CreatedAt)
		mockPool.ExpectQuery("SELECT id, username, role, created_at FROM users WHERE id").WithArgs(testID).WillReturnRows(row)

		user, err := repo.GetByID(ctx, testID)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
		assert.NoError(t, mockPool.ExpectationsWereMet())
	})

	t.Run("DB error", func(t *testing.T) {
		dbErr := errors.New("some database error")
		mockPool.ExpectQuery("SELECT id, username, role, created_at FROM users WHERE id").WithArgs(testID).WillReturnError(dbErr)

		_, err := repo.GetByID(ctx, testID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "UserRepository - GetByID - row.Scan()")
		assert.ErrorIs(t, err, dbErr)
		assert.NoError(t, mockPool.ExpectationsWereMet())
	})
}

func TestUserPostgres_GetRole(t *testing.T) {
	ctx := context.Background()
	logger := zerolog.Nop()

	mockPool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockPool.Close()

	pg := postgres.NewWithPool(mockPool)
	repo := NewUserRepository(pg, &logger)

	testID := int64(1)
	expectedRole := "user"

	t.Run("Succes", func(t *testing.T) {
		row := pgxmock.NewRows([]string{"role"}).AddRow(expectedRole)
		mockPool.ExpectQuery("SELECT role FROM users WHERE id").WithArgs(testID).WillReturnRows(row)

		role, err := repo.GetRole(ctx, testID)
		assert.NoError(t, err)
		assert.Equal(t, expectedRole, role)
		assert.NoError(t, mockPool.ExpectationsWereMet())
	})

	t.Run("DB error", func(t *testing.T) {
		dbErr := errors.New("some database error")
		mockPool.ExpectQuery("SELECT role FROM users WHERE id").WithArgs(testID).WillReturnError(dbErr)

		_, err := repo.GetRole(ctx, testID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "UserRepository - GetRole - row.Scan()")
		assert.ErrorIs(t, err, dbErr)
		assert.NoError(t, mockPool.ExpectationsWereMet())
	})
}
