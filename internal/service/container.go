package service

import (
	"cdr.dev/slog"
	"github.com/CloudyKit/jet/v6"
	"github.com/CloudyKit/jet/v6/loaders/httpfs"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"grove/internal/config"
	"grove/internal/resource"
	"grove/internal/util"
	"net/http"
)

var _ Container = (*container)(nil)

type Container interface {
	Logger() slog.Logger
	Config() *config.Config
	Router() chi.Router
	Server() *http.Server
	Jet() *jet.Set
}

type container struct {
	logger slog.Logger
	config *config.Config

	router chi.Router
	server *http.Server
	jet    *jet.Set
}

func (c *container) Jet() *jet.Set {
	return c.jet
}

func (c *container) Server() *http.Server {
	return c.server
}

func (c *container) Router() chi.Router {
	return c.router
}

func (c *container) Logger() slog.Logger {
	return c.logger
}

func (c *container) Config() *config.Config {
	return c.config
}

func New(logger slog.Logger, config *config.Config) Container {
	router := createRouter(logger)

	loader, err := httpfs.NewLoader(resource.Templates)
	if err != nil {
		panic(err)
	}

	return &container{
		logger: logger,
		config: config,
		router: router,
		server: &http.Server{
			Addr:    config.Http.Listen,
			Handler: router,
		},
		jet: jet.NewSet(loader, jet.InDevelopmentMode()), // development mode for templates ignores cache - which is not a performace issue in production (evereything is loaded from memory)
	}
}

func createRouter(logger slog.Logger) *chi.Mux {
	router := chi.NewRouter()
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
		middleware.RequestLogger(util.NewSlogChiFormatter(logger.Named("http"))),

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
