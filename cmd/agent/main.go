package main

import (
	"encoding/json"
	"github.com/AndIsaev/go-metrics-alerter/internal/service"
	"github.com/AndIsaev/go-metrics-alerter/internal/service/agent/client"
	"github.com/AndIsaev/go-metrics-alerter/internal/service/agent/metrics"

	"fmt"
	"time"
)

func sendReport(reportInterval time.Duration, address string, m metrics.List) error {
	time.Sleep(reportInterval)

	for _, v := range m {
		url := fmt.Sprintf("http://%v/update/", address)
		body, err := json.Marshal(v)

		if err != nil {
			return err
		}

		e := client.SendMetricsClient(url, "application/json", body)
		if e != nil {
			return e
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
