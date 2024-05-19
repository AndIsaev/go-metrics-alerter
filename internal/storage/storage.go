package storage

import (
	"github.com/AndIsaev/go-metrics-alerter/internal/common"
)

type MetricKey string

type MetricValue struct {
	IntValue   int64
	FloatValue float64
}

type MemStorage struct {
	Metrics map[MetricKey]interface{}
}

// NewMemStorage - return new var of MemStorage
func NewMemStorage() *MemStorage {
	return &MemStorage{
		Metrics: make(map[MetricKey]interface{}),
	}
}

var MS = NewMemStorage()

func (metric *MetricValue) SetValue(metricType string, value interface{}) error {
	if value == nil {
		return ErrIncorrectMetricValue
	}

	switch v := value.(type) {
	case int64:
		if metricType == common.Counter {
			metric.IntValue = v
			return nil
		} else {
			return ErrIncorrectMetricValue
		}

	case float64:
		if metricType == common.Gauge {
			metric.FloatValue = v
			return nil
		} else {
			return ErrIncorrectMetricValue
		}
	}
	return ErrIncorrectMetricValue
}

func (ms *MemStorage) Add(metricType, metricName string, metricValue interface{}) error {
	key := MetricKey(metricName)

	newMetricValue := &MetricValue{}
	if err := newMetricValue.SetValue(metricType, metricValue); err != nil {
		return ErrIncorrectMetricValue
	}

	switch metricType {
	case common.Gauge:
		ms.Metrics[key] = newMetricValue.FloatValue
		return nil
	case common.Counter:
		if val, ok := ms.Metrics[key].(int64); ok {
			ms.Metrics[key] = val + newMetricValue.IntValue
			return nil
		} else {
			ms.Metrics[key] = metricValue
			return nil
		}
	}
	return ErrIncorrectMetricValue
}

func (ms *MemStorage) Ping() error {
	if err := ms.Metrics; err == nil {
		return ErrNotInitializedStorage
	}
	return nil
}

func (ms *MemStorage) Get(metricName string) (interface{}, error) {
	if err := ms.Ping(); err != nil {
		return nil, err
	}
	key := MetricKey(metricName)

	if val, ok := ms.Metrics[key]; !ok {
		return nil, ErrKeyErrorStorage
	} else {
		return val, nil
	}
}
