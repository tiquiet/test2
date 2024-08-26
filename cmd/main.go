package main

import (
	"context"
	"github.com/joho/godotenv"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sber/internal/config"
	"sber/internal/logger"
	"sber/internal/repository"
	"sber/internal/server"
	"sber/internal/service"
	"syscall"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	log := logger.SetLogger(cfg.Env)

	log.Debug("Debug logging enabled")

	log.Info("Init storage")
	storage := repository.NewRepository(log, cfg.ClearStorageRefreshRate, cfg.FilePath)

	service := service.NewService(log, storage, cfg.MemSize)

	log.Info("Init router")
	router := server.NewRouter(log, service)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.RWTimeout,
		WriteTimeout: cfg.HTTPServer.RWTimeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	log.Info("Starting http server", slog.String("address", srv.Addr))

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("failed to start server", "error", err)
		}
	}()

	log.Info("Server started")

	<-quit
	log.Info("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.HTTPServer.StopTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", "error", err)

		return
	}

	storage.SaveToFile()

	log.Info("server stopped")
}
