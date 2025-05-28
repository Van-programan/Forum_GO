package usecase

import (
	"context"
	"fmt"

	"github.com/Van-programan/Forum_GO/internal/repo"
	"github.com/rs/zerolog"
)

type UserUsecase interface {
	GetUsernamesByIds(ctx context.Context, ids []int64) (map[int64]string, error)
	GetUsernameById(ctx context.Context, id int64) (string, error)
}

type userUsecase struct {
	repo repo.UserGprcRepository
	log  *zerolog.Logger
}

func New(repo repo.UserGprcRepository, log *zerolog.Logger) UserUsecase {
	return &userUsecase{repo, log}
}

func (u *userUsecase) GetUsernamesByIds(ctx context.Context, ids []int64) (map[int64]string, error) {
	if len(ids) == 0 {
		u.log.Warn().Str("op", "UserUsecase.GetUsernamesByIds").Msg("empty ids")
		return make(map[int64]string), nil
	}

	usernames, err := u.repo.GetUsernamesByIds(ctx, ids)
	if err != nil {
		u.log.Error().Err(err).Str("op", "UserUsecase.GetUsernamesByIds").Msg("failed to get usernames by ids")
		return nil, fmt.Errorf("UserService - UserUsecase - GetUsernamesByIds - repo.GetUsernamesByIds: %w", err)
	}

	u.log.Info().Str("op", "UserUsecase.GetUsernamesByIds").Msg("success")
	return usernames, nil
}

func (u *userUsecase) GetUsernameById(ctx context.Context, id int64) (string, error) {
	username, err := u.repo.GetUsernameById(ctx, id)
	if err != nil {
		u.log.Error().Err(err).Str("op", "UserUsecase.GetUsernameById").Msg("failed to get username by id")
		return "", fmt.Errorf("UserService - UserUsecase - GetUsernameById - repo.GetUsernameById: %w", err)
	}

	u.log.Info().Str("op", "UserUsecase.GetUsernameById").Msg("success")
	return username, nil
}
