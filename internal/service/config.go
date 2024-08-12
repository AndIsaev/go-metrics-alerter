package service

import (
	"flag"
	"fmt"
	"github.com/AndIsaev/go-metrics-alerter/internal/manager/file"
	"github.com/AndIsaev/go-metrics-alerter/internal/service/agent/metrics"
	"github.com/AndIsaev/go-metrics-alerter/internal/service/server/handlers"
	"github.com/AndIsaev/go-metrics-alerter/internal/storage"
	"github.com/go-chi/chi"
	"os"
	"strconv"
	"time"
)

type ServerConfig struct {
	Address      string `env:"ADDRESS"`
	Route        chi.Router
	MemStorage   *storage.MemStorage
	FileProducer *file.Producer
	FileConsumer *file.Consumer

	StoreInterval   time.Duration `env:"STORE_INTERVAL"`
	FileStoragePath string        `env:"FILE_STORAGE_PATH"`
	Restore         bool          `env:"RESTORE"`
}

func NewServerConfig() *ServerConfig {
	cfg := &ServerConfig{}
	var storeInterval uint64
	var fileStoragePath string

	cfg.MemStorage = storage.NewMemStorage()

	flag.StringVar(&cfg.Address, "a", "localhost:8080", "server address")
	flag.BoolVar(&cfg.Restore, "r", true, "Load metrics from file")
	flag.StringVar(&fileStoragePath, "f", "./metrics", "path of metrics on disk")
	flag.Uint64Var(&storeInterval, "i", 300, "interval for save metrics on file")

	flag.Parse()
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		cfg.Address = envRunAddr
	}

	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		createDir(envFileStoragePath)
		cfg.FileStoragePath = envFileStoragePath
	} else {
		createDir(fileStoragePath)
		cfg.FileStoragePath = fileStoragePath
	}

	if envStoreInterval := os.Getenv("STORE_INTERVAL"); envStoreInterval != "" {
		if val, err := strconv.Atoi(envStoreInterval); err == nil {
			cfg.StoreInterval = time.Duration(val) * time.Second
		}
	} else {
		cfg.StoreInterval = time.Duration(storeInterval) * time.Second
	}

	// set producer and consumer for file manager
	producer, _ := file.NewProducer(cfg.FileStoragePath)
	consumer, _ := file.NewConsumer(cfg.FileStoragePath)

	cfg.FileProducer = producer
	cfg.FileConsumer = consumer

	// set main server router
	cfg.Route = handlers.ServerRouter(cfg.MemStorage, cfg.FileProducer)

	// download metrics from disc to storage
	if cfg.Restore {
		downloadMetrics(cfg.FileConsumer, cfg.MemStorage)
		defer cfg.FileConsumer.Close()
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

func createDir(fileStoragePath string) {
	if _, err := os.Stat(fileStoragePath); os.IsNotExist(err) {
		os.Mkdir(fileStoragePath, 0755)
		dir := fmt.Sprintf("Directory %s created", fileStoragePath)
		fmt.Println(dir)
	} else {
		dir := fmt.Sprintf("Directory %s already exists", fileStoragePath)
		fmt.Println(dir)
	}
}

func downloadMetrics(consumer *file.Consumer, memStorage *storage.MemStorage) {

	fmt.Println("Read metrics from disk")
	for {
		m, err := consumer.ReadMetrics()
		if err != nil {
			break
		}
		memStorage.Set(m)
	}
	fmt.Println("Metrics downloaded")
}
