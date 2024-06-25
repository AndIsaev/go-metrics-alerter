package main

import (
	"fmt"
	"github.com/AndIsaev/go-metrics-alerter/internal/common"
	"github.com/AndIsaev/go-metrics-alerter/internal/service"
	"github.com/AndIsaev/go-metrics-alerter/internal/service/agent/client"
	"github.com/AndIsaev/go-metrics-alerter/internal/service/agent/metrics"
	"github.com/go-resty/resty/v2"
	"time"
)

func runPullReport(metrics *metrics.StorageMetrics) {
	metrics.Pull()
}

func runSendReport(address string, metrics *metrics.StorageMetrics) error {
	url := fmt.Sprintf("http://%s/update/", address)
	c := resty.New()

	for _, v := range metrics.Metrics {
		metric := common.Metrics{ID: v.ID, MType: v.MType, Value: &v.Value, Delta: &v.Delta}
		e := client.SendMetricsClient(c, url, metric)
		if e != nil {
			return e
		}
	}
	return nil
}

func main() {
	config := service.NewAgentConfig()
	fmt.Println("Start Agent")
	for {
		fmt.Println("Pull Metrics")

		runPullReport(config.StorageMetrics)

		time.Sleep(config.PollInterval)

		fmt.Println("Send Metrics to Server")

		if err := runSendReport(config.Address, config.StorageMetrics); err != nil {
			continue
		}

		time.Sleep(config.ReportInterval)

	}
}
