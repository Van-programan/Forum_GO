package repo

import (
	"context"
	"fmt"

	"github.com/Van-programan/Forum_GO/internal/entity"
	"github.com/Van-programan/Forum_GO/pkg/postgres"
	"github.com/rs/zerolog"
)

type (
	UserRepository interface {
		Create(ctx context.Context, user *entity.User) (int64, error)
		Delete(ctx context.Context, id int64) error
		GetByUsername(ctx context.Context, username string) (*entity.User, error)
		GetByID(ctx context.Context, id int64) (*entity.User, error)
		GetRole(ctx context.Context, id int64) (string, error)
	}

	RefreshTokenRepository interface {
		Save(ctx context.Context, token string, userID int64) error
		Delete(ctx context.Context, token string) error
		GetUserID(ctx context.Context, token string) (int64, error)
	}
)

const (
	createOp        = "UserRepository.Create"
	deleteOp        = "UserRepository.Delete"
	getByUsernameOp = "UserRepository.GetByUsername"
	getByIDOp       = "UserRepository.GetByID"
	getRoleOp       = "UserRepository.GetRole"
)

const (
	saveOp        = "RefreshTokenRepository.Save"
	deleteTokenOp = "RefreshTokenRepository.Delete"
	getUserIDOp   = "RefreshTokenRepository.GetUserID"
	isActiveOp    = "RefreshTokenRepository.IsActive"
)

type userRepository struct {
	pg  *postgres.Postgres
	log *zerolog.Logger
}

type refreshTokenRepository struct {
	pg  *postgres.Postgres
	log *zerolog.Logger
}

func NewUserRepository(pg *postgres.Postgres, log *zerolog.Logger) UserRepository {
	return &userRepository{pg, log}
}

func NewRefreshTokenRepository(pg *postgres.Postgres, log *zerolog.Logger) RefreshTokenRepository {
	return &refreshTokenRepository{pg, log}
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) (int64, error) {
	row := r.pg.Pool.QueryRow(ctx,
		"INSERT INTO users (username, password_hash) VALUES($1, $2) RETURNING id",
		user.Username, string(user.PasswordHash))

	var id int64
	if err := row.Scan(&id); err != nil {
		r.log.Error().Err(err).Str("op", createOp).Str("username", user.Username).Msg("Failed to scan user ID after insert")
		return 0, fmt.Errorf("UserRepository - Create - row.Scan(): %w", err)
	}

	return id, nil
}

func (r *userRepository) Delete(ctx context.Context, id int64) error {
	if _, err := r.pg.Pool.Exec(ctx, "DELETE FROM users WHERE id = $1", id); err != nil {
		r.log.Error().Err(err).Str("op", deleteOp).Int64("id", id).Msg("Failed to delete user")
		return fmt.Errorf("UserRepository - Delete - pg.Pool.Exec(): %w", err)
	}
	return nil
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	row := r.pg.Pool.QueryRow(ctx, "SELECT id, username, role, password_hash, created_at FROM users WHERE username = $1", username)

	var u entity.User
	if err := row.Scan(&u.ID, &u.Username, &u.Role, &u.PasswordHash, &u.CreatedAt); err != nil {
		r.log.Error().Err(err).Str("op", getByUsernameOp).Str("username", username).Msg("Failed to scan user")
		return nil, fmt.Errorf("UserRepository - GetByUsername - row.Scan(): %w", err)
	}

	return &u, nil
}

func (r *userRepository) GetByID(ctx context.Context, id int64) (*entity.User, error) {
	row := r.pg.Pool.QueryRow(ctx, "SELECT id, username, role, created_at FROM users WHERE id = $1", id)

	var u entity.User
	if err := row.Scan(&u.ID, &u.Username, &u.Role, &u.CreatedAt); err != nil {
		r.log.Error().Err(err).Str("op", getByIDOp).Int64("id", id).Msg("Failed to scan user")
		return nil, fmt.Errorf("UserRepository - GetByID - row.Scan(): %w", err)
	}

	return &u, nil
}

func (r *userRepository) GetRole(ctx context.Context, id int64) (string, error) {
	row := r.pg.Pool.QueryRow(ctx, "SELECT role FROM users WHERE id = $1", id)

	var role string
	if err := row.Scan(&role); err != nil {
		r.log.Error().Err(err).Str("op", getRoleOp).Int64("id", id).Msg("Failed to get role")
		return "", fmt.Errorf("UserRepository - GetRole - row.Scan(): %w", err)
	}

	return role, nil
}

func (r *refreshTokenRepository) Save(ctx context.Context, token string, userID int64) error {
	if _, err := r.pg.Pool.Exec(ctx, "INSERT INTO refresh_tokens (token, user_id) VALUES($1, $2)", token, userID); err != nil {
		r.log.Error().Err(err).Str("op", saveOp).Str("token", token).Int64("userID", userID).Msg("Failed to save refresh token")
		return fmt.Errorf("RefreshTokenRepository - Save - pg.Pool.Exec(): %w", err)
	}
	r.pg.Pool.Config()
	return nil
}

func (r *refreshTokenRepository) Delete(ctx context.Context, token string) error {
	if _, err := r.pg.Pool.Exec(ctx, `DELETE FROM refresh_tokens WHERE token = $1`, token); err != nil {
		r.log.Error().Err(err).Str("op", deleteTokenOp).Str("token", token).Msg("Failed to delete refresh token")
		return fmt.Errorf("RefreshTokenRepository - Delete - pg.Pool.Exec(): %w", err)
	}
	return nil
}

func (r *refreshTokenRepository) GetUserID(ctx context.Context, token string) (int64, error) {
	row := r.pg.Pool.QueryRow(ctx, "SELECT user_id FROM refresh_tokens WHERE token = $1", token)

	var id int64
	if err := row.Scan(&id); err != nil {
		r.log.Error().Err(err).Str("op", getUserIDOp).Str("token", token).Msg("Failed to get user ID")
		return 0, fmt.Errorf("RefreshTokenRepository - GetByUserID - row.Scan(): %w", err)
	}

	return id, nil
}
