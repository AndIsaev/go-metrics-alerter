package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPull(t *testing.T) {
	storageMetrics := NewListMetrics()

	assert.Empty(t, storageMetrics.Metrics)

	storageMetrics.Pull()

	assert.NotEmpty(t, storageMetrics.Metrics)

	expectedMetrics := []string{"TotalMemory", "FreeMemory", "PollCount", "Alloc", "RandomValue"}
	for _, metricName := range expectedMetrics {
		metric, exists := storageMetrics.Metrics[metricName]
		assert.True(t, exists, "Metric %s should exist", metricName)
		switch metricName {
		case "PollCount":
			assert.NotNil(t, metric.Delta, "Metric value for %s should not be nil", metricName)
			assert.NotNil(t, metric.Delta)
			assert.Equal(t, int64(1), *metric.Delta) // First call, should be 1
		default:
			assert.NotNil(t, metric.Value, "Metric value for %s should not be nil", metricName)
		}
	}
}
