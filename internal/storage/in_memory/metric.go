package in_memory

import (
	"context"
	"github.com/AndIsaev/go-metrics-alerter/internal/common"
	"github.com/AndIsaev/go-metrics-alerter/internal/storage"
)

// UpsertByValue - method create new value or update exists value in map
func (m *MemStorage) UpsertByValue(ctx context.Context, metric common.Metrics, metricValue any) error {
	newValue := storage.MetricValue{}
	if err := newValue.Set(metric.MType, metricValue); err != nil {
		return err
	}

	switch metric.MType {
	case common.Gauge:
		metric.Value = &newValue.FloatValue
		return m.Create(ctx, metric)

	case common.Counter:
		existsMetric, err := m.GetByName(ctx, metric.ID)
		if err != nil {
			metric.Delta = &newValue.IntValue
			return m.Create(ctx, metric)
		}
		newVal := *existsMetric.Delta + newValue.IntValue
		existsMetric.Delta = &newVal
		return m.Create(ctx, existsMetric)
	}

	return storage.ErrMetricValue
}

func (m *MemStorage) List(_ context.Context) ([]common.Metrics, error) {
	var metrics []common.Metrics
	if m.Metrics == nil {
		return nil, storage.ErrMapNotAvailable
	}
	for _, val := range m.Metrics {
		metrics = append(metrics, val)
	}
	return metrics, nil
}

func (m *MemStorage) GetByName(_ context.Context, name string) (common.Metrics, error) {
	metric, ok := m.Metrics[name]
	if !ok {
		return common.Metrics{}, storage.ErrValueNotFound
	}
	return metric, nil
}

func (m *MemStorage) Create(ctx context.Context, metric common.Metrics) error {
	m.Metrics[metric.ID] = metric
	if m.fm != nil && m.syncSave {
		metrics, _ := m.List(ctx)
		err := m.fm.Overwrite(metrics)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *MemStorage) InsertBatch(ctx context.Context, metrics []common.Metrics) error {
	for _, metric := range metrics {
		switch metric.MType {
		case common.Gauge:
			_ = m.Create(ctx, metric)
		case common.Counter:
			existsMetric, err := m.GetByName(ctx, metric.ID)
			if err != nil {
				_ = m.Create(ctx, metric)
				break
			}
			newVal := *existsMetric.Delta + *metric.Delta
			existsMetric.Delta = &newVal
			_ = m.Create(ctx, existsMetric)
		}

	}

	if m.fm != nil && m.syncSave {
		existsValues, _ := m.List(ctx)
		err := m.fm.Overwrite(existsValues)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *MemStorage) GetByNameType(_ context.Context, name, mType string) (common.Metrics, error) {
	metric, ok := m.Metrics[name]
	if !ok || metric.MType != mType {
		return common.Metrics{}, storage.ErrValueNotFound
	}
	return metric, nil
}

func (m *MemStorage) Insert(ctx context.Context, metric common.Metrics) (common.Metrics, error) {
	m.Metrics[metric.ID] = metric
	if m.fm != nil {
		metrics, _ := m.List(ctx)
		err := m.fm.Overwrite(metrics)
		if err != nil {
			return metric, err
		}
	}
	return metric, nil
}
