// container.go Copyright (c) 2023 z0ne.
// All Rights Reserved.
// Licensed under the EUPL 1.2 License.
// See LICENSE the project root for license information.
//
// SPDX-License-Identifier: EUPL-1.2

package service

import (
	"net/http"
	"time"

	"github.com/z0ne-dev/grove/internal/config"
	"github.com/z0ne-dev/grove/internal/resource"
	"github.com/z0ne-dev/grove/internal/util"

	"cdr.dev/slog"
	"github.com/CloudyKit/jet/v6"
	"github.com/CloudyKit/jet/v6/loaders/httpfs"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/multierr"
)

var _ Container = (*container)(nil)

type containerInitFunction func(c *container) error

var containerInitFunctions = make([]containerInitFunction, 0, 10)

func addContainerInit(f containerInitFunction) {
	containerInitFunctions = append(containerInitFunctions, f)
}

type container struct {
	logger slog.Logger
	config *config.Config

	router chi.Router
	server *http.Server
	jet    *jet.Set

	pgxPool *pgxpool.Pool
}

// Jet returns the jet template engine.
func (c *container) Jet() *jet.Set {
	return c.jet
}

// Server returns the http server.
func (c *container) Server() *http.Server {
	return c.server
}

// Router returns the http router.
func (c *container) Router() chi.Router {
	return c.router
}

// Logger returns the logger.
func (c *container) Logger() slog.Logger {
	return c.logger
}

// Config returns the config.
func (c *container) Config() *config.Config {
	return c.config
}

// New creates a new service container.
func New(logger slog.Logger, config *config.Config) (Container, error) {
	router := createRouter(&logger)

	loader, err := httpfs.NewLoader(resource.Templates)
	if err != nil {
		panic(err)
	}

	const readHeaderTimeout = 30 * time.Second

	c := &container{
		logger: logger,
		config: config,
		router: router,
		server: &http.Server{
			Addr:              config.Http.Listen,
			Handler:           router,
			ReadHeaderTimeout: readHeaderTimeout,
		},

		/* development mode for templates ignores cache -
		which is not a performace issue in production
		(evereything is loaded from memory)	*/
		jet: jet.NewSet(loader, jet.InDevelopmentMode()),
	}

	// var err error
	for _, f := range containerInitFunctions {
		initErr := f(c)
		err = multierr.Append(err, initErr)
	}

	if err != nil {
		return nil, err
	}

	return c, nil
}

func createRouter(logger *slog.Logger) *chi.Mux {
	router := chi.NewRouter()
	namedLogger := logger.Named("http")
	router.Use(
		// catch panics
		middleware.Recoverer,

		// clean up paths
		middleware.CleanPath,
		middleware.StripSlashes,

		// logging relevant data structures
		middleware.RequestID,
		middleware.RealIP,

		// logger
		middleware.RequestLogger(util.NewSlogChiFormatter(&namedLogger)),

		// security
		cors.New(cors.Options{
			AllowCredentials: false,
			AllowedMethods:   []string{"GET", "HEAD", "OPTIONS", "POST", "PUT", "DELETE", "PATCH"},
			AllowedOrigins:   []string{"*"}, // Federation must support all origins
		}).Handler,

		// forward head if needed
		middleware.GetHead,

		middleware.Heartbeat("/healthz"),
	)

	router.Handle("/*", http.FileServer(resource.Assets))
	return router
}
