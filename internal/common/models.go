package common

//go:generate easyjson -all models.go
type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func (metrics *Metrics) IsValid() bool {
	if !metrics.IsValidType() {
		return false
	}
	if metrics.MType == Gauge && metrics.Value == nil {
		return false
	} else if metrics.MType == Counter && metrics.Delta == nil {
		return false
	}
	return true
}

func (metrics *Metrics) IsValidType() bool {
	if metrics.MType != Gauge && metrics.MType != Counter {
		return false
	}
	return true
}

type Response struct {
	Message string `json:"message"`
}
