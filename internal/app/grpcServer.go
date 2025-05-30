package app

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Van-programan/Forum_GO/config"
	"github.com/Van-programan/Forum_GO/internal/controller/grpc"
	"github.com/Van-programan/Forum_GO/internal/repo"
	"github.com/Van-programan/Forum_GO/internal/usecase"
	"github.com/Van-programan/Forum_GO/pkg/logger"
	"github.com/Van-programan/Forum_GO/pkg/postgres"
	Grpc "google.golang.org/grpc"
)

func RunGrpcServer() {

	cfg, err := config.NewConfigAuth()
	if err != nil {
		log.Fatal("Failed to load config", err)
	}

	logger := logger.New("user-service", cfg.Log.LogLevel)

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

	userRepo := repo.New(pg, logger)

	userUsecase := usecase.New(userRepo, logger)

	grpcServer := Grpc.NewServer()
	grpc.Register(grpcServer, userUsecase, logger)

	l, err := net.Listen("tcp", ":"+cfg.AuthInfo.GRPCPort)
	if err != nil {
		log.Fatalf("app - Run - net.Listen: %v", err)
	}

	fmt.Println("grpc run")

	go func() {
		if err := grpcServer.Serve(l); err != nil {
			log.Fatalf("app - Run - grpcServer.Serve: %v", err)
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	<-interrupt
}
