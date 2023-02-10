package util

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	pgxDriver "github.com/golang-migrate/migrate/v4/database/pgx"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/pgx/v4/stdlib"
	"grove/internal/resource"
)

type Migrator interface {
	Up() error
	Close() (source error, database error)
}

func NewMigrator(pool *pgxpool.Pool) (Migrator, error) {
	conn := stdlib.OpenDB(*pool.Config().ConnConfig)
	driver, err := pgxDriver.WithInstance(conn, new(pgxDriver.Config))
	if err != nil {
		return nil, err
	}

	source, err := httpfs.New(resource.Migrations, "/")
	if err != nil {
		return nil, fmt.Errorf("failed to create httpfs source: %w", err)
	}

	m, err := migrate.NewWithInstance(
		"httpfs", source,
		"pgx", driver)

	if err != nil {
		return nil, fmt.Errorf("failed to create migrator: %w", err)
	}

	return m, nil
}

func NoMigrationNeeded(err error) bool {
	return errors.Is(err, migrate.ErrNoChange)
}
