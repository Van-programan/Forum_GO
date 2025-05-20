package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
)

type (
	ConfigAuth struct {
		App      App
		AuthInfo AuthInfo
		Log      Log
		JWT      JWT
		PGAuth   PGAuth
		Swagger  Swagger
	}

	ConfigForum struct {
		App         App
		ForumInfo   ForumInfo
		Log         Log
		AddressAuth AddressAuth
		PGForum     PGForum
		Swagger     Swagger
	}

	App struct {
		AppName     string `env:"APP_NAME" envDefault:"Forum_go"`
		Environment string `env:"ENVIRONMENT" envDefault:"development"`
		AppVersion  string `env:"APP_VERSION" envDefault:"1.0.0"`
	}

	AuthInfo struct {
		GRPCPort string        `env:"AUTH_SERVICE_GRPC_PORT" envDefault:"50051"`
		Timeout  time.Duration `env:"TIMEOUT" envDefault:"30s"`
	}

	ForumInfo struct {
		GRPCPort string        `env:"FORUM_SERVICE_GRPC_PORT" envDefault:"50052"`
		WSPort   string        `env:"FORUM_SERVICE_WS_PORT" envDefault:"8042"`
		Timeout  time.Duration `env:"TIMEOUT" envDefault:"30s"`
	}

	Log struct {
		LogLevel string `env:"LOG_LEVEL" envDefault:"debug"`
	}

	PGConfig struct {
		DBHost     string `env:"DB_HOST" envDefault:"localhost"`
		DBPort     int    `env:"DB_PORT" envDefault:"5432"`
		DBUser     string `env:"DB_USER" envDefault:"postgres"`
		DBPassword string `env:"DB_PASSWORD" envDefault:"1234"`
		DBSSLMode  string `env:"DB_SSLMODE" envDefault:"disable"`
	}

	PGAuth struct {
		PGConfig
		DBName string `env:"AUTH_DB_NAME" envDefault:"auth"`
	}

	PGForum struct {
		PGConfig
		DBName string `env:"FORUM_DB_NAME" envDefault:"forum"`
	}

	JWT struct {
		JWTSecretKey      string        `env:"JWT_SECRET_KEY" envDefault:"very-secret-key-which-impossible-hack"`
		JWTExpirationTime time.Duration `env:"JWT_EXPIRATION" envDefault:"15m"`
	}

	AddressAuth struct {
		Address string `env:"AUTH_SERVICE_ADDR" envDefault:"localhost:50051"`
	}

	Swagger struct {
		Enabled bool `env:"SWAGGER_ENABLED" envDefault:"false"`
	}
)

func NewConfigAuth() (*ConfigAuth, error) {
	cfg := &ConfigAuth{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}

func NewConfigForum() (*ConfigForum, error) {
	cfg := &ConfigForum{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}
