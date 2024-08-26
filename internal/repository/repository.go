package repository

import (
	"context"
	"log/slog"
	"sber/internal/repository/storage"
)

//go:generate mockgen -source=repository.go -destination=mocks/mock.go
type JSONStorage interface {
	Get(ctx context.Context, key string) (string, error)
	Put(ctx context.Context, key string, data string, expTime int) (string, error)
	SaveToFile()
	LoadFromFile()
}

type Repository struct {
	JSONStorage
}

func NewRepository(log *slog.Logger, clearTime int, filePath string) *Repository {
	return &Repository{
		JSONStorage: storage.NewStorage(log, clearTime, filePath),
	}
}
