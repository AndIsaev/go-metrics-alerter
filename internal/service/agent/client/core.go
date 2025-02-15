package client

import "github.com/AndIsaev/go-metrics-alerter/internal/common"

// RequestClient use for sending metrics
type RequestClient interface {
	// SendMetrics use for send batch metrics
	SendMetrics(metrics []common.Metrics) error
}
