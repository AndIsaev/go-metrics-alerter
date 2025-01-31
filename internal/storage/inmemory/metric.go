package inmemory

import (
	"context"
	"log"

	"github.com/AndIsaev/go-metrics-alerter/internal/common"
	"github.com/AndIsaev/go-metrics-alerter/internal/storage"
)

func (m *MemStorage) UpsertByValue(ctx context.Context, metric common.Metrics, metricValue any) error {
	if ctx.Err() != nil {
		log.Println("context is done -> exit from UpsertByValue")
		return ctx.Err()
	}

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

func (m *MemStorage) List(ctx context.Context) ([]common.Metrics, error) {
	if ctx.Err() != nil {
		log.Println("context is done -> exit from List")
		return []common.Metrics{}, ctx.Err()
	}
	metrics := make([]common.Metrics, 0, len(m.Metrics))

	if m.Metrics == nil {
		return nil, storage.ErrMapNotAvailable
	}
	m.mu.RLock()
	for _, val := range m.Metrics {
		metrics = append(metrics, val)
	}
	m.mu.RUnlock()
	return metrics, nil
}

func (m *MemStorage) GetByName(ctx context.Context, name string) (common.Metrics, error) {
	if ctx.Err() != nil {
		log.Println("context is done -> exit from GetByName")
		return common.Metrics{}, ctx.Err()
	}
	m.mu.RLock()
	metric, ok := m.Metrics[name]
	m.mu.RUnlock()
	if !ok {
		return common.Metrics{}, storage.ErrValueNotFound
	}
	return metric, nil
}
func (m *MemStorage) InsertBatch(ctx context.Context, metrics []common.Metrics) error {
	if ctx.Err() != nil {
		log.Println("context is done -> exit from InsertBatch")
		return ctx.Err()
	}

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

	if err := m.saveMetricsToDisc(ctx); err != nil {
		log.Println("error saving metrics to disk:", err)
		return err
	}

	return nil
}

func (m *MemStorage) GetByNameType(ctx context.Context, name, mType string) (common.Metrics, error) {
	if ctx.Err() != nil {
		log.Println("context is done -> exit from GetByNameType")
		return common.Metrics{}, ctx.Err()
	}
	metric, ok := m.Metrics[name]
	if !ok || metric.MType != mType {
		return common.Metrics{}, storage.ErrValueNotFound
	}
	return metric, nil
}

func (m *MemStorage) Insert(ctx context.Context, metric common.Metrics) (common.Metrics, error) {
	if ctx.Err() != nil {
		log.Println("context is done -> exit from Insert")
		return common.Metrics{}, ctx.Err()
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Metrics[metric.ID] = metric

	err := m.saveMetricsToDisc(ctx)
	if err != nil {
		log.Println("error save metrics to disc")
		return common.Metrics{}, err
	}

	return metric, nil
}

func (m *MemStorage) create(ctx context.Context, metric common.Metrics) error {
	if ctx.Err() != nil {
		log.Println("context is done -> exit from create")
		return ctx.Err()
	}
	m.mu.Lock()
	m.Metrics[metric.ID] = metric
	m.mu.Unlock()
	err := m.saveMetricsToDisc(ctx)
	if err != nil {
		log.Println("error save metrics to disc")
		return err
	}
	return nil
}

func (m *MemStorage) saveMetricsToDisc(ctx context.Context) error {
	if ctx.Err() != nil {
		log.Println("context is done -> exit from saveMetricsToDisc")
		return ctx.Err()
	}

	if m.syncSave {
		metrics, _ := m.List(ctx)
		return m.fm.Overwrite(ctx, metrics)
	}
	return nil
}
