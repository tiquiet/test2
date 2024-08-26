package logger

import (
	"log/slog"
	"os"
)

const (
	dev  = "dev"
	prod = "prod"
)

func SetLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case dev:
		log = slog.New(slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{Level: slog.LevelDebug},
		))
	case prod:
		log = slog.New(slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{Level: slog.LevelInfo},
		))
	default:
		log = slog.New(slog.NewTextHandler(
			os.Stdout,
			&slog.HandlerOptions{Level: slog.LevelInfo},
		))
	}

	return log
}
