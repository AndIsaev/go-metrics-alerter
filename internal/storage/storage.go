package storage

import (
	"github.com/AndIsaev/go-metrics-alerter/internal/common"
)

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
	switch metricType {
	case common.Gauge:
		ms.metrics[key] = metricValue
	case common.Counter:
		if newVal, res := metricValue.(int64); res {

			if val, ok := ms.metrics[key].(int64); ok {
				ms.metrics[key] = val + newVal
			} else {
				ms.metrics[key] = metricValue
			}
		}
	}
}
