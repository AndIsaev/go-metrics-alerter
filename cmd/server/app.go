package main

import (
	"context"
	"fmt"
	"github.com/AndIsaev/go-metrics-alerter/internal/manager/file"
	"github.com/AndIsaev/go-metrics-alerter/internal/service"
	"github.com/AndIsaev/go-metrics-alerter/internal/service/server/handlers"
	"github.com/AndIsaev/go-metrics-alerter/internal/storage"
	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v5"
	"log"
	"os"
)

// ServerApp - structure of application
type ServerApp struct {
	Route        chi.Router
	MemStorage   *storage.MemStorage
	FileProducer *file.Producer
	FileConsumer *file.Consumer
	DbConn       storage.PgStorage
}

// NewServerApp - create new app
func NewServerApp(cfg *service.ServerConfig) *ServerApp {
	app := &ServerApp{}

	// create directory
	if err := createDir(cfg.FileStoragePath); err != nil {
		log.Fatal(err)
	}

	// init storage
	app.MemStorage = storage.NewMemStorage()

	// set producer and consumer for file manager
	producer, _ := file.NewProducer(cfg.FileStoragePath)
	consumer, _ := file.NewConsumer(cfg.FileStoragePath)

	app.FileProducer = producer
	app.FileConsumer = consumer

	// download metrics from disc to storage
	if cfg.Restore {
		downloadMetrics(app.FileConsumer, app.MemStorage)
		defer app.FileConsumer.Close()
	}

	// connect to DB
	conn, err := pgx.Connect(context.Background(), cfg.DbDsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	app.DbConn = conn

	// init route
	app.Route = handlers.ServerRouter(app.MemStorage, app.FileProducer, app.DbConn)

	return app
}

// downloadMetrics - Read metrics from disk
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

// createDir - create directory for metrics
func createDir(fileStoragePath string) error {
	if _, err := os.Stat(fileStoragePath); os.IsNotExist(err) {
		if err = os.Mkdir(fileStoragePath, 0755); err != nil {
			fmt.Printf("the directory %s not created\n", fileStoragePath)
			return err
		}
	}
	fmt.Printf("the directory %s is done\n", fileStoragePath)
	return nil
}
