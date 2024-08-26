package handlers

import (
	"log/slog"
	"net/http"
	"sber/internal/server/api/response"
	"sber/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type GetDataResponse struct {
	response.Response
	Data string `json:"data,omitempty"`
}

func GetDataHandler(log *slog.Logger, service service.Objects) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		const op = "handlers.load"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		key := chi.URLParam(r, "key")

		data, err := service.GetData(ctx, key)
		if err != nil {
			log.Error("failed to get data", "error", err)

			render.JSON(w, r, response.Error("failed to get data"))

			return
		}

		render.JSON(w, r, GetDataResponse{
			Response: response.Ok(),
			Data:     data,
		})
	}
}
