package inmemory

import (
	"context"
	"log"

	"github.com/AndIsaev/go-metrics-alerter/internal/common"
	"github.com/AndIsaev/go-metrics-alerter/internal/storage"
)

func (m *MemStorage) UpsertByValue(ctx context.Context, metric common.Metrics, metricValue any) error {
	newValue := storage.MetricValue{}
	if err := newValue.Set(metric.MType, metricValue); err != nil {
		return err
	}

	switch metric.MType {
	case common.Gauge:
		metric.Value = &newValue.FloatValue
		return m.create(ctx, metric)

	case common.Counter:
		existsMetric, err := m.GetByName(ctx, metric.ID)
		if err != nil {
			metric.Delta = &newValue.IntValue
			return m.create(ctx, metric)
		}
		newVal := *existsMetric.Delta + newValue.IntValue
		existsMetric.Delta = &newVal
		return m.create(ctx, existsMetric)
	}

	return storage.ErrMetricValue
}

func (m *MemStorage) List(_ context.Context) ([]common.Metrics, error) {
	metrics := make([]common.Metrics, 0, len(m.Metrics))
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

func (m *MemStorage) InsertBatch(ctx context.Context, metrics []common.Metrics) error {
	for _, metric := range metrics {
		switch metric.MType {
		case common.Gauge:
			_ = m.create(ctx, metric)
		case common.Counter:
			existsMetric, err := m.GetByName(ctx, metric.ID)
			if err != nil {
				_ = m.create(ctx, metric)
				break
			}
			newVal := *existsMetric.Delta + *metric.Delta
			existsMetric.Delta = &newVal
			_ = m.create(ctx, existsMetric)
		}
	}

	err := m.saveMetricsToDisc(ctx)
	if err != nil {
		log.Println("error save metrics to disc")
		return err
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

	err := m.saveMetricsToDisc(ctx)
	if err != nil {
		log.Println("error save metrics to disc")
		return common.Metrics{}, err
	}

	return metric, nil
}

func (m *MemStorage) create(ctx context.Context, metric common.Metrics) error {
	m.Metrics[metric.ID] = metric
	err := m.saveMetricsToDisc(ctx)
	if err != nil {
		log.Println("error save metrics to disc")
		return err
	}
	return nil
}

func (m *MemStorage) saveMetricsToDisc(ctx context.Context) error {
	if m.syncSave {
		metrics, _ := m.List(ctx)
		return m.fm.Overwrite(metrics)
	}
	return nil
}
