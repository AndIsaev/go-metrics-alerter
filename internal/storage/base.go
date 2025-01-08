package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/AndIsaev/go-metrics-alerter/internal/common"
)

var ErrMapNotAvailable = errors.New("map not available")
var ErrMetricValue = errors.New("incorrect value for metric")
var ErrValueNotFound = errors.New("value not found")

type MetricValue struct {
	IntValue   int64
	FloatValue float64
}

// Set - set value for key in map
func (mv *MetricValue) Set(metricType string, value any) error {
	switch v := value.(type) {
	case int64:
		if metricType == common.Counter {
			mv.IntValue = v
			return nil
		}

	case float64:
		if metricType == common.Gauge {
			mv.FloatValue = v
			return nil
		}
	}
	return fmt.Errorf("%w", ErrMetricValue)
}

type Storage interface {
	System() SystemRepository
	Metric() MetricRepository
}

type SystemRepository interface {
	Close(ctx context.Context) error
	Ping(ctx context.Context) error
	RunMigrations(ctx context.Context) error
}

type MetricRepository interface {
	List(ctx context.Context) ([]common.Metrics, error)
	UpsertByValue(ctx context.Context, metric common.Metrics, value any) error
	Create(ctx context.Context, metric common.Metrics) error
	InsertBatch(ctx context.Context, metrics []common.Metrics) error
	GetByName(ctx context.Context, name string) (common.Metrics, error)
	GetByNameType(ctx context.Context, name, mType string) (common.Metrics, error)
	Insert(ctx context.Context, metric common.Metrics) (common.Metrics, error)
}
