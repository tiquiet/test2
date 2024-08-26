package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sber/internal/repository"
	"strconv"
)

type DataService struct {
	log   *slog.Logger
	repos repository.JSONStorage
}

const opData = "service.data"

func (s *DataService) GetData(ctx context.Context, key string) (string, error) {
	if key == "" {
		s.log.Error("empty key")
		return "", fmt.Errorf("invalid key")
	}

	data, err := s.repos.Get(ctx, key)
	if err != nil {
		s.log.Error("failed to get data", "error", err)

		return "", fmt.Errorf("failed to get data")
	}

	return data, nil
}

func (s *DataService) SaveData(ctx context.Context, key string, body []byte, expTime string) (string, error) {
	if key == "" {
		s.log.Error("empty key")

		return "", fmt.Errorf("invalid key")
	}

	var js interface{}
	err := json.Unmarshal(body, &js)
	if err != nil {
		s.log.Error("invalid json", "error", err)

		return "", fmt.Errorf("invalid json")
	}

	expireTime, err := strconv.Atoi(expTime)
	if err != nil && expTime != "" {
		expireTime = 0
		s.log.Error("invalid expire time")
	}

	key, err = s.repos.Put(ctx, key, string(body), expireTime)
	if err != nil {
		s.log.Error("failed to store data", "error", "error", err)

		return "", fmt.Errorf("failed to store data")
	}

	return key, nil
}

func NewDataService(log *slog.Logger, repos repository.JSONStorage) *DataService {
	return &DataService{
		log:   log.With(slog.String("op", opData)),
		repos: repos,
	}
}
