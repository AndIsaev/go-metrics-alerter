package main

import (
	"fmt"
	"github.com/AndIsaev/go-metrics-alerter/internal/common"
	"github.com/AndIsaev/go-metrics-alerter/internal/service"
	"github.com/AndIsaev/go-metrics-alerter/internal/service/agent/client"
	"github.com/AndIsaev/go-metrics-alerter/internal/service/agent/metrics"
	"github.com/go-resty/resty/v2"
	"sync"
	"time"
)

func runPullReport(metrics *metrics.StorageMetrics, interval *time.Ticker) {
	for range interval.C {
		metrics.Pull()
	}
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
	var wg sync.WaitGroup
	wg.Add(2)

	config := service.NewAgentConfig()
	pollInterval := time.NewTicker(config.PollInterval)
	reportInterval := time.NewTicker(config.ReportInterval)

	go runPullReport(config.StorageMetrics, pollInterval)
	go func() {
		for range reportInterval.C {
			err := runSendReport(config.Address, config.StorageMetrics)
			if err != nil {
				panic(err)
			}
		}
	}()
	wg.Wait()

}
