package service

import (
	"flag"
	"fmt"
	"github.com/AndIsaev/go-metrics-alerter/internal/service/agent/metrics"
	"github.com/AndIsaev/go-metrics-alerter/internal/service/server/handlers"
	"github.com/AndIsaev/go-metrics-alerter/internal/storage"
	"github.com/go-chi/chi"
	"os"
	"strconv"
	"time"
)

type ServerConfig struct {
	Address    string `env:"ADDRESS"`
	Route      chi.Router
	MemStorage *storage.MemStorage
}

func NewServerConfig() *ServerConfig {
	cfg := &ServerConfig{}
	cfg.MemStorage = storage.NewMemStorage()
	cfg.Route = handlers.ServerRouter(cfg.MemStorage)

	flag.StringVar(&cfg.Address, "a", "localhost:8080", "server address")

	flag.Parse()
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		cfg.Address = envRunAddr
	}

	return cfg
}

type AgentConfig struct {
	Address             string        `env:"ADDRESS"`
	ReportInterval      time.Duration `env:"REPORT_INTERVAL"`
	PollInterval        time.Duration `env:"POLL_INTERVAL"`
	StorageMetrics      *metrics.StorageMetrics
	UpdateMetricAddress string
	ProtocolHttp        string
}

func NewAgentConfig() *AgentConfig {
	cfg := &AgentConfig{StorageMetrics: metrics.NewListMetrics(), ProtocolHttp: "http"}
	var pollIntervalSeconds uint64
	var reportIntervalSeconds uint64

	flag.StringVar(&cfg.Address, "a", "localhost:8080", "address")
	flag.Uint64Var(&reportIntervalSeconds, "r", 10, "seconds of report interval")
	flag.Uint64Var(&pollIntervalSeconds, "p", 2, "seconds of poll interval")

	flag.Parse()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		cfg.Address = envRunAddr
	}

	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		if val, err := strconv.Atoi(envReportInterval); err == nil {
			cfg.ReportInterval = time.Duration(val) * time.Second
		}
	} else {
		cfg.ReportInterval = time.Duration(reportIntervalSeconds) * time.Second
	}

	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		if val, err := strconv.Atoi(envPollInterval); err == nil {
			cfg.PollInterval = time.Duration(val) * time.Second
		}
	} else {
		cfg.PollInterval = time.Duration(pollIntervalSeconds) * time.Second
	}
	// set address for update metric
	cfg.UpdateMetricAddress = fmt.Sprintf("%s://%s/update/", cfg.ProtocolHttp, cfg.Address)

	return cfg
}
