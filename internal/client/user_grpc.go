package client

import (
	"context"
	"fmt"
	"strings"
	"time"

	userpb "github.com/Van-programan/Forum_GO/pkg/proto"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserClient interface {
	GetUsernames(ctx context.Context, userIDs []int64) (map[int64]string, error)
	GetUsername(ctx context.Context, userID int64) (string, error)
	Close() error
}

type userClient struct {
	client userpb.UserServiceClient
	conn   *grpc.ClientConn
	log    *zerolog.Logger
}

func New(address string, log *zerolog.Logger) (UserClient, error) {
	if !strings.Contains(address, ":") {
		address = ":" + address
	}

	conn, err := grpc.Dial(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithTimeout(5*time.Second),
	)
	if err != nil {
		return nil, fmt.Errorf("grpc.Dial failed: %w", err)
	}

	client := userpb.NewUserServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err = client.GetUsernames(ctx, &userpb.GetUsernamesRequest{UserIds: []int64{0}})
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("health check failed: %w", err)
	}

	return &userClient{
		client: client,
		conn:   conn,
		log:    log,
	}, nil
}

func (c *userClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *userClient) GetUsernames(ctx context.Context, userIDs []int64) (map[int64]string, error) {
	if len(userIDs) == 0 {
		c.log.Warn().Str("op", "UserClient.GetUsernames").Msg("Empty userIDs")
		return make(map[int64]string), nil
	}

	req := &userpb.GetUsernamesRequest{
		UserIds: userIDs,
	}

	callCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := c.client.GetUsernames(callCtx, req)
	if err != nil {
		c.log.Error().Err(err).Str("op", "UserClient.GetUsernames").Any("userIDs", userIDs).Msg("Failed to get usernames")
		return nil, fmt.Errorf("clients.user - GetUsernames - c.client.GetUsernames: %w", err)
	}

	c.log.Info().Str("op", "UserClient.GetUsernames").Msg("Successfully got usernames")
	return res.GetUsernames(), nil
}

func (c *userClient) GetUsername(ctx context.Context, userID int64) (string, error) {
	req := &userpb.GetUsernameRequest{
		UserId: userID,
	}

	callCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := c.client.GetUsername(callCtx, req)
	if err != nil {
		c.log.Error().Err(err).Str("op", "UserClient.GetUsername").Any("userID", userID).Msg("Failed to get username")
		return "", fmt.Errorf("clients.user - GetUsername - c.client.GetUsername: %w", err)
	}

	c.log.Info().Str("op", "UserClient.GetUsername").Msg("Successfully got username")
	return res.GetUsername(), nil
}
