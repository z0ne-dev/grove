package application

import (
	"cdr.dev/slog"
	"context"
	"github.com/ztrue/shutdown"
	"grove/internal/service"
	"net/http"
	"os"
	"syscall"
	"time"
)

var _ App = (*app)(nil)

type App interface {
	ConfigureRouter() App
	ListenAndServe()
}

type app struct {
	container service.Container
}

func (a *app) ConfigureRouter() App {
	r := a.container.Router()
	set := a.container.Jet()
	NewGenericRoutes(set).Routes(r)

	return a
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
