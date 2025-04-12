package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/Van-programan/Forum_GO/config"
	"github.com/Van-programan/Forum_GO/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	Pool *pgxpool.Pool
	log  logger.Logger
}

var poolConfig = &pgxpool.Config{
	MaxConns:          10,
	MinConns:          2,
	MaxConnLifetime:   time.Hour,
	MaxConnIdleTime:   30 * time.Minute,
	HealthCheckPeriod: time.Minute,
}

func newPostgres(ctx context.Context, dsn string, logger logger.Logger) *Postgres {
	const op = "storage.postgres.New"

	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		logger.Error("%s: %w", op, err)
		return nil
	}

	poolCfg.MaxConns = poolConfig.MaxConns
	poolCfg.MinConns = poolConfig.MinConns
	poolCfg.MaxConnLifetime = poolConfig.MaxConnLifetime
	poolCfg.MaxConnIdleTime = poolConfig.MaxConnIdleTime
	poolCfg.HealthCheckPeriod = poolConfig.HealthCheckPeriod

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		logger.Error("%s: %w", op, err)
		return nil
	}

	if err := pool.Ping(ctx); err != nil {
		logger.Error("%s: ping failed: %w", op, err)
		return nil
	}

	logger.Info("PostgreSQL connected successfully")

	return &Postgres{
		Pool: pool,
		log:  logger,
	}
}

func NewPostgresAuth(ctx context.Context, cfg *config.ConfigAuth, logger logger.Logger) *Postgres {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.PGAuth.DBHost,
		cfg.PGAuth.DBPort,
		cfg.PGAuth.DBUser,
		cfg.PGAuth.DBPassword,
		cfg.PGAuth.DBName,
		cfg.PGAuth.DBSSLMode,
	)
	return newPostgres(ctx, dsn, logger)
}

func NewPostgresForum(ctx context.Context, cfg *config.ConfigForum, logger logger.Logger) *Postgres {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.PGForum.DBHost,
		cfg.PGForum.DBPort,
		cfg.PGForum.DBUser,
		cfg.PGForum.DBPassword,
		cfg.PGForum.DBName,
		cfg.PGForum.DBSSLMode,
	)
	return newPostgres(ctx, dsn, logger)
}

func (p *Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
		p.log.Info("PostgreSQL connection closed")
	}
}
