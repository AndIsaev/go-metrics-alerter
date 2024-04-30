package storage

import "github.com/AndIsaev/go-metrics-alerter/internal/common"

type MemStorage struct {
	metrics map[string]interface{}
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		metrics: make(map[string]interface{}),
	}
}

func (ms *MemStorage) Update(metricType, metricName string, metricValue interface{}) {
	key := metricType + "/" + metricName
	if val, ok := ms.metrics[key]; ok {
		switch metricType {
		case common.Gauge:
			ms.metrics[key] = metricValue
		case common.Counter:
			ms.metrics[key] = val.(int64) + metricValue.(int64)
		}
	}
}
