package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"flag"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/AndIsaev/go-metrics-alerter/internal/common"
)

// Config use for setting server application
type Config struct {
	// Address is host for application
	Address string `env:"ADDRESS" json:"address"`
	// StoreInterval interval for save metrics on file
	StoreInterval time.Duration `env:"STORE_INTERVAL" json:"store_interval"`
	// FileStoragePath path of metrics on disk
	FileStoragePath string `env:"FILE_STORAGE_PATH" json:"file_store_file"`
	// Restore load metrics from file when start app
	Restore bool `env:"RESTORE" json:"restore"`
	// Dsn database dsn
	Dsn string `env:"DATABASE_DSN" json:"database_dsn"`
	// Key for access to metrics
	Key string `env:"KEY" json:"key"`
	// PrivateKey path of private key of server
	PrivateKey string `env:"CRYPTO_KEY" json:"crypto_key"`
	ConfigPath string `env:"CONFIG"`
}

// NewConfig create new config
func NewConfig() *Config {
	cfg := &Config{}
	var storeInterval uint64
	var fileStoragePath string
	var dbDsn string
	var restore bool

	flag.StringVar(&cfg.Address, "a", ":8080", "server address")
	flag.BoolVar(&restore, "r", true, "load metrics from file")
	flag.StringVar(&fileStoragePath, "f", "./metrics", "path of metrics on disk")
	flag.Uint64Var(&storeInterval, "i", 300, "interval for save metrics on file")
	flag.StringVar(&dbDsn, "d", "", "database dsn")
	flag.StringVar(&cfg.Key, "k", "", "set key")
	flag.StringVar(&cfg.PrivateKey, "crypto-key", "", "set path of private key")
	// config path
	configFile := flag.String("c", "", "Path to the configuration file")
	flag.StringVar(configFile, "config", "", "Path to the configuration file (alias for -c)")

	flag.Parse()
	cfg.ConfigPath = *configFile

	if envPrivateKey := os.Getenv("CRYPTO_KEY"); envPrivateKey != "" {
		cfg.PrivateKey = envPrivateKey
	}

	if envKey := os.Getenv("KEY"); envKey != "" {
		cfg.Key = envKey
	}

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		cfg.Address = envRunAddr
	}

	if envRestore := os.Getenv("RESTORE"); envRestore != "" {
		val, err := strconv.ParseBool(envRestore)
		if err != nil {
			log.Println("error parse r flag, must be boolean value")
		}
		cfg.Restore = val
	} else {
		cfg.Restore = restore
	}

	if envDBDsn := os.Getenv("DATABASE_DSN"); envDBDsn != "" {
		cfg.Dsn = envDBDsn
	} else {
		cfg.Dsn = dbDsn
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
	if envConfig := os.Getenv("CONFIG"); envConfig != "" {
		cfg.ConfigPath = envConfig
	}
	// загружаем конфиг
	if cfg.ConfigPath != "" {
		keys := []string{"store_interval"}
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

func (c *Config) GetPrivateKey() (*rsa.PrivateKey, error) {
	privateKey, err := os.ReadFile(c.PrivateKey)
	if err != nil {
		log.Println("error read file of private key")
		return nil, errors.Unwrap(err)
	}

	privateKeyDecode, _ := pem.Decode(privateKey)
	parsedPrivateKey, err := x509.ParsePKCS1PrivateKey(privateKeyDecode.Bytes)
	if err != nil {
		log.Println("error parse private key")
		return nil, errors.Unwrap(err)
	}

	return parsedPrivateKey, nil
}
