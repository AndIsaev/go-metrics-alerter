package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidValue(t *testing.T) {
	var delta int64 = 100
	var value float64 = 10.5

	tests := []struct {
		name     string
		metrics  Metrics
		expected bool
	}{
		{
			name: "Valid gauge with value",
			metrics: Metrics{
				ID:    "metric1",
				MType: Gauge,
				Value: &value,
			},
			expected: true,
		},
		{
			name: "Invalid gauge without value",
			metrics: Metrics{
				ID:    "metric2",
				MType: Gauge,
				Value: nil,
			},
			expected: false,
		},
		{
			name: "Valid counter with delta",
			metrics: Metrics{
				ID:    "metric3",
				MType: Counter,
				Delta: &delta,
			},
			expected: true,
		},
		{
			name: "Invalid counter without delta",
			metrics: Metrics{
				ID:    "metric4",
				MType: Counter,
				Delta: nil,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.metrics.IsValidValue()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsValidType(t *testing.T) {
	tests := []struct {
		name     string
		metrics  Metrics
		expected bool
	}{
		{
			name: "Valid gauge type",
			metrics: Metrics{
				ID:    "metric1",
				MType: Gauge,
			},
			expected: true,
		},
		{
			name: "Valid counter type",
			metrics: Metrics{
				ID:    "metric2",
				MType: Counter,
			},
			expected: true,
		},
		{
			name: "Invalid type",
			metrics: Metrics{
				ID:    "metric3",
				MType: "invalidType",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.metrics.IsValidType()
			assert.Equal(t, tt.expected, result)
		})
	}
}
