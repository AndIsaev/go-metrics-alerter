package server

import (
	"flag"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Address         string        `env:"ADDRESS"`
	StoreInterval   time.Duration `env:"STORE_INTERVAL"`
	FileStoragePath string        `env:"FILE_STORAGE_PATH"`
	Restore         bool          `env:"RESTORE"`
	DBDsn           string        `env:"DATABASE_DSN"`
	Key             string        `env:"KEY"`
}

func NewConfig() *Config {
	cfg := &Config{}
	var storeInterval uint64
	var fileStoragePath string
	var dbDsn string

	flag.StringVar(&cfg.Address, "a", ":8080", "server address")
	flag.BoolVar(&cfg.Restore, "r", true, "load metrics from file")
	flag.StringVar(&fileStoragePath, "f", "./metrics", "path of metrics on disk")
	flag.Uint64Var(&storeInterval, "i", 300, "interval for save metrics on file")
	flag.StringVar(&dbDsn, "d", "", "database dsn")
	flag.StringVar(&cfg.Key, "k", "", "set key")

	flag.Parse()

	if envKey := os.Getenv("KEY"); envKey != "" {
		cfg.Key = envKey
	}

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
