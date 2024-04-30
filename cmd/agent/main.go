package main

import (
	"github.com/AndIsaev/go-metrics-alerter/internal/service/agent/metrics"
	"github.com/AndIsaev/go-metrics-alerter/internal/service/agent/request"

	"fmt"
	"time"
)

const (
	address        string        = "http://localhost:8080/update/%v/%v/%v"
	reportInterval time.Duration = 10
)

func sendReport(m metrics.Metrics) {
	time.Sleep(reportInterval * time.Second)

	for _, v := range m {
		url := fmt.Sprintf(address, v.MetricType, v.Name, v.Value)
		request.SendMetricsHandler(url, "text/plain", nil)
	}
}

func main() {
	newMetrics := metrics.GetMetrics()
	sendReport(newMetrics)
}
