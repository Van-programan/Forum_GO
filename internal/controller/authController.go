package controller

import (
	"context"
	"fmt"

	"github.com/Van-programan/Forum_GO/internal/entity"
	"github.com/Van-programan/Forum_GO/internal/usecase"
	"github.com/Van-programan/Forum_GO/pkg/proto/authservice"
	"github.com/Van-programan/Forum_GO/pkg/tokens"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AuthGRPCController struct {
	authservice.UnimplementedAuthServiceServer
	authUC usecase.AuthUseCase
	tm     tokens.TokenManager
}

func NewAuthController(authUC usecase.AuthUseCase, tm tokens.TokenManager) AuthGRPCController {
	return AuthGRPCController{authUC: authUC, tm: tm}
}

func (c *AuthGRPCController) Register(ctx context.Context, req *authservice.RegisterRequest) (*authservice.RegisterResponse, error) {
	user, err := c.authUC.Register(ctx, req.Username, req.Email, req.Password)
	if err != nil {
		return nil, fmt.Errorf("registration failed: %v", err)
	}

	return &authservice.RegisterResponse{
		User: &authservice.User{
			Id:           user.ID,
			Username:     user.Username,
			Email:        user.Email,
			RegisteredAt: timestamppb.New(user.RegisteredAt),
		},
	}, nil
}

func (c *AuthGRPCController) Login(ctx context.Context, req *authservice.LoginRequest) (*authservice.LoginResponse, error) {
	user, accessToken, refreshToken, err := c.authUC.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, fmt.Errorf("login failed: %v", err)
	}

	return &authservice.LoginResponse{
		User:         toProtoUser(user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (c *AuthGRPCController) RefreshToken(ctx context.Context, req *authservice.RefreshTokenRequest) (*authservice.RefreshTokenResponse, error) {
	user, newRefreshToken, err := c.authUC.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("token refresh failed: %v", err)
	}

	accessToken, err := c.tm.GenerateAccessToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %v", err)
	}

	return &authservice.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (c *AuthGRPCController) ValidateToken(ctx context.Context, req *authservice.ValidateTokenRequest) (*authservice.ValidateTokenResponse, error) {
	userID, err := c.tm.ParseAccessToken(req.AccessToken)
	if err != nil {
		return &authservice.ValidateTokenResponse{
			Valid: false,
		}, nil
	}

	return &authservice.ValidateTokenResponse{
		UserId: userID,
		Valid:  true,
	}, nil
}

func toProtoUser(user *entity.User) *authservice.User {
	return &authservice.User{
		Id:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		RegisteredAt: timestamppb.New(user.RegisteredAt),
	}
}
