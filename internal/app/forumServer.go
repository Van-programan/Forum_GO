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
	"github.com/Van-programan/Forum_GO/internal/client"
	"github.com/Van-programan/Forum_GO/internal/controller/route"
	"github.com/Van-programan/Forum_GO/internal/repo"
	"github.com/Van-programan/Forum_GO/internal/usecase"
	"github.com/Van-programan/Forum_GO/internal/ws"
	"github.com/Van-programan/Forum_GO/pkg/httpserver"
	"github.com/Van-programan/Forum_GO/pkg/jwt"
	"github.com/Van-programan/Forum_GO/pkg/logger"
	"github.com/Van-programan/Forum_GO/pkg/migrator"
	"github.com/Van-programan/Forum_GO/pkg/postgres"
)

func RunForumServer() {

	cfg, err := config.NewConfigForum()
	if err != nil {
		log.Fatal("Failed to load config", err)
	}

	logger := logger.New("forum-service", cfg.Log.LogLevel)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		cfg.PGForum.DBUser,
		cfg.PGForum.DBPassword,
		cfg.PGForum.DBHost,
		cfg.PGForum.DBPort,
		cfg.PGForum.DBName,
	)

	pg := postgres.NewPostgresForum(ctx, dbURL)
	if pg == nil {
		log.Fatalf("Failed to connect to database")
	}
	defer pg.Close()

	log.Println("Successfully connected to database")

	migrationsPath := "migrations/forum"

	migrator := migrator.NewMigrator(dbURL, migrationsPath, *logger)
	defer migrator.Close()
	migrator.Up()

	categoryRepo := repo.NewCategoryRepository(pg, logger)
	topicRepo := repo.NewTopicRepository(pg, logger)
	postRepo := repo.NewPostRepository(pg, logger)
	chatRepo := repo.NewChatRepository(pg, logger)

	userClient, err := client.New(cfg.ForumInfo.GRPCPort, logger)
	if err != nil {
		log.Fatalf("app - Run - client.New: %v", err)
	}
	defer userClient.Close()

	categoryUC := usecase.NewCategoryUsecase(categoryRepo, logger)
	topicUC := usecase.NewTopicUsecase(topicRepo, categoryRepo, userClient, logger)
	postUC := usecase.NewPostUsecase(postRepo, topicRepo, userClient, logger)

	jwt := jwt.New(cfg.JWT.Secret, cfg.JWT.Access_TTL, cfg.JWT.Refresh_TTL)

	hub := ws.NewHub(logger)
	go hub.Run()
	chatUC := usecase.NewChatUsecase(chatRepo, logger)

	httpServer := httpserver.New(cfg.ForumInfo.Server)
	route.NewForumRouter(httpServer.Engine, categoryUC, topicUC, postUC, jwt, logger, hub, chatUC, userClient)
	httpServer.Run()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	<-interrupt
}
