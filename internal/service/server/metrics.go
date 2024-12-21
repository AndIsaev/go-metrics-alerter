package server

import (
	"context"
	"github.com/AndIsaev/go-metrics-alerter/internal/common"
	"github.com/AndIsaev/go-metrics-alerter/internal/storage"
	"log"
)

type Service interface {
	PingStorage(ctx context.Context) error
	CloseStorage(ctx context.Context) error
	RunMigrationsStorage(ctx context.Context) error
	ListMetrics(ctx context.Context) ([]common.Metrics, error)
	UpdateMetricByValue(ctx context.Context, metric common.Metrics, value any) error
	GetMetricByName(ctx context.Context, name string) (common.Metrics, error)
	GetMetricByNameType(ctx context.Context, name, mType string) (common.Metrics, error)
	InsertMetric(ctx context.Context, metric common.Metrics) (common.Metrics, error)
	InsertMetrics(ctx context.Context, metrics []common.Metrics) error
}

type Methods struct {
	Storage storage.Storage
}

func (m *Methods) PingStorage(ctx context.Context) error {
	err := m.Storage.System().Ping(ctx)
	if err != nil {
		log.Printf("error ping storage: %v", err.Error())
		return err
	}
	return nil
}

func (m *Methods) CloseStorage(ctx context.Context) error {
	err := m.Storage.System().Close(ctx)
	if err != nil {
		log.Printf("error close storage: %v", err.Error())
		return err
	}
	return nil
}

func (m *Methods) RunMigrationsStorage(ctx context.Context) error {
	err := m.Storage.System().RunMigrations(ctx)
	if err != nil {
		log.Printf("error run migrations for storage: %v", err.Error())
		return err
	}
	return nil
}

func (m *Methods) ListMetrics(ctx context.Context) ([]common.Metrics, error) {
	return m.Storage.Metric().List(ctx)
}

func (m *Methods) UpdateMetricByValue(ctx context.Context, metric common.Metrics, value any) error {
	return m.Storage.Metric().UpsertByValue(ctx, metric, value)
}

func (m *Methods) GetMetricByName(ctx context.Context, name string) (common.Metrics, error) {
	metric, err := m.Storage.Metric().GetByName(ctx, name)
	if err != nil {
		return common.Metrics{}, err
	}
	return metric, nil
}

func (m *Methods) GetMetricByNameType(ctx context.Context, name, mType string) (common.Metrics, error) {
	metric, err := m.Storage.Metric().GetByNameType(ctx, name, mType)
	if err != nil {
		return common.Metrics{}, err
	}
	return metric, nil
}

func (m *Methods) InsertMetric(ctx context.Context, metric common.Metrics) (common.Metrics, error) {
	result, err := m.Storage.Metric().Insert(ctx, metric)
	if err != nil {
		return common.Metrics{}, nil
	}
	return result, nil
}

func (m *Methods) InsertMetrics(ctx context.Context, metrics []common.Metrics) error {
	err := m.Storage.Metric().InsertBatch(ctx, metrics)
	if err != nil {
		return err
	}
	return nil
}
