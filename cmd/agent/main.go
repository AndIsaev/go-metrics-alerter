package main

import (
	"encoding/json"
	"fmt"
	"github.com/AndIsaev/go-metrics-alerter/internal/service"
	"github.com/AndIsaev/go-metrics-alerter/internal/service/agent/client"
	"github.com/AndIsaev/go-metrics-alerter/internal/service/agent/metrics"
	"github.com/go-resty/resty/v2"
	"time"
)

func sendReport(reportInterval time.Duration, address string, m metrics.List) error {
	time.Sleep(reportInterval)
	c := resty.New()

	for _, v := range m {
		url := fmt.Sprintf("http://%v/update/", address)
		body, err := json.Marshal(v)

		if err != nil {
			return err
		}

		e := client.SendMetricsClient(c, url, body)
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

	fmt.Println("---------------------------------------------------")
	fmt.Println(err)
	fmt.Println("---------------------------------------------------")

	if err != nil {
		panic(err)
	}
}
