package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/AndIsaev/go-metrics-alerter/internal/common"
)

// ErrMapNotAvailable error map not available
var ErrMapNotAvailable = errors.New("map not available")

// ErrMetricValue error incorrect value for metric
var ErrMetricValue = errors.New("incorrect value for metric")

// ErrValueNotFound error value not found
var ErrValueNotFound = errors.New("value not found")

// MetricValue use for define value of metric
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

// Storage provide methods of storage
type Storage interface {
	// System implement system methods
	System() SystemRepository
	// Metric implement metrics methods
	Metric() MetricRepository
}

// SystemRepository system methods
type SystemRepository interface {
	// Close storage connection
	Close(ctx context.Context) error
	// Ping check connection storage
	Ping(ctx context.Context) error
	// RunMigrations run migrations
	RunMigrations(ctx context.Context) error
}

// MetricRepository methods for operations with metrics
type MetricRepository interface {
	// List get list metrics
	List(ctx context.Context) ([]common.Metrics, error)
	// UpsertByValue update metric
	UpsertByValue(ctx context.Context, metric common.Metrics, value any) error
	// InsertBatch insert batch metrics
	InsertBatch(ctx context.Context, metrics []common.Metrics) error
	// GetByName get metric by name
	GetByName(ctx context.Context, name string) (common.Metrics, error)
	// GetByNameType get metric by name and type
	GetByNameType(ctx context.Context, name, mType string) (common.Metrics, error)
	// Insert create new metric
	Insert(ctx context.Context, metric common.Metrics) (common.Metrics, error)
}
