// interface.go Copyright (c) 2023 z0ne.
// All Rights Reserved.
// Licensed under the EUPL 1.2 License.
// See LICENSE the project root for license information.
//
// SPDX-License-Identifier: EUPL-1.2

package service

import (
	"net/http"

	"github.com/z0ne-dev/grove/internal/config"

	"github.com/CloudyKit/jet/v6"
	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/exp/slog"
)

// Container is the service container.
type Container interface {
	// Logger returns the logger.
	Logger() *slog.Logger
	// Config returns the config.
	Config() *config.Config
	// Router returns the http router.
	Router() chi.Router
	// Server returns the http server.
	Server() *http.Server
	// Jet returns the jet template engine.
	Jet() *jet.Set
	// PostgresPool returns the postgres pool.
	PostgresPool() (*pgxpool.Pool, error)
	// Postgres returns a new connection from the pool.
	Postgres() (*pgxpool.Conn, error)
}
