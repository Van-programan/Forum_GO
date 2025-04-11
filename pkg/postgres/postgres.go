package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/Van-programan/Forum_GO/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type Postgres struct {
	Pool *pgxpool.Pool
	log  *zerolog.Logger
}

var poolConfig = &pgxpool.Config{
	MaxConns:          10,
	MinConns:          2,
	MaxConnLifetime:   time.Hour,
	MaxConnIdleTime:   30 * time.Minute,
	HealthCheckPeriod: time.Minute,
}

func newPostgres(ctx context.Context, dsn string, logger *zerolog.Logger) (*Postgres, error) {
	const op = "storage.postgres.New"

	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	poolCfg.MaxConns = poolConfig.MaxConns
	poolCfg.MinConns = poolConfig.MinConns
	poolCfg.MaxConnLifetime = poolConfig.MaxConnLifetime
	poolCfg.MaxConnIdleTime = poolConfig.MaxConnIdleTime
	poolCfg.HealthCheckPeriod = poolConfig.HealthCheckPeriod

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("%s: ping failed: %w", op, err)
	}

	logger.Info().Msg("PostgreSQL connected successfully")

	return &Postgres{
		Pool: pool,
		log:  logger,
	}, nil
}

func NewPostgresAuth(ctx context.Context, cfg *config.ConfigAuth, logger *zerolog.Logger) (*Postgres, error) {
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

func NewPostgresForum(ctx context.Context, cfg *config.ConfigForum, logger *zerolog.Logger) (*Postgres, error) {
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
		p.log.Info().Msg("PostgreSQL connection closed")
	}
}
