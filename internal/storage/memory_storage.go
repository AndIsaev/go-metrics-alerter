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
		return fmt.Errorf("%w", ErrMetricValue)
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
	return fmt.Errorf("%w", ErrMetricValue)
}

func (ms *MemStorage) Add(metricType, metricName string, metricValue interface{}) error {
	newMetricValue := &MetricValue{}
	if err := newMetricValue.setValue(metricType, metricValue); err != nil {
		return fmt.Errorf("%w", err)
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
	return fmt.Errorf("%w", ErrMetricValue)
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

	return metric, fmt.Errorf("%w", ErrKeyStorage)
}

// Set - insert new value or update exists values
func (ms *MemStorage) Set(metric *common.Metrics) error {
	switch metric.MType {
	case common.Counter:

		if value, ok := ms.Metrics[metric.ID]; !ok {
			ms.Metrics[metric.ID] = *metric.Delta
		} else {
			v, err := strconv.ParseInt(fmt.Sprintf("%v", value), 10, 64)
			if err != nil {
				return err
			}

			ms.Metrics[metric.ID] = *metric.Delta + v
			*metric.Delta += v
		}
	case common.Gauge:
		ms.Metrics[metric.ID] = *metric.Value
	}
	return nil
}

func (ms *MemStorage) GetMetricByName(metricName string) (interface{}, error) {
	val, ok := ms.Metrics[metricName]
	if !ok {
		return nil, fmt.Errorf("%w", ErrKeyStorage)
	}
	return val, nil
}

func (ms *MemStorage) InsertBatch(metrics *[]common.Metrics) error {
	for _, m := range *metrics {
		if err := ms.Set(&m); err != nil {
			return fmt.Errorf("%w", err)
		}
	}
	return nil
}
