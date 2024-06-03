package main

import (
	"encoding/json"
	"github.com/AndIsaev/go-metrics-alerter/internal/common/models"
	"github.com/AndIsaev/go-metrics-alerter/internal/service"
	"github.com/AndIsaev/go-metrics-alerter/internal/service/agent/client"

	"github.com/AndIsaev/go-metrics-alerter/internal/service/agent/metrics"

	"fmt"
	"time"
)

func sendReport(reportInterval time.Duration, address string, m []models.Metrics) error {
	time.Sleep(reportInterval)

	for _, v := range m {
		if body, err := json.Marshal(v); err == nil {

			url := fmt.Sprintf("http://%v/update/", address)
			if e := client.SendMetricsClient(url, "application/json", body); e != nil {
				return e
			}

		} else {
			return err
		}

	}
	return nil
}

func main() {
	config := service.NewAgentConfig()

	newMetrics := metrics.GetMetrics(config.PollInterval)
	m := metrics.ConvertMetrics(newMetrics)

	err := sendReport(config.ReportInterval, config.Address, m)
	if err != nil {
		panic(err)
	}
}
