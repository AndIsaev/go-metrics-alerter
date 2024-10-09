package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/AndIsaev/go-metrics-alerter/internal/service/agent/metrics"
	"github.com/AndIsaev/go-metrics-alerter/internal/service/agent/utils"

	"github.com/AndIsaev/go-metrics-alerter/internal/common"
	"github.com/AndIsaev/go-metrics-alerter/internal/service/agent"
	"github.com/AndIsaev/go-metrics-alerter/internal/service/agent/middleware"
	"github.com/go-resty/resty/v2"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

// AgentApp - structure of application
type AgentApp struct {
	Config *agent.Config
	Client *resty.Client
	mu     sync.RWMutex
	wg     sync.WaitGroup
	jobs   chan metrics.StorageMetrics
}

func New() *AgentApp {
	app := &AgentApp{}
	config := agent.NewConfig()
	app.Config = config

	return app
}

func (a *AgentApp) StartApp() {
	a.Client = a.initHTTPClient()
	a.jobs = make(chan metrics.StorageMetrics)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	a.runReport(ctx)
}

func (a *AgentApp) initHTTPClient() *resty.Client {
	cli := resty.New()
	cli.SetTimeout(time.Second * 5)
	cli.OnBeforeRequest(middleware.GzipRequestMiddleware)
	cli.OnBeforeRequest(a.HashMiddleware)
	return cli
}

func (a *AgentApp) runReport(ctx context.Context) {
	defer close(a.jobs)
	a.wg.Add(1)
	go a.pullMetrics(ctx)

	a.wg.Add(1)
	go a.collectAdditionalMetrics(ctx)

	a.wg.Add(int(a.Config.RateLimit))
	go a.runWorkers(ctx)

	a.wg.Wait()
}

// pullMetrics - get metrics from runtime and save to StorageMetrics
func (a *AgentApp) pullMetrics(ctx context.Context) {
	defer a.wg.Done()
	for {
		select {
		case <-ctx.Done(): // Завершение по сигналу отмены
			return
		default:
			a.mu.Lock()
			log.Println("pull metrics")
			a.Config.StorageMetrics.Pull()
			log.Println("sent metrics to channel")
			a.mu.Unlock()
			a.jobs <- *a.Config.StorageMetrics
			time.Sleep(a.Config.PollInterval)
		}
	}
}

func (a *AgentApp) runWorkers(ctx context.Context) {
	for w := 1; w <= int(a.Config.RateLimit); w++ {
		go func() {
			defer a.wg.Done()
			for {
				select {
				case <-ctx.Done(): // Завершение по сигналу отмены
					return
				case m, ok := <-a.jobs:
					if !ok {
						log.Printf("jobs channel closed")
						return
					}
					if err := utils.Retry(a.SendMetrics)(m); err != nil {
						log.Printf("Error sending metrics: %v", err)
					}
				}
			}
		}()
	}
}

func (a *AgentApp) SendMetrics(m metrics.StorageMetrics) error {
	values := make([]common.Metrics, 0, len(m.Metrics))
	var result common.Metrics

	for _, v := range m.Metrics {
		metric := common.Metrics{ID: v.ID, MType: v.MType, Value: v.Value, Delta: v.Delta}
		values = append(values, metric)
	}
	log.Println("send metrics")

	res, err := a.Client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(&values).
		SetResult(&result).
		Post(a.Config.UpdateMetricsAddress)

	if err != nil {
		return errors.Unwrap(err)
	}

	if res.StatusCode() != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode())
	}

	return nil
}

func (a *AgentApp) HashMiddleware(c *resty.Client, r *resty.Request) error {

	if a.Config.Key == "" {
		return nil
	}
	value, ok := r.Body.(*[]common.Metrics)
	if !ok {
		log.Printf("not expected type: %T", r.Body)
		return nil
	}
	if value == nil {
		return nil
	}

	v, err := json.Marshal(value)
	if err != nil {
		log.Printf("can't serialize value: %v", err)
		return err
	}

	sha256sum := common.Sha256sum(v, a.Config.Key)

	// Устанавливаем заголовок с хэшем
	c.Header.Set("HashSHA256", sha256sum)
	return nil
}

func (a *AgentApp) collectAdditionalMetrics(ctx context.Context) {
	defer a.wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		default:

			vmStat, err := mem.VirtualMemory()
			if err != nil {
				log.Printf("error getting memory stats: %v", err)
				continue
			}

			totalMemory := float64(vmStat.Total)
			freeMemory := float64(vmStat.Free)

			cpuUtilization, err := cpu.Percent(0, true)
			if err != nil {
				log.Printf("error getting CPU stats: %v", err)
				continue
			}

			a.mu.Lock()

			TotalMemory := metrics.StorageMetric{ID: "TotalMemory", MType: common.Gauge, Value: &totalMemory}
			FreeMemory := metrics.StorageMetric{ID: "FreeMemory", MType: common.Gauge, Value: &freeMemory}
			a.Config.StorageMetrics.AddMetric(TotalMemory)
			a.Config.StorageMetrics.AddMetric(FreeMemory)

			for i, utilization := range cpuUtilization {
				name := fmt.Sprintf("CPUutilization%d", i+1)
				a.Config.StorageMetrics.AddMetric(metrics.StorageMetric{ID: name, MType: common.Gauge, Value: &utilization})
			}
			a.mu.Unlock()

			log.Println("collected additional metrics")
			time.Sleep(a.Config.PollInterval)
		}
	}
}
