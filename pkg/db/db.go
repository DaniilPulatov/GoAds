package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Pool interface {
	Ping(ctx context.Context) error
	Close()

	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	Begin(ctx context.Context) (pgx.Tx, error)
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
}

type conn struct {
	*pgxpool.Pool
}
/*
// NewDB creates a new DB instance
func NewDB(dsn string) (Pool, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse dsn: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping pool: %w", err)
	}

	return &conn{
		Pool: pool,
	}, nil
}
*/
func (c *conn) Ping(ctx context.Context) error {
	if err := c.Pool.Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping pool: %w", err)
	}
	return nil
}

func (c *conn) Close() {
	c.Pool.Close()
}

func (c *conn) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return c.Pool.QueryRow(ctx, sql, args...)
}

func (c *conn) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	rows, err := c.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}
	return rows, nil
}

func (c *conn) Begin(ctx context.Context) (pgx.Tx, error) {
	tx, err := c.Pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	return tx, nil
}

func (c *conn) Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
	tag, err := c.Pool.Exec(ctx, sql, arguments...)
	if err != nil {
		return tag, fmt.Errorf("failed to exec: %w", err)
	}
	return tag, nil
}

func NewDB(dsn string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse dsn: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping pool: %w", err)
	}

	return pool, nil
}
