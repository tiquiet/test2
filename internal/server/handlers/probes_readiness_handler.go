package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"sber/internal/service"
)

func ReadinessHandler(log *slog.Logger, service service.Probes) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.readiness"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		err := service.Readiness()
		if err != nil {
			log.Error("service not ready", "error", err)
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintf(w, "Not ready\n")

			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Ready\n")
	}
}
