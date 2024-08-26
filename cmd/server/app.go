package main

import (
	"context"
	"fmt"

	"github.com/go-chi/chi/middleware"
	"github.com/mailru/easyjson"

	"github.com/AndIsaev/go-metrics-alerter/internal/common"
	"github.com/AndIsaev/go-metrics-alerter/internal/service/server/handlers"
	mid "github.com/AndIsaev/go-metrics-alerter/internal/service/server/middleware"

	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"

	"github.com/AndIsaev/go-metrics-alerter/internal/logger"

	"github.com/go-chi/chi"

	"github.com/AndIsaev/go-metrics-alerter/internal/manager/file"
	"github.com/AndIsaev/go-metrics-alerter/internal/service"
	"github.com/AndIsaev/go-metrics-alerter/internal/storage"
)

// ServerApp - structure of application
type ServerApp struct {
	Router       chi.Router
	MemStorage   *storage.MemStorage
	FileProducer *file.Producer
	FileConsumer *file.Consumer
	DbConn       storage.PgStorage
	Config       *service.ServerConfig
	Server       *http.Server
}

// New - create new app
func New() *ServerApp {
	app := &ServerApp{}
	config := service.NewServerConfig()
	app.Config = config

	// init file storage
	app.MemStorage = storage.NewMemStorage()

	// set producer and consumer for file manager
	producer, _ := file.NewProducer(config.FileStoragePath)
	consumer, _ := file.NewConsumer(config.FileStoragePath)

	app.FileProducer = producer
	app.FileConsumer = consumer

	app.Router = chi.NewRouter()

	return app
}

func (a *ServerApp) StartApp(ctx context.Context) error {
	if err := logger.Initialize(); err != nil {
		return err
	}

	// create directory
	if err := a.createMetricsDir(); err != nil {
		log.Fatalf("can't create directory: %s\n", err.Error())
	}

	// download metrics from disc to storage
	a.downloadMetrics()

	// connect to DB
	a.connectDB(ctx)

	// init router
	a.initRouter()

	// init http server
	a.initHTTPServer()

	return a.startHTTPServer()
}

// startHTTPServer - start http server
func (a *ServerApp) startHTTPServer() error {
	fmt.Printf("start server on: %s\n", a.Config.Address)
	return a.Server.ListenAndServe()
}

// initHTTPServer - init http server
func (a *ServerApp) initHTTPServer() {
	server := &http.Server{}
	server.Handler = a.Router
	server.Addr = a.Config.Address
	a.Server = server
}

// downloadMetrics - Read metrics from disk
func (a *ServerApp) downloadMetrics() {
	if a.Config.Restore {
		fmt.Println("read metrics from disk")
		for {
			m, err := a.FileConsumer.ReadMetrics()
			if err != nil {
				break
			}
			a.MemStorage.Set(m)
		}
		fmt.Println("metrics downloaded")
	}
	defer a.FileConsumer.Close()
}

// createMetricsDir - create directory for metrics
func (a *ServerApp) createMetricsDir() error {
	if _, err := os.Stat(a.Config.FileStoragePath); os.IsNotExist(err) {
		if err = os.Mkdir(a.Config.FileStoragePath, 0755); err != nil {
			fmt.Printf("the directory %s not created\n", a.Config.FileStoragePath)
			return err
		}
	}
	fmt.Printf("the directory %s is done\n", a.Config.FileStoragePath)
	return nil
}

// connectDB - connection to database
func (a *ServerApp) connectDB(ctx context.Context) {
	if a.Config.DbDsn != "" {
		conn, err := pgx.Connect(ctx, a.Config.DbDsn)
		if err != nil {
			log.Fatalf("unable to connect to database: %s\n", err.Error())
		}

		a.DbConn = conn
	}
}

func (a *ServerApp) Shutdown(ctx context.Context) {
	if err := a.FileProducer.Close(); err != nil {
		fmt.Printf("%s\n", err.Error())
	}
	if err := a.FileConsumer.Close(); err != nil {
		fmt.Printf("%s\n", err.Error())
	}
	if err := a.DbConn.Close(ctx); err != nil {
		fmt.Printf("%s\n", err.Error())
	}
}

// initRouter - initialize new router
func (a *ServerApp) initRouter() {
	r := a.Router
	r.Use(logger.RequestLogger, logger.ResponseLogger)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.Recoverer)
	r.Use(middleware.SetHeader("Content-Type", "application/json"))

	r.NotFound(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		body := common.Response{Message: "route does not exist"}
		response, _ := easyjson.Marshal(body)
		w.Write(response)
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		body := common.Response{Message: "method is not valid"}
		response, _ := easyjson.Marshal(body)
		w.Write(response)
	})

	// Routes

	r.Group(func(r chi.Router) {
		r.Use(mid.GzipMiddleware)

		// Ping db connection
		r.Get(`/ping`, handlers.PingHandler(a.DbConn))

		// update
		r.Post(`/update/{MetricType}/{MetricName}/{MetricValue}`, handlers.SetMetricHandler(a.MemStorage))
		r.Post(`/update`, handlers.UpdateHandler(a.MemStorage, a.FileProducer))

		// value
		r.Get(`/value/{MetricType}/{MetricName}`, handlers.GetMetricHandler(a.MemStorage))
		r.Post(`/value`, handlers.GetHandler(a.MemStorage))

		// main page
		r.Get(`/`, handlers.MainPageHandler(a.MemStorage))
	})
}
