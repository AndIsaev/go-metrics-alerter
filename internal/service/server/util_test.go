package server

import (
	"github.com/AndIsaev/go-metrics-alerter/internal/common"
	"github.com/stretchr/testify/assert"
	"testing"
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
