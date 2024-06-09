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

func (ms *MemStorage) GetV1(metricName string) (interface{}, error) {
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

func (ms *MemStorage) Set(metric *common.Metrics) {
	key := MetricKey(metric.MType + "-" + metric.ID)

	switch metric.MType {
	case common.Counter:

		if value, ok := ms.Metrics[key]; !ok {
			ms.Metrics[key] = *metric.Delta
		} else {
			var result int64

			str := fmt.Sprintf("%v", value)
			v, e := strconv.ParseInt(str, 10, 32)
			if e != nil {
				return
			}
			result = *metric.Delta + v

			ms.Metrics[key] = result
			metric.Delta = &result
		}
	case common.Gauge:
		ms.Metrics[key] = *metric.Value
	}

}

func (ms *MemStorage) Get(metric *common.Metrics) error {
	key := MetricKey(metric.MType + "-" + metric.ID)
	if _, ok := ms.Metrics[key]; !ok {
		return ErrKeyErrorStorage
	}

	switch metric.MType {
	case common.Gauge:
		var result float64
		result, err := strconv.ParseFloat(fmt.Sprintf("%v", ms.Metrics[key]), 64)
		if err != nil {
			return err
		}

		metric.Value = &result
	case common.Counter:
		var result int64
		result, err := strconv.ParseInt(fmt.Sprintf("%v", ms.Metrics[key]), 10, 64)
		if err != nil {
			return err
		}
		metric.Delta = &result
	}
	return nil
}
