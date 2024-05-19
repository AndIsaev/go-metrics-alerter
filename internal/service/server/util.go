package server

import (
	"errors"
	"fmt"
	"github.com/AndIsaev/go-metrics-alerter/internal/common"
	"strconv"
)

// IsCorrectType - check correct type for Metrics
func IsCorrectType(MetricType string) bool {
	for _, v := range []string{common.Counter, common.Gauge} {
		if v == MetricType {

			return true
		}
	}
	return false
}

// DefineMetricValue - define correct value for type of Metrics
func DefineMetricValue(MetricType string, MetricValue string) (interface{}, error) {
	switch MetricType {
	case common.Gauge:
		if val, err := strconv.ParseFloat(MetricValue, 64); err == nil {
			return val, nil
		}
	case common.Counter:
		if val, err := strconv.ParseInt(MetricValue, 10, 64); err == nil {
			return val, nil
		}
	}
	err := fmt.Sprintf("incorrect value for %v type", MetricType)
	return nil, errors.New(err)
}
