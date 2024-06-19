package storage

import (
	"fmt"
	"github.com/AndIsaev/go-metrics-alerter/internal/common"
	"strconv"
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

func (metric *MetricValue) setValue(metricType string, value interface{}) error {
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
	key := MetricKey(metricType + "-" + metricName)

	newMetricValue := &MetricValue{}
	if err := newMetricValue.setValue(metricType, metricValue); err != nil {
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

func (ms *MemStorage) GetV1(metric *common.Metrics) error {
	key := MetricKey(metric.MType + "-" + metric.ID)
	var result interface{}

	if value, ok := ms.Metrics[key]; !ok {
		return ErrKeyErrorStorage
	} else {
		result = value
	}

	switch metric.MType {
	case common.Counter:
		val, _ := strconv.ParseInt(fmt.Sprintf("%v", result), 10, 64)
		metric.Delta = &val
	case common.Gauge:
		val, _ := strconv.ParseFloat(fmt.Sprintf("%v", result), 64)
		metric.Value = &val
	}
	return nil
}

func (ms *MemStorage) Set(metric *common.Metrics) {
	key := MetricKey(metric.MType + "-" + metric.ID)

	switch metric.MType {
	case common.Counter:

		if value, ok := ms.Metrics[key]; !ok {
			ms.Metrics[key] = *metric.Delta
		} else {
			v, e := strconv.ParseInt(fmt.Sprintf("%v", value), 10, 64)
			if e != nil {
				return
			}

			ms.Metrics[key] = *metric.Delta + v
			*metric.Delta = *metric.Delta + v
		}
	case common.Gauge:
		ms.Metrics[key] = *metric.Value
	}

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
