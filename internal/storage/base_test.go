package storage

import (
	"testing"

	"github.com/AndIsaev/go-metrics-alerter/internal/common"
)

func TestMetricValue_Set(t *testing.T) {
	tests := []struct {
		name       string
		metricType string
		value      any
		wantInt    int64
		wantFloat  float64
		wantErr    bool
	}{
		{
			name:       "Valid Counter",
			metricType: common.Counter,
			value:      int64(5),
			wantInt:    5,
			wantFloat:  0.0,
			wantErr:    false,
		},
		{
			name:       "Valid Gauge",
			metricType: common.Gauge,
			value:      5.5,
			wantInt:    0,
			wantFloat:  5.5,
			wantErr:    false,
		},
		{
			name:       "Invalid MetricType with Int64",
			metricType: common.Gauge,
			value:      5,
			wantErr:    true,
		},
		{
			name:       "Invalid MetricType with Float64",
			metricType: common.Counter,
			value:      5.5,
			wantErr:    true,
		},
		{
			name:       "Unsupported Type",
			metricType: common.Counter,
			value:      "string",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mv := &MetricValue{}
			err := mv.Set(tt.metricType, tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("MetricValue.Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if mv.IntValue != tt.wantInt {
					t.Errorf("MetricValue.Set() IntValue = %v, want %v", mv.IntValue, tt.wantInt)
				}
				if mv.FloatValue != tt.wantFloat {
					t.Errorf("MetricValue.Set() FloatValue = %v, want %v", mv.FloatValue, tt.wantFloat)
				}
			}
		})
	}
}
