package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/Van-programan/Forum_GO/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	Pool *pgxpool.Pool
}

var poolConfig = &pgxpool.Config{
	MaxConns:          10,
	MinConns:          2,
	MaxConnLifetime:   time.Hour,
	MaxConnIdleTime:   30 * time.Minute,
	HealthCheckPeriod: time.Minute,
}

func newPostgres(ctx context.Context, dsn string) *Postgres {
	const op = "storage.postgres.New"

	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil
	}

	poolCfg.MaxConns = poolConfig.MaxConns
	poolCfg.MinConns = poolConfig.MinConns
	poolCfg.MaxConnLifetime = poolConfig.MaxConnLifetime
	poolCfg.MaxConnIdleTime = poolConfig.MaxConnIdleTime
	poolCfg.HealthCheckPeriod = poolConfig.HealthCheckPeriod

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil
	}

	if err := pool.Ping(ctx); err != nil {
		return nil
	}

	return &Postgres{
		Pool: pool,
	}
}

func NewPostgresAuth(ctx context.Context, cfg *config.ConfigAuth) *Postgres {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.PGAuth.DBHost,
		cfg.PGAuth.DBPort,
		cfg.PGAuth.DBUser,
		cfg.PGAuth.DBPassword,
		cfg.PGAuth.DBName,
		cfg.PGAuth.DBSSLMode,
	)
	return newPostgres(ctx, dsn)
}

func NewPostgresForum(ctx context.Context, cfg *config.ConfigForum) *Postgres {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.PGForum.DBHost,
		cfg.PGForum.DBPort,
		cfg.PGForum.DBUser,
		cfg.PGForum.DBPassword,
		cfg.PGForum.DBName,
		cfg.PGForum.DBSSLMode,
	)
	return newPostgres(ctx, dsn)
}

func (p *Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}
