package main

import (
	"errors"
	"log"
	"time"

	"github.com/AndIsaev/go-metrics-alerter/internal/service/agent/metrics"
	"github.com/AndIsaev/go-metrics-alerter/internal/service/agent/utils"
)

func runPullReport(metrics *metrics.StorageMetrics) {
	log.Println("pull metrics")
	metrics.Pull()
}

func main() {
	app := New()
	app.StartApp()

	log.Println("start agent")
	for {
		runPullReport(app.Config.StorageMetrics)

		time.Sleep(app.Config.PollInterval)

		if err := utils.Retry(app.SendMetrics)(); err != nil {
			log.Fatalln("original error -", errors.Unwrap(err))
		}

		time.Sleep(app.Config.ReportInterval)
	}
}
