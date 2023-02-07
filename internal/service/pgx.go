package service

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"grove/internal/util"
)

func init() {
	addContainerInit(func(c *container) error {
		_, err := c.createPostgresPool()

		return err
	})
}

func (c *container) createPostgresPool() (*pgxpool.Pool, error) {
	if c.pgxPool == nil {
		dsn, err := pgxpool.ParseConfig(c.Config().Postgres)
		if err != nil {
			return nil, fmt.Errorf("failed to parse postgres dsn: %w", err)
		}

		dsn.ConnConfig.Logger = pgx.LoggerFunc(util.SlogPgxLogger(c.Logger().Named("sql")))
		dsn.ConnConfig.LogLevel = pgx.LogLevelDebug

		pool, err := pgxpool.ConnectConfig(context.Background(), dsn)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to postgres: %w", err)
		}

		c.pgxPool = pool
	}

	return c.pgxPool, nil
}

func (c *container) PostgresPool() (*pgxpool.Pool, error) {
	return c.createPostgresPool()
}

func (c *container) Postgres() (*pgxpool.Conn, error) {
	pool, err := c.PostgresPool()
	if err != nil {
		return nil, err
	}

	return pool.Acquire(context.Background())
}

func (c *container) Migrator() (util.Migrator, error) {
	p, err := c.PostgresPool()
	if err != nil {
		return nil, err
	}
	return util.NewMigrator(p)
}
