package storage

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"os"
	"testing"
	"time"

	"log/slog"
)

func TestPutAndGet(t *testing.T) {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	s := NewStorage(log, 60, "testfile.json")
	defer os.Remove("testfile.json")

	key, err := s.Put(context.Background(), "key1", `{"key":"value"}`, 0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if key != "key1" {
		t.Fatalf("expected key to be 'key1', got %v", key)
	}

	value, err := s.Get(context.Background(), "key1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if value != `{"key":"value"}` {
		t.Fatalf("expected value to be '%v', got '%v'", `{"key":"value"}`, value)
	}
}

func TestPutDuplicate(t *testing.T) {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	s := NewStorage(log, 60, "testfile.json")
	defer os.Remove("testfile.json")

	_, err := s.Put(context.Background(), "key1", `{"key":"value"}`, 0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = s.Put(context.Background(), "key1", `{"key":"value2"}`, 0)
	if err == nil {
		t.Fatalf("expected error for duplicate key, got nil")
	}
}

func TestExpiration(t *testing.T) {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	s := NewStorage(log, 1, "testfile.json")
	defer os.Remove("testfile.json")

	_, err := s.Put(context.Background(), "key1", `{"key":"value"}`, 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	time.Sleep(2 * time.Second)

	_, err = s.Get(context.Background(), "key1")
	if err == nil {
		t.Fatalf("expected error for expired object, got nil")
	}
}

func TestSaveAndLoadFromFile(t *testing.T) {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	s := NewStorage(log, 60, "testfile.json")
	defer os.Remove("testfile.json")

	_, err := s.Put(context.Background(), "key1", `{"key":"value"}`, 0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	s.SaveToFile()

	prometheus.Unregister(s.size)

	s2 := NewStorage(log, 60, "testfile.json")
	s2.LoadFromFile()

	value, err := s2.Get(context.Background(), "key1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if value != `{"key":"value"}` {
		t.Fatalf("expected value to be '%v', got '%v'", `{"key":"value"}`, value)
	}
}
