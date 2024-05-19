package service

import (
	"flag"
	"os"
	"strconv"
	"time"
)

type ServerConfig struct {
	Address string `env:"ADDRESS"`
}

func NewServerConfig() *ServerConfig {
	cfg := &ServerConfig{}

	flag.StringVar(&cfg.Address, "a", "0.0.0.0:8080", "server address")

	flag.Parse()
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		cfg.Address = envRunAddr
	}

	return cfg
}

type AgentConfig struct {
	Address        string        `env:"ADDRESS"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
}

func NewAgentConfig() *AgentConfig {
	cfg := &AgentConfig{}
	var pollIntervalSeconds uint64
	var reportIntervalSeconds uint64

	flag.StringVar(&cfg.Address, "a", "0.0.0.0:8080", "address")
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

	return cfg
}
