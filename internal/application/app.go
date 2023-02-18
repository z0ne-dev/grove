// app.go Copyright (c) 2023 z0ne.
// All Rights Reserved.
// Licensed under the EUPL 1.2 License.
// See LICENSE the project root for license information.
//
// SPDX-License-Identifier: EUPL-1.2

package application

import (
	"context"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/z0ne-dev/grove/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/ztrue/shutdown"
	"golang.org/x/exp/slog"
)

var _ App = (*app)(nil)

// App is the main application interface.
type App interface {
	// ConfigureRouter configures the router.
	ConfigureRouter() error
	// MigrateDatabase migrates the database.
	MigrateDatabase() error
	// ListenAndServe starts the server.
	ListenAndServe()
}

type app struct {
	container service.Container
}

func (a *app) ConfigureRouter() error {
	r := a.container.Router()
	set := a.container.Jet()
	pool, err := a.container.PostgresPool()
	if err != nil {
		return err
	}

	r.Group(func(r chi.Router) {
		r.With(jetGlobalsMiddleware(set, pool))
		NewGenericRoutes(set).Routes(r)
	})

	return nil
}

func (a *app) ListenAndServe() {
	server := a.container.Server()
	configuration := a.container.Config()
	slogger := a.container.Logger().WithGroup("server")

	go func(h *http.Server) {
		err := h.ListenAndServe()
		slogger.Error("server error", err)
	}(server)
	defer func(s *http.Server) {
		_ = s.Close() // No need to handle error here
	}(server)

	slogger.Info("server started", slog.String("address", configuration.HTTP.PublicAddress), slog.String("listen", configuration.HTTP.Listen))

	shutdown.Add(func() {
		slogger.Debug("shutting down server")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err := server.Shutdown(ctx)
		if err != nil {
			slogger.Error("error shutting down server", err)
			os.Exit(1)
		}
		cancel()
	})
	shutdown.Listen(syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
}

func New(container service.Container) App {
	return &app{
		container: container,
	}
}
