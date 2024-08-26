package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"io"
	"log/slog"
	"net/http"
	"sber/internal/server/api/response"
	"sber/internal/service"
)

type Response struct {
	response.Response
	Key string `json:"Key,omitempty"`
}

func SaveDataHandler(log *slog.Logger, service service.Objects) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		const op = "handlers.save"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		key := chi.URLParam(r, "key")
		expireTime := r.Header.Get("Expires")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Error("failed to parse request body", "error", err)

			render.JSON(w, r, response.Error("failed to parse request"))

			return
		}

		log.Info("Request: ", slog.Any("request", body))

		key, err = service.SaveData(ctx, key, body, expireTime)
		if err != nil {
			log.Error("failed to save data", "error", err)

			render.JSON(w, r, response.Error("failed to save data"))

			return
		}

		render.JSON(w, r, Response{
			Response: response.Ok(),
			Key:      key,
		})
	}
}
