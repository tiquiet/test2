package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"sber/internal/service"
)

func LivenessHandler(log *slog.Logger, service service.Probes) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.liveness"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		err := service.Liveness()
		if err != nil {
			log.Error("Service not healthy", "error", err)
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintf(w, "Not alive\n")

			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Alive\n")
	}
}
