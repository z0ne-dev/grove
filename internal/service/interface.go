package service

import (
	"cdr.dev/slog"
	"github.com/CloudyKit/jet/v6"
	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v4/pgxpool"
	"grove/internal/config"
	"net/http"
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
