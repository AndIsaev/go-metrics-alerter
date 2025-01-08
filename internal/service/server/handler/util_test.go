package handler

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/AndIsaev/go-metrics-alerter/internal/common"
)

func mockedIncorrectMetricTypes() []string {
	testCases := make([]string, 0)
	testCases = append(testCases, "Error")
	testCases = append(testCases, "Some type")

	return testCases
}

func mockedCorrectMetricTypes() []string {
	testCases := make([]string, 0)
	testCases = append(testCases, common.Counter)
	testCases = append(testCases, common.Gauge)

	return testCases
}

func TestIsIncorrectType(t *testing.T) {
	mockedMetrics := mockedIncorrectMetricTypes()
	check := assert.New(t)

	for _, v := range mockedMetrics {
		result := IsCorrectType(v)
		{
			check.False(result)
		}
	}
}

func TestIsCorrectType(t *testing.T) {
	mockedMetrics := mockedCorrectMetricTypes()
	check := assert.New(t)

	for _, v := range mockedMetrics {
		result := IsCorrectType(v)
		{
			check.True(result)
		}
	}
}

func BenchmarkIsCorrectType(b *testing.B) {
	for _, bm := range []struct {
		name       string
		metricName string
	}{
		{"Gauge with valid float", common.Gauge},
		{"Counter with valid int", common.Counter},
		{"Counter with invalid int", "incorrect name"},
	} {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = IsCorrectType(bm.metricName)
			}
		})
	}
}

func TestDefineMetricValue(t *testing.T) {
	tests := []struct {
		metricType  string
		metricValue string
		expected    any
		expectErr   bool
	}{
		{common.Gauge, "123.45", 123.45, false},
		{common.Gauge, "not-a-float", nil, true},
		{common.Counter, "123", int64(123), false},
		{common.Counter, "not-an-int", nil, true},
		{"unknownType", "123", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.metricType+"/"+tt.metricValue, func(t *testing.T) {
			value, err := DefineMetricValue(tt.metricType, tt.metricValue)
			if (err != nil) != tt.expectErr {
				t.Errorf("unexpected error status: got error = %v, want error = %v", err, tt.expectErr)
			}

			if !tt.expectErr && value != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, value)
			}
		})
	}
}

func BenchmarkDefineMetricValue(b *testing.B) {
	for _, bm := range []struct {
		name        string
		metricType  string
		metricValue string
	}{
		{"Gauge with valid float", common.Gauge, "123.45"},
		{"Gauge with invalid float", common.Gauge, "not-a-float"},
		{"Counter with valid int", common.Counter, "123"},
		{"Counter with invalid int", common.Counter, "not-an-int"},
		{"Unknown type", "unknownType", "123"},
	} {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = DefineMetricValue(bm.metricType, bm.metricValue)
			}
		})
	}
}
