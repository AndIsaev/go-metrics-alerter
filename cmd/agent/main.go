package main

import (
	"errors"
	"log"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/AndIsaev/go-metrics-alerter/internal/common"
	"github.com/AndIsaev/go-metrics-alerter/internal/service/agent"
	"github.com/AndIsaev/go-metrics-alerter/internal/service/agent/client"
	"github.com/AndIsaev/go-metrics-alerter/internal/service/agent/metrics"
	"github.com/AndIsaev/go-metrics-alerter/internal/service/agent/middleware"
	"github.com/AndIsaev/go-metrics-alerter/internal/service/agent/utils"
)

func runPullReport(metrics *metrics.StorageMetrics) {
	log.Println("pull metrics")
	metrics.Pull()
}

func SendMetrics(url string, db *metrics.StorageMetrics) error {
	c := resty.New().SetTimeout(time.Second * 5)
	c.OnBeforeRequest(middleware.GzipRequestMiddleware)
	values := make([]common.Metrics, 0, 100)

	for _, v := range db.Metrics {
		metric := common.Metrics{ID: v.ID, MType: v.MType, Value: v.Value, Delta: v.Delta}
		values = append(values, metric)
	}

	if err := client.SendMetricsHandler(c, url, &values); err != nil {
		return errors.Unwrap(err)
	}

	return nil
}

func main() {
	config := agent.NewConfig()
	log.Println("start agent")
	for {
		runPullReport(config.StorageMetrics)

		time.Sleep(config.PollInterval)

		if err := utils.Retry(SendMetrics)(config.UpdateMetricsAddress, config.StorageMetrics); err != nil {
			log.Fatalln("original error -", errors.Unwrap(err))
		}

		time.Sleep(config.ReportInterval)
	}
}
