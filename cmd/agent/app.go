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

	"github.com/AndIsaev/go-metrics-alerter/internal/service/agent/utils"

	"github.com/go-resty/resty/v2"

	"github.com/AndIsaev/go-metrics-alerter/internal/common"
)

// AgentApp structure of application
type AgentApp struct {
	// Config use for settings of app
	Config *Config
	// Client use for requests to server
	Client *resty.Client
	mu     sync.RWMutex
	wg     sync.WaitGroup
	jobs   chan common.Metrics
}

// New create and return new AgentApp
func New() *AgentApp {
	app := &AgentApp{}
	config := NewConfig()
	app.Config = config

	return app
}

// StartApp user for start application
func (a *AgentApp) StartApp() {
	a.Client = a.initHTTPClient()
	a.jobs = make(chan common.Metrics)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	a.runReport(ctx)
	defer close(a.jobs)
}

func (a *AgentApp) initHTTPClient() *resty.Client {
	cli := resty.New()
	cli.SetTimeout(time.Second * 5)
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
		case <-ctx.Done():
			return
		default:
			log.Println("pull metrics")

			a.Config.StorageMetrics.Pull()

			a.mu.RLock()
			for _, val := range a.Config.StorageMetrics.Metrics {
				a.jobs <- val
			}
			a.mu.RUnlock()
			log.Println("the metrics have been sent to the channel")

			time.Sleep(a.Config.PollInterval)
		}
	}
}

func (a *AgentApp) runWorkers(ctx context.Context) {
	a.wg.Add(int(a.Config.RateLimit)) // Увеличивайте на количество горутин

	for w := 1; w <= int(a.Config.RateLimit); w++ {
		go func() {
			defer a.wg.Done()
			for {
				select {
				case <-ctx.Done(): // Завершение по сигналу отмены
					return
				default:
					tasks := make([]common.Metrics, 0, 3)
					for i := 0; i < 3; i++ {
						metric, ok := <-a.jobs
						if !ok {
							log.Printf("jobs channel closed")
							return
						}
						tasks = append(tasks, metric)
					}
					if len(tasks) > 0 {
						if err := utils.Retry(a.sendMetrics)(tasks); err != nil {
							log.Printf("Error sending metrics: %v", err)
						}
					}

					interval := a.Config.ReportInterval
					time.Sleep(interval)
				}
			}
		}()
	}
}

func (a *AgentApp) sendMetrics(metrics []common.Metrics) error {
	body, err := json.Marshal(metrics)
	if err != nil {
		return errors.Unwrap(fmt.Errorf("error encoding metric: %w", err))
	}

	client := a.Client.R()
	if a.Config.Key != "" {
		sha256sum := common.Sha256sum(body, a.Config.Key)
		client.SetHeader("HashSHA256", sha256sum)
	}

	if a.Config.PublicKeyPath != "" {
		body, err = utils.Encrypt(a.Config.PublicKey, body)
		if err != nil {
			return fmt.Errorf("error encrypting metrics: %w", errors.Unwrap(err))
		}
	}

	log.Println("send metrics")

	// Используйте правильный Content-Type
	res, err := client.
		SetHeader("Accept-Encoding", "gzip").
		SetHeader("Content-Type", "application/octet-stream").
		SetBody(body).
		Post(a.Config.UpdateMetricsAddress)

	if err != nil {
		log.Printf("error sending request: %v\n", err)
		return fmt.Errorf("error sending request: %w", err)
	}

	if res.StatusCode() != http.StatusOK {
		log.Printf("error sending request: status: %v, response: %v\n", res.StatusCode(), res)
		return fmt.Errorf("error sending request: response: %v", res)
	}

	return nil
}
