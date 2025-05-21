package app

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/Van-programan/Forum_GO/config"
	"github.com/Van-programan/Forum_GO/internal/controller"
	"github.com/Van-programan/Forum_GO/internal/repo"
	"github.com/Van-programan/Forum_GO/internal/transport/ws"
	"github.com/Van-programan/Forum_GO/internal/usecase"
	"github.com/Van-programan/Forum_GO/pkg/logger"
	"github.com/Van-programan/Forum_GO/pkg/migrator"
	"github.com/Van-programan/Forum_GO/pkg/postgres"
	"github.com/Van-programan/Forum_GO/pkg/proto/forumservice"
	"github.com/Van-programan/Forum_GO/pkg/tokens"
	"google.golang.org/grpc"
)

func RunForumService() {
	logger := logger.New("info")
	logger.Info("Starting forum service...")

	cfg, err := config.NewConfigForum()
	if err != nil {
		logger.Fatal("Failed to load config", err)
	}

	cfgAuth, err := config.NewConfigAuth()
	if err != nil {
		logger.Fatal("Failed to load config", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pg := postgres.NewPostgresForum(ctx, cfg)
	if pg == nil {
		logger.Fatal("Failed to connect to database", nil)
	}
	defer pg.Close()

	logger.Info("Successfully connected to database")

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		cfg.PGForum.DBUser,
		cfg.PGForum.DBPassword,
		cfg.PGForum.DBHost,
		cfg.PGForum.DBPort,
		cfg.PGForum.DBName,
	)

	migrationsPath := filepath.Join("migrations", "forum")

	migrator := migrator.NewMigrator(dbURL, migrationsPath, logger)
	defer migrator.Close()
	migrator.Up()

	messageRepo := repo.NewMessageRepository(pg)
	topicRepo := repo.NewTopicRepository(pg)

	wsHub := ws.NewHub()
	go wsHub.Run()

	tokenManager := tokens.NewTokenManager(cfgAuth.JWT.JWTSecretKey)

	forumUC := usecase.NewForumUseCase(topicRepo, messageRepo, wsHub, tokenManager)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(loggingInterceptor(logger)),
	)

	forumController := controller.NewForumController(forumUC)
	forumservice.RegisterForumServiceServer(grpcServer, forumController.UnimplementedForumServiceServer)

	listener, err := net.Listen("tcp", net.JoinHostPort("", cfg.ForumInfo.GRPCPort))
	if err != nil {
		logger.Fatal("Failed to create listener", err)
	}

	go func() {
		logger.Info("Forum gRPC server started")
		if err := grpcServer.Serve(listener); err != nil {
			logger.Fatal("gRPC server failed", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	stopped := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(stopped)
	}()

	select {
	case <-time.After(5 * time.Second):
		logger.Warn("Server shutdown timed out, forcing exit")
		grpcServer.Stop()
	case <-stopped:
		logger.Info("Server stopped gracefully")
	}

	logger.Info("Server shutdown completed")
}
