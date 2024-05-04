package storage

import (
	"errors"
	"github.com/AndIsaev/go-metrics-alerter/internal/common"
)

var ErrIncorrectMetricValue = errors.New("incorrect value for metric")

type MetricKey string

type MetricValue struct {
	IntValue   int64
	FloatValue float64
}

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

type MemStorage struct {
	Metrics map[MetricKey]interface{}
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		Metrics: make(map[MetricKey]interface{}),
	}
}

func (ms *MemStorage) Add(metricType, metricName string, metricValue interface{}) error {
	key := MetricKey(metricType + "/" + metricName)

	newMetricValue := &MetricValue{}
	if err := newMetricValue.SetValue(metricType, metricValue); err != nil {
		println(err)
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
