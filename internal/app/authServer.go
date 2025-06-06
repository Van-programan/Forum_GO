package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"syscall"
	"time"

	"github.com/Van-programan/Forum_GO/config"
	"github.com/Van-programan/Forum_GO/internal/controller/route"
	"github.com/Van-programan/Forum_GO/internal/repo"
	"github.com/Van-programan/Forum_GO/internal/usecase"
	"github.com/Van-programan/Forum_GO/pkg/httpserver"
	"github.com/Van-programan/Forum_GO/pkg/jwt"
	"github.com/Van-programan/Forum_GO/pkg/logger"
	"github.com/Van-programan/Forum_GO/pkg/migrator"
	"github.com/Van-programan/Forum_GO/pkg/postgres"
)

func RunAuthServer() {

	cfg, err := config.NewConfigAuth()
	if err != nil {
		log.Fatal("Failed to load config", err)
	}

	logger := logger.New("auth-service", cfg.Log.LogLevel)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		cfg.PGAuth.DBUser,
		cfg.PGAuth.DBPassword,
		cfg.PGAuth.DBHost,
		cfg.PGAuth.DBPort,
		cfg.PGAuth.DBName,
	)

	pg := postgres.NewPostgresAuth(ctx, dbURL)
	if pg == nil {
		log.Fatalf("Failed to connect to database")
	}
	defer pg.Close()

	log.Println("Successfully connected to database")

	migrationsPath := "migrations/auth"

	migrator := migrator.NewMigrator(dbURL, migrationsPath, *logger)
	defer migrator.Close()
	migrator.Up()
	userRepo := repo.NewUserRepository(pg, logger)
	tokenRepo := repo.NewRefreshTokenRepository(pg, logger)

	jwt := jwt.New(cfg.JWT.Secret, cfg.JWT.Access_TTL, cfg.JWT.Refresh_TTL)

	authUC := usecase.NewAuthUsecase(userRepo, tokenRepo, jwt, logger)
	fmt.Println("04")
	httpServer := httpserver.New(cfg.AuthInfo.Server)
	route.NewAuthRouter(httpServer.Engine, authUC, jwt, logger)

	httpServer.Run()
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	<-interrupt
}
