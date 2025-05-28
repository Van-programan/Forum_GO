package postgres

import (
	"context"
	"time"

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

func NewPostgresAuth(ctx context.Context, url string) *Postgres {
	return newPostgres(ctx, url)
}

func NewPostgresForum(ctx context.Context, url string) *Postgres {
	return newPostgres(ctx, url)
}

func (p *Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}
