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
	"github.com/Van-programan/Forum_GO/internal/usecase"
	"github.com/Van-programan/Forum_GO/pkg/logger"
	"github.com/Van-programan/Forum_GO/pkg/migrator"
	"github.com/Van-programan/Forum_GO/pkg/postgres"
	"github.com/Van-programan/Forum_GO/pkg/proto/authservice"
	"github.com/Van-programan/Forum_GO/pkg/tokens"
	"google.golang.org/grpc"
)

func RunAuthServer() {
	logger := logger.New("info")
	logger.Info("Starting auth service...")

	cfg, err := config.NewConfigAuth()
	if err != nil {
		logger.Fatal("Failed to load config", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pg := postgres.NewPostgresAuth(ctx, cfg)
	if pg == nil {
		logger.Fatal("Failed to connect to database", nil)
	}
	defer pg.Close()

	logger.Info("Successfully connected to database")

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		cfg.PGAuth.DBUser,
		cfg.PGAuth.DBPassword,
		cfg.PGAuth.DBHost,
		cfg.PGAuth.DBPort,
		cfg.PGAuth.DBName,
	)

	migrationsPath := filepath.Join("migrations", "auth")

	migrator := migrator.NewMigrator(dbURL, migrationsPath, logger)
	defer migrator.Close()
	migrator.Up()

	userRepo := repo.NewUserRepository(pg)
	sessionRepo := repo.NewSessionRepository(pg)

	tokenManager := tokens.NewTokenManager(cfg.JWT.JWTSecretKey)

	authUC := usecase.NewAuthUseCase(userRepo, sessionRepo, tokenManager)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(loggingInterceptor(logger)),
	)

	authController := controller.NewAuthController(authUC, tokenManager)
	authservice.RegisterAuthServiceServer(grpcServer, authController.UnimplementedAuthServiceServer)

	listener, err := net.Listen("tcp", net.JoinHostPort("", cfg.AuthInfo.GRPCPort))
	if err != nil {
		logger.Fatal("Failed to create listener", err)
	}

	go func() {
		logger.Info("Auth gRPC server started")
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

func loggingInterceptor(logger logger.Interface) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		start := time.Now()

		logger.Debug("gRPC request started",
			"method", info.FullMethod,
			"request", req)

		resp, err = handler(ctx, req)

		logger.Debug("gRPC request completed",
			"method", info.FullMethod,
			"duration", time.Since(start),
			"error", err)

		return resp, err
	}
}
