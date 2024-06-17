package main

import (
	"fmt"
	"github.com/AndIsaev/go-metrics-alerter/internal/service"
	"github.com/AndIsaev/go-metrics-alerter/internal/service/agent/client"
	"github.com/AndIsaev/go-metrics-alerter/internal/service/agent/metrics"
	"github.com/go-resty/resty/v2"
	"time"
)

func sendReport(reportInterval time.Duration, address string, m metrics.List) error {
	time.Sleep(reportInterval)
	url := fmt.Sprintf("http://%v/update/", address)
	c := resty.New()

	for _, v := range m {
		e := client.SendMetricsClient(c, url, v)
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
