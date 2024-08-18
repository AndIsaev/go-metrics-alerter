package storage

import (
	"fmt"
	"strconv"

	"github.com/AndIsaev/go-metrics-alerter/internal/common"
)

type MetricKey string

type MetricValue struct {
	IntValue   int64
	FloatValue float64
}

type MemStorage struct {
	Metrics map[string]interface{}
}

// NewMemStorage - return new var of MemStorage
func NewMemStorage() *MemStorage {
	return &MemStorage{
		Metrics: make(map[string]interface{}),
	}
}

func (metric *MetricValue) setValue(metricType string, value interface{}) error {
	if value == nil {
		return ErrIncorrectMetricValue
	}

	switch v := value.(type) {
	case int64:
		if metricType == common.Counter {
			metric.IntValue = v
			return nil
		}

	case float64:
		if metricType == common.Gauge {
			metric.FloatValue = v
			return nil
		}
	}
	return ErrIncorrectMetricValue
}

func (ms *MemStorage) Add(metricType, metricName string, metricValue interface{}) error {
	newMetricValue := &MetricValue{}
	if err := newMetricValue.setValue(metricType, metricValue); err != nil {
		return ErrIncorrectMetricValue
	}

	switch metricType {
	case common.Gauge:
		ms.Metrics[metricName] = newMetricValue.FloatValue
		return nil
	case common.Counter:
		val, ok := ms.Metrics[metricName].(int64)
		if ok {
			ms.Metrics[metricName] = val + newMetricValue.IntValue
			return nil
		}
		ms.Metrics[metricName] = metricValue
		return nil
	}
	return ErrIncorrectMetricValue
}

func (ms *MemStorage) GetMetric(MType, ID string) (common.Metrics, error) {
	metric := common.Metrics{ID: ID, MType: MType}

	if value, ok := ms.Metrics[ID]; ok {
		switch MType {
		case common.Counter:
			val, _ := strconv.ParseInt(fmt.Sprintf("%v", value), 10, 64)
			metric.Delta = &val
			return metric, nil
		case common.Gauge:
			val, _ := strconv.ParseFloat(fmt.Sprintf("%v", value), 64)
			metric.Value = &val
			return metric, nil
		}
	}

	return metric, ErrKeyErrorStorage
}

func (ms *MemStorage) Set(metric *common.Metrics) {
	switch metric.MType {
	case common.Counter:

		if value, ok := ms.Metrics[metric.ID]; !ok {
			ms.Metrics[metric.ID] = *metric.Delta
		} else {
			v, e := strconv.ParseInt(fmt.Sprintf("%v", value), 10, 64)
			if e != nil {
				return
			}

			ms.Metrics[metric.ID] = *metric.Delta + v
			*metric.Delta += v
		}
	case common.Gauge:
		ms.Metrics[metric.ID] = *metric.Value
	}
}

func (ms *MemStorage) GetMetricByName(metricName string) (interface{}, error) {
	val, ok := ms.Metrics[metricName]
	if !ok {
		return nil, ErrKeyErrorStorage
	}
	return val, nil
}
