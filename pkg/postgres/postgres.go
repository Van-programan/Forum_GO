package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	_defaultMaxPoolSize  = 5
	_defaultConnAttempts = 10
	_defaultConnTimeout  = time.Second
)

type Postgres struct {
	maxPoolSize  int
	connAttempts int
	connTimeout  time.Duration
	Pool         DBPool
}

type DBPool interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Close()
	Config() *pgxpool.Config
}

type Option func(*Postgres)

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

func NewWithPool(pool DBPool) *Postgres {
	return &Postgres{
		Pool: pool,
	}
}
