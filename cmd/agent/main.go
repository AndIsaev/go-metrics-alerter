package main

import (
	"fmt"
	"log"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/AndIsaev/go-metrics-alerter/internal/common"
	"github.com/AndIsaev/go-metrics-alerter/internal/service"
	"github.com/AndIsaev/go-metrics-alerter/internal/service/agent/client"
	"github.com/AndIsaev/go-metrics-alerter/internal/service/agent/metrics"
	"github.com/AndIsaev/go-metrics-alerter/internal/service/agent/middleware"
)

func runPullReport(metrics *metrics.StorageMetrics) {
	metrics.Pull()
}

func SendMetric(url string, metrics *metrics.StorageMetrics) error {
	c := resty.New()
	c.OnBeforeRequest(middleware.GzipRequestMiddleware)

	for _, v := range metrics.Metrics {
		metric := common.Metrics{ID: v.ID, MType: v.MType, Value: v.Value, Delta: v.Delta}
		err := client.SendMetricHandler(c, url, metric)
		if err != nil {
			return err
		}
	}
	return nil
}

func SendMetrics(url string, storage *metrics.StorageMetrics) error {
	c := resty.New()
	c.OnBeforeRequest(middleware.GzipRequestMiddleware)
	values := make([]common.Metrics, 0, 100)

	for _, v := range storage.Metrics {
		metric := common.Metrics{ID: v.ID, MType: v.MType, Value: v.Value, Delta: v.Delta}
		values = append(values, metric)
	}

	if err := client.SendMetricsHandler(c, url, values); err != nil {
		return err
	}
	fmt.Println("=====")
	log.Println(values)
	fmt.Println("=====")
	return nil
}

func main() {
	config := service.NewAgentConfig()
	log.Println("Start Agent")
	for {
		log.Println("Pull Metrics")

		runPullReport(config.StorageMetrics)

		time.Sleep(config.PollInterval)

		log.Println("Send Metrics to Server")

		if err := SendMetric(config.UpdateMetricAddress, config.StorageMetrics); err != nil {
			continue
		}

		if err := SendMetrics(config.UpdateMetricsAddress, config.StorageMetrics); err != nil {
			continue
		}

		time.Sleep(config.ReportInterval)
	}
}
