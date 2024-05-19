package main

import (
	"github.com/AndIsaev/go-metrics-alerter/internal/service"
	"github.com/AndIsaev/go-metrics-alerter/internal/service/agent/client"
	"github.com/AndIsaev/go-metrics-alerter/internal/service/agent/metrics"

	"fmt"
	"time"
)

func sendReport(reportInterval time.Duration, address string, m metrics.Metrics) error {
	time.Sleep(reportInterval)

	for _, v := range m {
		url := fmt.Sprintf("http://%v/update/%v/%v/%v", address, v.MetricType, v.Name, v.Value)
		err := client.SendMetricsClient(url, "text/plain", []byte{})
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	config := service.NewAgentConfig()

	newMetrics := metrics.GetMetrics(config.PollInterval)
	err := sendReport(config.ReportInterval, config.Address, newMetrics)
	if err != nil {
		panic(err)
	}
}
