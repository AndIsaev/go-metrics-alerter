package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/AndIsaev/go-metrics-alerter/internal/common"

	"github.com/AndIsaev/go-metrics-alerter/internal/service/agent/metrics"
)

// Config use for setting agent application
type Config struct {
	// Address is host for application
	Address string `env:"ADDRESS" json:"address"`
	// ReportInterval set interval sending metrics to server
	ReportInterval time.Duration `env:"REPORT_INTERVAL" json:"report_interval"`
	// PollInterval how often to update metrics
	PollInterval time.Duration `env:"POLL_INTERVAL" json:"poll_interval"`
	// Key secret key for connect to server
	Key string `env:"KEY" json:"key"`
	// RateLimit limit connections to server
	RateLimit int `env:"RATE_LIMIT" json:"rate_limit"`
	// StorageMetrics storage with metrics
	StorageMetrics *metrics.StorageMetrics
	// UpdateMetricAddress address for update one metric
	UpdateMetricAddress string
	// UpdateMetricsAddress address for update batch with metrics
	UpdateMetricsAddress string
	// PublicKey public key
	PublicKey *rsa.PublicKey
	// PublicKeyPath path of public key
	PublicKeyPath string `env:"CRYPTO_KEY" json:"crypto_key"`
	// RPCClient define bool variable for use rpc client
	RPCClient bool `env:"RPC_CLIENT" json:"rpc_client"`
	// ConfigPath path of config file
	ConfigPath string `env:"CONFIG"`
}

// NewConfig create new config
func NewConfig() *Config {
	cfg := &Config{StorageMetrics: metrics.NewListMetrics()}
	var pollIntervalSeconds uint64
	var reportIntervalSeconds uint64
	var rateLimit uint64

	flag.StringVar(&cfg.Address, "a", ":8080", "address")
	flag.Uint64Var(&reportIntervalSeconds, "r", 10, "seconds of report interval")
	flag.Uint64Var(&pollIntervalSeconds, "p", 2, "seconds of poll interval")
	flag.StringVar(&cfg.Key, "k", "", "set key")
	flag.Uint64Var(&rateLimit, "l", 10, "rate limit")
	flag.StringVar(&cfg.PublicKeyPath, "crypto-key", "", "set path of public key")
	flag.BoolVar(&cfg.RPCClient, "rpc", false, "set true if yor want use rrc")
	// config path
	configFile := flag.String("c", "", "Path to the configuration file")
	flag.StringVar(configFile, "config", "", "Path to the configuration file (alias for -c)")

	flag.Parse()
	cfg.ConfigPath = *configFile

	if envRPCClient := os.Getenv("RPC_CLIENT"); envRPCClient != "" {
		parseBool, err := strconv.ParseBool(envRPCClient)
		if err == nil {
			cfg.RPCClient = parseBool
		} else {
			log.Println("error parse RPC_CLIENT variable, must be bool value")
		}
	}

	if envPrivateKeyPath := os.Getenv("CRYPTO_KEY"); envPrivateKeyPath != "" {
		cfg.PublicKeyPath = envPrivateKeyPath
	}

	if envRateLimit := os.Getenv("RATE_LIMIT"); envRateLimit != "" {
		if val, err := strconv.Atoi(envRateLimit); err == nil {
			cfg.RateLimit = val
		}
	} else {
		cfg.RateLimit = int(rateLimit)
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
	cfg.UpdateMetricAddress = fmt.Sprintf("http://localhost%s/update/", cfg.Address)
	cfg.UpdateMetricsAddress = fmt.Sprintf("http://localhost%s/updates/", cfg.Address)

	if cfg.PublicKeyPath != "" {
		publicKey, err := cfg.getPublicKey()

		if err != nil {
			log.Printf("error with get public key: %v\n", err)
		}
		cfg.PublicKey = publicKey
	}
	if envConfig := os.Getenv("CONFIG"); envConfig != "" {
		cfg.ConfigPath = envConfig
	}
	// загружаем конфиг
	if cfg.ConfigPath != "" {
		keys := []string{"report_interval", "poll_interval"}
		body, err := common.LoadConfigFromJSON(cfg.ConfigPath, keys)
		if err != nil {
			log.Printf("failed to load configuration from file: %v\n", err)
		}
		err = json.Unmarshal(body, cfg)
		if err != nil {
			log.Printf("could not unmarshal config file: %s\n", err)
		}
	}

	return cfg
}

func (c *Config) getPublicKey() (*rsa.PublicKey, error) {
	publicKeyPEM, err := os.ReadFile(c.PublicKeyPath)
	if err != nil {
		log.Printf("error reading public key file: %v", err)
		return nil, err
	}

	publicKeyDecode, _ := pem.Decode(publicKeyPEM)
	if publicKeyDecode == nil {
		err = errors.New("failed to decode PEM block containing public key")
		log.Println(err)
		return nil, err
	}

	parsedPublicKey, err := x509.ParsePKIXPublicKey(publicKeyDecode.Bytes)
	if err != nil {
		log.Printf("error parsing public key: %v", err)
		return nil, err
	}

	rsaPubKey, ok := parsedPublicKey.(*rsa.PublicKey)
	if !ok {
		err = errors.New("not an RSA public key")
		log.Println(err)
		return nil, err
	}

	return rsaPubKey, nil
}
