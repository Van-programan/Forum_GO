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
		App       App
		ForumInfo ForumInfo
		Log       Log
		PGForum   PGForum
		Swagger   Swagger
		JWT       JWT
	}

	App struct {
		AppName     string `env:"APP_NAME" envDefault:"Forum_go"`
		Environment string `env:"ENVIRONMENT" envDefault:"development"`
		AppVersion  string `env:"APP_VERSION" envDefault:"1.0.0"`
	}

	AuthInfo struct {
		Server   string `env:"AUTH_SERVICE" envDefault:"3100"`
		GRPCPort string `env:"GRPC_PORT" envDefault:"50051"`
	}

	ForumInfo struct {
		Server   string `env:"FORUM_SERVICE" envDefault:"3101"`
		GRPCPort string `env:"GRPC_PORT" envDefault:"50051"`
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
		Access_TTL  time.Duration `env:"ACCESS_TTL" envDefault:"15m"`
		Refresh_TTL time.Duration `env:"REFRESH_TTL" envDefault:"720h"`
		Secret      string        `env:"SECRET" envDefault:"very-secret-value-impossible-hack"`
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
