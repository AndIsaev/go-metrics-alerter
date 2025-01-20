package common

const (
	Gauge   string = "gauge"
	Counter string = "counter"
)

// Metrics stores metric parameters
//
//go:generate easyjson -all models.go
type Metrics struct {
	// ID name of metric
	ID string `json:"id" db:"id"`
	// MType taking the value gauge or counter
	MType string `json:"type" db:"type"`
	// Delta the value of the metric in the case of a counter transfer
	Delta *int64 `json:"delta,omitempty" db:"delta"`
	// Value the value of the metric in the case of a gauge transfer
	Value *float64 `json:"value,omitempty" db:"value"`
}

// IsValidValue check value of metric
func (m *Metrics) IsValidValue() bool {
	if m.MType == Gauge && m.Value == nil {
		return false
	} else if m.MType == Counter && m.Delta == nil {
		return false
	}
	return true
}

// IsValidType check valid type
func (m *Metrics) IsValidType() bool {
	if m.MType != Gauge && m.MType != Counter {
		return false
	}
	return true
}

type Response struct {
	Message string `json:"message"`
}
