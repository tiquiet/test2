package service

import (
	"fmt"
	"log/slog"
	"runtime"
)

const (
	opProbes     = "service.data"
	mbMultiplier = 1024 * 1024
)

type ProbesService struct {
	log     *slog.Logger
	MemSize uint64
}

func (s *ProbesService) Liveness() error {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	if m.HeapAlloc > s.MemSize*mbMultiplier {
		s.log.Error("high memory usage")
		return fmt.Errorf("high memory usage")
	}

	return nil
}

func (s *ProbesService) Readiness() error {
	// Здесь можно добавить пинг в бд
	// Не захотел ничего лишнего здесь кастылить поэтому просто вернул nil
	return nil
}

func NewProbesService(log *slog.Logger, memSize uint64) *ProbesService {
	return &ProbesService{
		log:     log.With(slog.String("op", opProbes)),
		MemSize: memSize,
	}
}
