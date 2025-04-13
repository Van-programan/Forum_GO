package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/Van-programan/Forum_GO/internal/entity"
	"github.com/Van-programan/Forum_GO/internal/repo"
	"github.com/Van-programan/Forum_GO/pkg/tokens"
)

type Auth interface {
	Register(ctx context.Context, username, email, password string) (*entity.User, error)
	Login(ctx context.Context, email, password string) (*entity.User, string, string, error)
	RefreshToken(ctx context.Context, refreshToken string) (*entity.User, string, error)
	Logout(ctx context.Context, refreshToken string) error
	GetUser(ctx context.Context, id int64) (*entity.User, error)
}

type AuthUseCase struct {
	userRepo     repo.UserRepository
	sessionRepo  repo.SessionRepository
	tokenManager tokens.TokenManager
}

func NewAuthUseCase(userRepo repo.UserRepository, sessionRepo repo.SessionRepository,
	tokenManager tokens.TokenManager) Auth {
	return &AuthUseCase{
		userRepo:     userRepo,
		sessionRepo:  sessionRepo,
		tokenManager: tokenManager,
	}
}

func (uc *AuthUseCase) Register(ctx context.Context, username, email, password string) (*entity.User, error) {
	existingUser, err := uc.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	hashedPassword, err := uc.tokenManager.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		Username:     username,
		Email:        email,
		Password:     hashedPassword,
		RegisteredAt: time.Now(),
	}

	err = uc.userRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *AuthUseCase) Login(ctx context.Context, email, password string) (*entity.User, string, string, error) {
	user, err := uc.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, "", "", err
	}
	if user == nil {
		return nil, "", "", errors.New("user not found")
	}

	if !uc.tokenManager.CheckPasswordHash(password, user.Password) {
		return nil, "", "", errors.New("invalid password")
	}

	accessToken, err := uc.tokenManager.GenerateAccessToken(user.ID)
	if err != nil {
		return nil, "", "", err
	}

	refreshToken, err := uc.tokenManager.GenerateRefreshToken()
	if err != nil {
		return nil, "", "", err
	}

	session := &entity.Session{
		UserID:                user.ID,
		RefreshToken:          refreshToken,
		ExpiresAtRefreshToken: time.Now().Add(tokens.RefreshTokenTTL),
	}

	err = uc.sessionRepo.CreateSession(ctx, session)
	if err != nil {
		return nil, "", "", err
	}

	return user, accessToken, refreshToken, nil
}

func (uc *AuthUseCase) RefreshToken(ctx context.Context, refreshToken string) (*entity.User, string, error) {
	session, err := uc.sessionRepo.GetSessionByRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, "", err
	}
	if session == nil {
		return nil, "", errors.New("invalid refresh token")
	}

	if time.Now().After(session.ExpiresAtRefreshToken) {
		return nil, "", errors.New("refresh token expired")
	}

	user, err := uc.userRepo.GetUserByID(ctx, session.UserID)
	if err != nil {
		return nil, "", err
	}
	if user == nil {
		return nil, "", errors.New("user not found")
	}

	newRefreshToken, err := uc.tokenManager.GenerateRefreshToken()
	if err != nil {
		return nil, "", err
	}

	session.RefreshToken = newRefreshToken
	session.ExpiresAtRefreshToken = time.Now().Add(tokens.RefreshTokenTTL)

	err = uc.sessionRepo.UpdateSession(ctx, session)
	if err != nil {
		return nil, "", err
	}

	return user, newRefreshToken, nil
}

func (uc *AuthUseCase) Logout(ctx context.Context, refreshToken string) error {
	session, err := uc.sessionRepo.GetSessionByRefreshToken(ctx, refreshToken)
	if err != nil {
		return err
	}
	if session == nil {
		return nil
	}

	return uc.sessionRepo.DeleteSession(ctx, session.ID)
}

func (uc *AuthUseCase) GetUser(ctx context.Context, id int64) (*entity.User, error) {
	return uc.userRepo.GetUserByID(ctx, id)
}
