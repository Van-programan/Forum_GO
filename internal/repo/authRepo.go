package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/Van-programan/Forum_GO/internal/entity"
	"github.com/Van-programan/Forum_GO/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *entity.User) error
	GetUsers(ctx context.Context) ([]entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	UpdateUser(ctx context.Context, user *entity.User) error
	DeleteUser(ctx context.Context, id int64) error
}
type SessionRepository interface {
	CreateSession(ctx context.Context, session *entity.Session) error
	GetSessionByRefreshToken(ctx context.Context, token string) (*entity.Session, error)
	UpdateSession(ctx context.Context, session *entity.Session) error
	DeleteSession(ctx context.Context, id int64) error
}

type sessionRepo struct {
	pg *postgres.Postgres
}

type userRepo struct {
	pg *postgres.Postgres
}

func NewAuthRepository(pg *postgres.Postgres) UserRepository {
	return &userRepo{pg: pg}
}

func NewSessionRepository(pg *postgres.Postgres) SessionRepository {
	return &sessionRepo{pg: pg}
}

func (r *userRepo) CreateUser(ctx context.Context, user *entity.User) error {
	query := `
		INSERT INTO users (username, email, password_hash, registered_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id`

	err := r.pg.Pool.QueryRow(ctx, query,
		user.Username,
		user.Email,
		user.Password,
		user.RegisteredAt,
	).Scan(&user.ID)

	if err != nil {
		return fmt.Errorf("authRepo.CreateUser: %w", err)
	}
	return nil
}

func (r *userRepo) GetUsers(ctx context.Context) ([]entity.User, error) {
	query := `
        SELECT id, username, email, registered_at
        FROM users
        ORDER BY id`

	rows, err := r.pg.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("authRepo.GetUsers: %w", err)
	}
	defer rows.Close()

	var users []entity.User
	for rows.Next() {
		var user entity.User
		if err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.RegisteredAt,
		); err != nil {
			return nil, fmt.Errorf("authRepo.GetUsers: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("authRepo.GetUsers: %w", err)
	}

	return users, nil
}

func (r *userRepo) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `
		SELECT id, username, email, password_hash, registered_at
		FROM users
		WHERE email = $1`

	var user entity.User
	err := r.pg.Pool.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.RegisteredAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("authRepo.GetUserByEmail: %w", err)
	}
	return &user, nil
}

func (r *userRepo) UpdateUser(ctx context.Context, user *entity.User) error {
	query := `
		UPDATE users
		SET username = $1, email = $2, password_hash = $3
		WHERE id = $4`

	_, err := r.pg.Pool.Exec(ctx, query,
		user.Username,
		user.Email,
		user.Password,
		user.ID,
	)

	if err != nil {
		return fmt.Errorf("authRepo.UpdateUser: %w", err)
	}
	return nil
}

func (r *userRepo) DeleteUser(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.pg.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("authRepo.DeleteUser: %w", err)
	}
	return nil
}

func (r *sessionRepo) CreateSession(ctx context.Context, session *entity.Session) error {
	query := `
		INSERT INTO sessions (
			user_id, refresh_token, refresh_token_expires_at
		) VALUES ($1, $2, $3)
		RETURNING id`

	err := r.pg.Pool.QueryRow(ctx, query,
		session.UserID,
		session.RefreshToken,
		session.ExpiresAtRefreshToken,
	).Scan(&session.ID)

	if err != nil {
		return fmt.Errorf("sessionRepo.CreateSession: %w", err)
	}
	return nil
}

func (r *sessionRepo) GetSessionByRefreshToken(ctx context.Context, token string) (*entity.Session, error) {
	query := `
		SELECT id, user_id, refresh_token, refresh_token_expires_at
		FROM sessions
		WHERE refresh_token = $1`

	var session entity.Session
	err := r.pg.Pool.QueryRow(ctx, query, token).Scan(
		&session.ID,
		&session.UserID,
		&session.RefreshToken,
		&session.ExpiresAtRefreshToken,
	)

	if err != nil {
		return nil, fmt.Errorf("sessionRepo.GetSessionByRefreshToken: %w", err)
	}
	return &session, nil
}

func (r *sessionRepo) UpdateSession(ctx context.Context, session *entity.Session) error {
	query := `
        UPDATE sessions 
        SET 
            refresh_token = $1,
            refresh_token_expires_at = $2
        WHERE id = $3`

	_, err := r.pg.Pool.Exec(ctx, query,
		session.RefreshToken,
		session.ExpiresAtRefreshToken,
		session.ID,
	)

	if err != nil {
		return fmt.Errorf("sessionRepo.UpdateSession: %w", err)
	}
	return nil
}

func (r *sessionRepo) DeleteSession(ctx context.Context, id int64) error {
	query := `DELETE FROM sessions WHERE id = $1`
	_, err := r.pg.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("sessionRepo.DeleteSession: %w", err)
	}
	return nil
}
