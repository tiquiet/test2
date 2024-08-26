package storage

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type Object struct {
	Data      []byte     `json:"Data"`
	ExpiresAt *time.Time `json:"ExpiresAt"`
}

type Storage struct {
	log      *slog.Logger
	mu       sync.RWMutex
	data     map[string]Object
	size     prometheus.Gauge
	filePath string
}

const op = "repository.storage"

func (s *Storage) Put(ctx context.Context, key string, json string, expTime int) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data[key]; ok {
		return "", ErrorAlreadyExists
	}

	obj := Object{
		Data: []byte(json),
	}

	if expTime > 0 {
		expiresAt := time.Now().Add(time.Second * time.Duration(expTime))
		obj.ExpiresAt = &expiresAt
	}

	s.size.Inc()
	s.data[key] = obj

	return key, nil
}

func (s *Storage) Get(ctx context.Context, key string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, ok := s.data[key]
	if !ok {
		return "", ErrorNotFound
	}

	return string(value.Data), nil
}

func (s *Storage) removeExpiredObjects(clearTime int) {
	for {
		time.Sleep(time.Duration(clearTime) * time.Second)
		s.mu.Lock()
		for key, obj := range s.data {
			if obj.ExpiresAt != nil && time.Now().After(*obj.ExpiresAt) {
				delete(s.data, key)
				s.size.Dec()
			}
		}
		s.mu.Unlock()
	}
}

func (s *Storage) SaveToFile() {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := json.Marshal(s.data)
	if err != nil {
		s.log.Error("Error saving to file:", "error", err)
		return
	}

	err = os.WriteFile(s.filePath, data, 0644)
	if err != nil {
		s.log.Error("Error writing file:", "error", err)
	}
}

func (s *Storage) LoadFromFile() {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.filePath)
	if err != nil {

		s.log.Error("No storage file found, starting fresh.")
		return
	}

	err = json.Unmarshal(data, &s.data)
	if err != nil {
		s.log.Error("Error reading file:", "error", err)
		return
	}

	s.size.Set(float64(len(s.data)))
}

func NewStorage(log *slog.Logger, clearTime int, filePath string) *Storage {
	st := &Storage{
		log:      log.With(slog.String("op", op)),
		data:     make(map[string]Object),
		filePath: filePath,
		size: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "object_count",
			Help: "Number of objects in storage",
		}),
	}

	prometheus.MustRegister(st.size)

	st.LoadFromFile()

	go st.removeExpiredObjects(clearTime)

	return st
}
