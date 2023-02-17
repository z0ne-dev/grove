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

	"grove/internal/service"

	"cdr.dev/slog"
	"github.com/go-chi/chi"
	"github.com/ztrue/shutdown"
)

var _ App = (*app)(nil)

type App interface {
	ConfigureRouter() error
	MigrateDatabase() error
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
	slogger := a.container.Logger().Named("server")

	go func(h *http.Server) {
		err := h.ListenAndServe()
		slogger.Critical(context.Background(), "server error", slog.Error(err))
	}(server)
	defer func(s *http.Server) {
		_ = s.Close() // No need to handle error here
	}(server)

	slogger.Info(context.Background(), "server started", slog.F("address", configuration.Http.PublicAddress), slog.F("listen", configuration.Http.Listen))

	shutdown.Add(func() {
		slogger.Debug(context.Background(), "shutting down server")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err := server.Shutdown(ctx)
		if err != nil {
			slogger.Fatal(context.Background(), "error shutting down server", slog.Error(err))
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
