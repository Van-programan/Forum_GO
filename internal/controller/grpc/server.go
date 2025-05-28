package grpc

import (
	"context"
	"fmt"

	"github.com/Van-programan/Forum_GO/internal/usecase"
	userpb "github.com/Van-programan/Forum_GO/pkg/proto"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

type serverAPI struct {
	userpb.UnimplementedUserServiceServer
	usecase usecase.UserUsecase
	log     *zerolog.Logger
}

func Register(grpcServer *grpc.Server, usecase usecase.UserUsecase, log *zerolog.Logger) {
	userpb.RegisterUserServiceServer(grpcServer, &serverAPI{usecase: usecase, log: log})
}

func (s *serverAPI) GetUsernames(ctx context.Context, req *userpb.GetUsernamesRequest) (*userpb.GetUsernamesResponse, error) {
	usernames, err := s.usecase.GetUsernamesByIds(ctx, req.UserIds)
	if err != nil {
		s.log.Error().Err(err).Str("op", "UserServer.GetUsernames").Msg("failed to get usernames")
		return nil, fmt.Errorf("UserService - UserServer - GetUsernames - usecase.GetUsernamesByIds: %w", err)
	}

	s.log.Info().Str("op", "UserServer.GetUsernames").Msg("success")
	return &userpb.GetUsernamesResponse{Usernames: usernames}, nil
}

func (s *serverAPI) GetUsername(ctx context.Context, req *userpb.GetUsernameRequest) (*userpb.GetUsernameResponse, error) {
	username, err := s.usecase.GetUsernameById(ctx, req.UserId)
	if err != nil {
		s.log.Error().Err(err).Str("op", "UserServer.GetUsername").Msg("failed to get username")
		return nil, fmt.Errorf("UserService - UserServer - GetUsername - usecase.GetUsernameById: %w", err)
	}

	s.log.Info().Str("op", "UserServer.GetUsername").Msg("success")
	return &userpb.GetUsernameResponse{UserId: req.UserId, Username: username}, nil
}
