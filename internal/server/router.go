package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log/slog"
	"sber/internal/server/handlers"
	"sber/internal/service"
)

func NewRouter(log *slog.Logger, service *service.Service) chi.Router {
	router := chi.NewRouter()

	// middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/objects", func(router chi.Router) {
		router.Route("/{key}", func(router chi.Router) {
			router.Put("/", handlers.SaveDataHandler(log, service.Objects))
			router.Get("/", handlers.GetDataHandler(log, service.Objects))
		})
	})

	router.Route("/probes", func(router chi.Router) {
		router.Put("/liveness", handlers.LivenessHandler(log, service.Probes))
		router.Get("/readiness", handlers.ReadinessHandler(log, service.Probes))
	})

	router.Handle("/metrics", promhttp.Handler())

	return router
}
