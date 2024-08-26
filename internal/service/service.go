package service

import (
	"context"
	"log/slog"
	"sber/internal/repository"
)

type Objects interface {
	GetData(ctx context.Context, key string) (string, error)
	SaveData(ctx context.Context, key string, body []byte, expTime string) (string, error)
}

type Probes interface {
	Liveness() error
	Readiness() error
}

type Service struct {
	Objects
	Probes
}

func NewService(log *slog.Logger, repos *repository.Repository, memSize uint64) *Service {
	return &Service{
		Objects: NewDataService(log, repos.JSONStorage),
		Probes:  NewProbesService(log, memSize),
	}
}
