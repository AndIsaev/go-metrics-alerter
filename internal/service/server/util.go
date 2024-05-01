package server

import (
	"github.com/AndIsaev/go-metrics-alerter/internal/common"
	"strconv"
)

func IsCorrectType(metricType string) bool {
	for _, v := range []string{common.Counter, common.Gauge} {
		if v == metricType {

			return true
		}
	}
	return false
}

func DefineMetricValue(metricType string, metricValue string) interface{} {
	switch metricType {
	case common.Gauge:
		val, err := strconv.ParseFloat(metricValue, 64)
		if err == nil {
			return val
		}
	case common.Counter:
		val, err := strconv.ParseInt(metricValue, 10, 64)
		if err == nil {
			return val
		}
	}
	return nil
}
