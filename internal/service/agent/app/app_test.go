package app

import (
	"context"
	"testing"
	"time"

	"github.com/AndIsaev/go-metrics-alerter/internal/service/agent/config"

	"github.com/AndIsaev/go-metrics-alerter/internal/common"

	"github.com/stretchr/testify/assert"

	"github.com/AndIsaev/go-metrics-alerter/internal/service/agent/client"
	"github.com/AndIsaev/go-metrics-alerter/internal/service/agent/client/http"
	"github.com/AndIsaev/go-metrics-alerter/internal/service/agent/client/rpc"
)

func TestNew(t *testing.T) {
	app := New()
	assert.NotNil(t, app.Config)
	assert.IsType(t, &config.Config{}, app.Config)
}

func TestInitRequestClient(t *testing.T) {
	app := New()
	app.Config = config.NewConfig()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tests := []struct {
		name   string
		http   bool
		client client.RequestClient
	}{
		{name: "http client", http: true, client: &http.Client{}},
		{name: "grpc client", http: false, client: &rpc.Client{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.http {
				app.Config.RPCClient = true
			}
			app.initRequestClient(ctx, cancel)
			assert.IsType(t, tt.client, app.Client)
		})
	}
}

func TestInitHTTPClient(t *testing.T) {
	app := New()
	app.Config = config.NewConfig()

	app.initHTTPClient()

	assert.IsType(t, &http.Client{}, app.Client)
}

func TestInitGRPCClient(t *testing.T) {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()
	app := New()
	app.Config = config.NewConfig()

	app.initGRPCClient(cancel)

	assert.IsType(t, &rpc.Client{}, app.Client)
}

type MockStorageMetrics struct {
	Metrics map[string]common.Metrics
}

func (m *MockStorageMetrics) Pull() {
}

func (m *MockStorageMetrics) List() map[string]common.Metrics {
	return m.Metrics
}

func NewListMockStorageMetrics() *MockStorageMetrics {
	mockMetrics := make(map[string]common.Metrics)
	mockMetrics["metric1"] = common.Metrics{ID: "metric1", MType: common.Counter, Delta: common.LinkInt64(1)}
	return &MockStorageMetrics{Metrics: mockMetrics}
}

func TestPullMetrics(t *testing.T) {
	mapMetrics := NewListMockStorageMetrics()

	mockStorage := &MockStorageMetrics{Metrics: mapMetrics.Metrics}
	config := &config.Config{
		StorageMetrics: mockStorage,
		PollInterval:   100 * time.Millisecond,
	}

	app := &AgentApp{
		Config: config,
		jobs:   make(chan common.Metrics, 10),
	}

	ctx, cancel := context.WithCancel(context.Background())
	app.wg.Add(1)
	go app.pullMetrics(ctx)

	time.Sleep(config.PollInterval * 2)

	for _, expMetric := range mapMetrics.Metrics {
		select {
		case metric := <-app.jobs:
			assert.Equal(t, expMetric, metric, "Expected metric not received from channel.")
		default:
			t.Errorf("expected metric in channel, but got none")
		}
	}

	cancel()
	app.wg.Wait()
}

// Псевдо-реализация клиентского интерфейса
type MockClient struct {
	sendCalled int
	sendErr    error
}

func (c *MockClient) SendMetrics(_ context.Context, _ []common.Metrics) error {
	if c.sendErr != nil {
		return c.sendErr
	}
	c.sendCalled++
	return nil
}

func TestRunWorkers(t *testing.T) {
	mockClient := &MockClient{}
	config := &config.Config{
		RateLimit:      3,
		ReportInterval: 50 * time.Millisecond,
	}

	app := &AgentApp{
		Config: config,
		Client: mockClient,
		jobs:   make(chan common.Metrics, 3),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	app.wg.Add(1)
	go app.runWorkers(ctx)

	metrics := []common.Metrics{
		{ID: "metric1", MType: common.Counter, Delta: common.LinkInt64(1)},
		{ID: "metric2", MType: common.Counter, Delta: common.LinkInt64(5)},
		{ID: "metric3", MType: common.Gauge, Value: common.LinkFloat64(12.1)},
	}

	for _, metric := range metrics {
		app.jobs <- metric
	}

	close(app.jobs)

	time.Sleep(300 * time.Millisecond)

	assert.Greater(t, mockClient.sendCalled, 0, "Expected SendMetrics to be called")

	app.wg.Wait()
}
