package in_memory

import (
	"context"
	"errors"
	"github.com/AndIsaev/go-metrics-alerter/internal/storage"
	"io"
	"log"
)

func (m *MemStorage) Close(_ context.Context) error {
	return nil
}

func (m *MemStorage) Ping(_ context.Context) error {
	if m.Metrics == nil {
		log.Fatalln(storage.ErrMapNotAvailable)
	}
	return nil
}

func (m *MemStorage) RunMigrations(ctx context.Context) error {
	log.Println("run migrations")
	metrics, err := m.fm.ReadFile()
	if err != nil && !errors.Is(err, io.EOF) {
		log.Println("error read metrics from file")
		return err
	}
	err = m.InsertBatch(ctx, metrics)
	if err != nil {
		log.Println("error insert batch metrics from file")
		return err
	}
	log.Println("migrations completed")

	return nil
}
