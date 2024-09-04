package service

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/AndIsaev/go-metrics-alerter/internal/service/agent/metrics"
)

type ServerConfig struct {
	Address         string        `env:"ADDRESS"`
	StoreInterval   time.Duration `env:"STORE_INTERVAL"`
	FileStoragePath string        `env:"FILE_STORAGE_PATH"`
	Restore         bool          `env:"RESTORE"`
	DBDsn           string        `env:"DATABASE_DSN"`
}

func NewServerConfig() *ServerConfig {
	cfg := &ServerConfig{}
	var storeInterval uint64
	var fileStoragePath string
	var dbDsn string

	flag.StringVar(&cfg.Address, "a", "localhost:8080", "server address")
	flag.BoolVar(&cfg.Restore, "r", true, "load metrics from file")
	flag.StringVar(&fileStoragePath, "f", "./metrics", "path of metrics on disk")
	flag.Uint64Var(&storeInterval, "i", 300, "interval for save metrics on file")
	flag.StringVar(&dbDsn, "d", "", "database dsn")

	flag.Parse()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		cfg.Address = envRunAddr
	}

	if envDBDsn := os.Getenv("DATABASE_DSN"); envDBDsn != "" {
		cfg.DBDsn = envDBDsn
	} else {
		cfg.DBDsn = dbDsn
	}

	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		cfg.FileStoragePath = envFileStoragePath
	} else {
		cfg.FileStoragePath = fileStoragePath
	}

	if envStoreInterval := os.Getenv("STORE_INTERVAL"); envStoreInterval != "" {
		if val, err := strconv.Atoi(envStoreInterval); err == nil {
			cfg.StoreInterval = time.Duration(val) * time.Second
		}
	} else {
		cfg.StoreInterval = time.Duration(storeInterval) * time.Second
	}

	return cfg
}

type AgentConfig struct {
	Address             string        `env:"ADDRESS"`
	ReportInterval      time.Duration `env:"REPORT_INTERVAL"`
	PollInterval        time.Duration `env:"POLL_INTERVAL"`
	StorageMetrics      *metrics.StorageMetrics
	UpdateMetricAddress string
	ProtocolHTTP        string
}

func NewAgentConfig() *AgentConfig {
	cfg := &AgentConfig{StorageMetrics: metrics.NewListMetrics(), ProtocolHTTP: "http"}
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
	cfg.UpdateMetricAddress = fmt.Sprintf("%s://%s/update/", cfg.ProtocolHTTP, cfg.Address)

	return cfg
}
