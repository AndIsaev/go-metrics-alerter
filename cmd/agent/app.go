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

	"github.com/go-resty/resty/v2"

	"github.com/AndIsaev/go-metrics-alerter/internal/common"
	"github.com/AndIsaev/go-metrics-alerter/internal/service/agent/metrics"
	"github.com/AndIsaev/go-metrics-alerter/internal/service/agent/middleware"
	"github.com/AndIsaev/go-metrics-alerter/internal/service/agent/utils"
)

// AgentApp - structure of application
type AgentApp struct {
	Config *Config
	Client *resty.Client
	mu     sync.RWMutex
	wg     sync.WaitGroup
	jobs   chan metrics.StorageMetrics
}

func New() *AgentApp {
	app := &AgentApp{}
	config := NewConfig()
	app.Config = config

	return app
}

func (a *AgentApp) StartApp() {
	a.Client = a.initHTTPClient()
	a.jobs = make(chan metrics.StorageMetrics)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	a.runReport(ctx)
	defer close(a.jobs)
}

func (a *AgentApp) initHTTPClient() *resty.Client {
	cli := resty.New()
	cli.SetTimeout(time.Second * 5)
	cli.OnBeforeRequest(middleware.GzipRequestMiddleware)
	cli.OnBeforeRequest(a.HashMiddleware)
	return cli
}

func (a *AgentApp) runReport(ctx context.Context) {
	a.wg.Add(2 + int(a.Config.RateLimit))
	go a.pullMetrics(ctx)

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
			log.Println("pull metrics")
			a.mu.Lock()
			a.Config.StorageMetrics.Pull()
			a.mu.Unlock()
			a.jobs <- *a.Config.StorageMetrics
			log.Println("the metrics have been sent to the channel")

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
