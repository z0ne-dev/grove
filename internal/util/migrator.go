package util

import (
	"github.com/golang-migrate/migrate/v4"
	pgxDriver "github.com/golang-migrate/migrate/v4/database/pgx"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
	"github.com/jackc/pgx/v4/stdlib"
	"grove/internal/resource"
	"grove/internal/service"
)

type Migrator interface {
	Up() error
	Close() (source error, database error)
}

func NewMigrator(c service.Container) (Migrator, error) {
	pool, err := c.PostgresPool()
	if err != nil {
		return nil, err
	}

	conn := stdlib.OpenDB(*pool.Config().ConnConfig)
	driver, err := pgxDriver.WithInstance(conn, new(pgxDriver.Config))
	if err != nil {
		return nil, err
	}

	source, err := httpfs.New(resource.Migrations, "/")
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithInstance(
		"httpfs", source,
		"pgx", driver)

	return m, err
}
