// interface.go Copyright (c) 2023 z0ne.
// All Rights Reserved.
// Licensed under the EUPL 1.2 License.
// See LICENSE the project root for license information.
//
// SPDX-License-Identifier: EUPL-1.2

package service

import (
	"net/http"

	"grove/internal/config"

	"cdr.dev/slog"
	"github.com/CloudyKit/jet/v6"
	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Container interface {
	Logger() slog.Logger
	Config() *config.Config
	Router() chi.Router
	Server() *http.Server
	Jet() *jet.Set
	PostgresPool() (*pgxpool.Pool, error)
	Postgres() (*pgxpool.Conn, error)
}
