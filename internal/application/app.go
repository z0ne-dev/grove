package application

import (
	"cdr.dev/slog"
	"context"
	"grove/internal/service"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var _ App = (*app)(nil)

type App interface {
	ListenAndServe()
}

type app struct {
	container service.Container
}

func (a *app) ListenAndServe() {
	s := a.container.Server()
	configuration := a.container.Config()
	slogger := a.container.Logger().Named("server")

	errChan := make(chan error, 1)
	stopChan := make(chan os.Signal, 1)

	go func(h *http.Server, errChan chan error) {
		errChan <- h.ListenAndServe()
	}(s, errChan)
	defer func(s *http.Server) {
		_ = s.Close() // No need to handle error here
	}(s)

	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	slogger.Info(
		context.Background(),
		"server started",
		slog.F("addr", configuration.Http.PublicAddress),
	)

	select {
	case err := <-errChan:
		if err != nil {
			slogger.Fatal(
				context.Background(),
				"server returned with an error",
				slog.Error(err),
			)
		}
	case sig := <-stopChan:
		slogger.Debug(
			context.Background(),
			"received stop signal",
			slog.F("signal", sig),
		)
	}
}

func New(container service.Container) App {
	return &app{
		container: container,
	}
}

// TODO: setup routes
