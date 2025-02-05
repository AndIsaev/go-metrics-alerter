package inmemory

import (
	"context"
	"errors"
	"io"
	"log"

	"github.com/AndIsaev/go-metrics-alerter/internal/storage"
)

func (m *MemStorage) Close(_ context.Context) error {
	return nil
}

func (m *MemStorage) Ping(ctx context.Context) error {
	if ctx.Err() != nil {
		log.Println("context is done -> exit from Ping")
		return ctx.Err()
	}
	if m.Metrics == nil {
		log.Println(storage.ErrMapNotAvailable)
		return storage.ErrMapNotAvailable
	}
	return nil
}

func (m *MemStorage) RunMigrations(ctx context.Context) error {
	if ctx.Err() != nil {
		log.Println("context is done -> exit from RunMigrations")
		return ctx.Err()
	}
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
