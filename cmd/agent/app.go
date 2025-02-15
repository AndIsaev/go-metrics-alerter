package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/AndIsaev/go-metrics-alerter/internal/common"
	"github.com/AndIsaev/go-metrics-alerter/internal/service/agent/client"
	"github.com/AndIsaev/go-metrics-alerter/internal/service/agent/client/http"
	"github.com/AndIsaev/go-metrics-alerter/internal/service/agent/client/rpc"
	"github.com/AndIsaev/go-metrics-alerter/internal/service/agent/utils"
)

// AgentApp structure of application
type AgentApp struct {
	// Config use for settings of app
	Config *Config
	// Client use for requests to server
	Client client.RequestClient
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
	a.jobs = make(chan common.Metrics)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// init Client for requests
	a.initRequestClient(ctx, cancel)

	// обработка сигналов завершения
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	select {
	case <-ctx.Done():
		log.Println("context done -> exit from StartApp")
		return
	default:
		// горутина для прослушивания сигналов
		a.wg.Add(1)
		go func() {
			defer a.wg.Done()
			sig := <-sigs
			log.Printf("Received signal: %s", sig)
			cancel() // отменяем контекст
		}()
	}

	a.runReport(ctx)

	a.wg.Wait()

	close(a.jobs)

	log.Println("application stopped gracefully")
}

// initRequestClient init request client interface
func (a *AgentApp) initRequestClient(ctx context.Context, cancel context.CancelFunc) {
	select {
	case <-ctx.Done():
		log.Println("context done -> exit from initRequestClient")
		return
	default:

		if !a.Config.RPCClient {
			a.initHTTPClient()
		} else {
			a.initGRPCClient(cancel)
		}
	}
}

// initHTTPClient - init http client
func (a *AgentApp) initHTTPClient() {
	ipResolver := utils.NewDefaultIPResolver(a.Config.Address)
	httpClient := http.NewClient(
		a.Config.UpdateMetricsAddress,
		ipResolver,
		a.Config.Key,
		a.Config.PublicKey,
	)
	a.Client = httpClient
}

// initGRPCClient - init grpc client
func (a *AgentApp) initGRPCClient(cancel context.CancelFunc) {
	conn, err := grpc.NewClient(a.Config.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("error init grpc client, cancel context: %v\n", err)
		conn.Close()
		cancel()
		return
	}
	grpcClient := rpc.NewGRPCClient(conn)
	rpcClient := rpc.NewClient(grpcClient)
	a.Client = rpcClient
}

func (a *AgentApp) runReport(ctx context.Context) {
	a.wg.Add(2)
	go a.pullMetrics(ctx)

	go a.runWorkers(ctx)
}

// pullMetrics - get metrics from runtime and save to StorageMetrics
func (a *AgentApp) pullMetrics(ctx context.Context) {
	defer a.wg.Done()
	for {
		select {
		case <-ctx.Done():
			log.Println("context done -> exit from pullMetrics")
			return
		default:
			log.Println("pull metrics")

			a.Config.StorageMetrics.Pull()

			a.mu.RLock()
			// Передача метрик в канал
			for _, val := range a.Config.StorageMetrics.Metrics {
				select {
				case a.jobs <- val:
				case <-ctx.Done():
					log.Println("context done -> exit from pullMetrics")
					return
				}
			}
			a.mu.RUnlock()
			log.Println("the metrics have been sent to the channel")

			time.Sleep(a.Config.PollInterval)
		}
	}
}

func (a *AgentApp) runWorkers(ctx context.Context) {
	defer a.wg.Done()
	var subWg sync.WaitGroup
	subWg.Add(a.Config.RateLimit)

	for w := 1; w <= a.Config.RateLimit; w++ {
		go func() {
			defer subWg.Done()
			for {
				select {
				case <-ctx.Done(): // Завершение по сигналу отмены
					log.Println("context done -> exit from runWorkers")
					return
				default:
					tasks := make([]common.Metrics, 0, 3)
					for i := 0; i < 3; i++ {
						metric, ok := <-a.jobs
						if !ok {
							log.Printf("jobs channel closed")
							break
						}
						tasks = append(tasks, metric)
					}
					if len(tasks) > 0 {
						if err := utils.Retry(a.Client.SendMetrics)(ctx, tasks); err != nil {
							log.Printf("Error sending metrics: %v", err)
						}
					}

					interval := a.Config.ReportInterval
					time.Sleep(interval)
				}
			}
		}()
	}
	subWg.Wait()
	log.Println("all workers done")
}
