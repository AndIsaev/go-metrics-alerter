package tests

import (
	"github.com/AndIsaev/go-metrics-alerter/internal/service/server"
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
		result := server.IsCorrectType(v)
		{
			check.False(result)
		}
	}
}

func TestIsCorrectType(t *testing.T) {
	mockedMetrics := mockedCorrectMetricTypes()
	check := assert.New(t)

	for _, v := range mockedMetrics {
		result := server.IsCorrectType(v)
		{
			check.True(result)
		}
	}
}
