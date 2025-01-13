package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/AndIsaev/go-metrics-alerter/internal/service/agent/metrics"
)

// Config use for setting agent application
type Config struct {
	// Address is host for application
	Address string `env:"ADDRESS"`
	// ReportInterval set interval sending metrics to server
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
	// PollInterval how often to update metrics
	PollInterval time.Duration `env:"POLL_INTERVAL"`
	// Key secret key for connect to server
	Key string `env:"KEY"`
	// RateLimit limit connections to server
	RateLimit uint64 `env:"RATE_LIMIT"`
	// StorageMetrics storage with metrics
	StorageMetrics *metrics.StorageMetrics
	// UpdateMetricAddress address for update one metric
	UpdateMetricAddress string
	// UpdateMetricsAddress address for update batch with metrics
	UpdateMetricsAddress string
	ProtocolHTTP         string
}

// NewConfig create new config
func NewConfig() *Config {
	cfg := &Config{StorageMetrics: metrics.NewListMetrics(), ProtocolHTTP: "http"}
	var pollIntervalSeconds uint64
	var reportIntervalSeconds uint64
	var rateLimit uint64

	flag.StringVar(&cfg.Address, "a", ":8080", "address")
	flag.Uint64Var(&reportIntervalSeconds, "r", 10, "seconds of report interval")
	flag.Uint64Var(&pollIntervalSeconds, "p", 2, "seconds of poll interval")
	flag.StringVar(&cfg.Key, "k", "", "set key")
	flag.Uint64Var(&rateLimit, "l", 10, "rate limit")

	flag.Parse()

	if envRateLimit := os.Getenv("RATE_LIMIT"); envRateLimit != "" {
		if val, err := strconv.Atoi(envRateLimit); err == nil {
			cfg.RateLimit = uint64(val)
		}
	} else {
		cfg.RateLimit = rateLimit
	}

	if envKey := os.Getenv("KEY"); envKey != "" {
		cfg.Key = envKey
	}

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
	cfg.UpdateMetricsAddress = fmt.Sprintf("%s://%s/updates/", cfg.ProtocolHTTP, cfg.Address)

	return cfg
}
