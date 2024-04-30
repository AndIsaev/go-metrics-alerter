package storage

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
		case "gauge":
			ms.metrics[key] = metricValue
		case "counter":
			ms.metrics[key] = val.(int64) + metricValue.(int64)
		}
	}
}
