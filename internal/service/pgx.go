// pgx.go Copyright (c) 2023 z0ne.
// All Rights Reserved.
// Licensed under the EUPL 1.2 License.
// See LICENSE the project root for license information.
//
// SPDX-License-Identifier: EUPL-1.2

package service

import (
	"context"
	"fmt"

	"github.com/z0ne-dev/grove/internal/util"

	"github.com/jackc/pgx/v5/pgxpool"
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

		namedLogger := c.Logger().Named("sql")
		dsn.ConnConfig.Tracer = util.NewSlogPgxTracer(&namedLogger)

		pool, err := pgxpool.NewWithConfig(context.Background(), dsn)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to postgres: %w", err)
		}

		c.pgxPool = pool
	}

	return c.pgxPool, nil
}

// PostgresPool returns the postgres pool.
func (c *container) PostgresPool() (*pgxpool.Pool, error) {
	return c.createPostgresPool()
}

// Postgres returns a new connection from the pool.
func (c *container) Postgres() (*pgxpool.Conn, error) {
	pool, err := c.PostgresPool()
	if err != nil {
		return nil, err
	}

	return pool.Acquire(context.Background())
}
