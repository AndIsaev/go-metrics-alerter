package main

import (
	"errors"
	"log"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/AndIsaev/go-metrics-alerter/internal/service/agent/metrics"
	"github.com/AndIsaev/go-metrics-alerter/internal/service/agent/utils"
)

// runPullReport - get metrics from runtime and save to StorageMetrics
func runPullReport(metrics *metrics.StorageMetrics) {
	log.Println("pull metrics")
	metrics.Pull()
}

func generator(input *metrics.StorageMetrics, ch chan<- metrics.StorageMetrics) {
	go func() {
		ch <- *input
	}()
}

func worker(ch <-chan metrics.StorageMetrics, wg *sync.WaitGroup, duration time.Duration, f utils.Object) error {
	for {
		time.Sleep(duration)
		select {
		case job := <-ch:
			if err := utils.Retry(f)(job); err != nil {
				wg.Done()
				return err
			}
		default:
			break
		}
	}
}

func main() {
	app := New()
	app.StartApp()
	jobs := make(chan metrics.StorageMetrics)
	g := new(errgroup.Group)

	var (
		mu sync.RWMutex
		wg sync.WaitGroup
	)
	wg.Add(1)

	// pull metrics
	go func() {
		defer wg.Done()
		for {
			time.Sleep(app.Config.PollInterval)
			mu.Lock()
			runPullReport(app.Config.StorageMetrics)
			mu.Unlock()
		}
	}()

	wg.Add(1)

	// send metrics to chan
	go func() {
		defer wg.Done()
		for {
			time.Sleep(app.Config.ReportInterval)
			mu.RLock()
			generator(app.Config.StorageMetrics, jobs)
			mu.RUnlock()
		}
	}()

	// send metrics from chan to server
	for w := 1; w <= int(app.Config.RateLimit); w++ {
		wg.Add(1)
		g.Go(func() error {
			err := worker(jobs, &wg, app.Config.ReportInterval, app.SendMetrics)
			if err != nil {
				return err
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		log.Fatalln("original error -", errors.Unwrap(err))
	}
	wg.Wait()
}
