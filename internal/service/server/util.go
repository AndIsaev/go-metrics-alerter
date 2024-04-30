package server

import (
	"github.com/AndIsaev/go-metrics-alerter/internal/common"
)

func IsCorrectType(metricType string) bool {
	for _, v := range []string{common.Counter, common.Gauge} {
		if v == metricType {
			return true
		}
	}
	return false

}
