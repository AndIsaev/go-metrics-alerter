package main

import (
	"github.com/AndIsaev/go-metrics-alerter/internal/service/agent/client"
	"github.com/AndIsaev/go-metrics-alerter/internal/service/agent/metrics"

	"fmt"
	"time"
)

const (
	reportInterval time.Duration = 10
)

func sendReport(m metrics.Metrics) error {
	time.Sleep(reportInterval * time.Second)

	for _, v := range m {
		url := fmt.Sprintf("%v/update/%v/%v/%v", flagRunAddr, v.MetricType, v.Name, v.Value)
		err := client.SendMetricsClient(url, "text/plain", []byte{})
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	parseFlags()

	newMetrics := metrics.GetMetrics()
	err := sendReport(newMetrics)
	if err != nil {
		panic(err)
	}
}
