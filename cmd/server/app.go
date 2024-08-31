package main

import (
	"context"
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
	DBConn       storage.PgStorage
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

	app.Router = chi.NewRouter()

	return app
}

func (a *ServerApp) StartApp(ctx context.Context) error {
	if err := logger.Initialize(); err != nil {
		return err
	}

	if a.Config.DBDsn != "" {
		// connect to DB
		if err := a.connectToDB(ctx); err != nil {
			return err
		}

		// создаем таблицы
		if err := a.createTables(ctx); err != nil {
			return err
		}
	}

	if a.Config.FileStoragePath != "" {
		// create directory
		if err := createMetricsDir(a.Config.FileStoragePath); err != nil {
			log.Printf("can't create directory because of: %s\n", err.Error())
			return err
		}
		// set producer and consumer for file manager
		if err := a.initFileManagers(); err != nil {
			return err
		}

		// download metrics from disc to storage
		a.downloadMetrics()
	}

	// init router
	a.initRouter()

	// init http server
	a.initHTTPServer()

	return a.startHTTPServer()
}

// startHTTPServer - start http server
func (a *ServerApp) startHTTPServer() error {
	log.Printf("start server on: %s\n", a.Config.Address)
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
		log.Println("read metrics from disk")
		for {
			m, err := a.FileConsumer.ReadMetrics()
			if err != nil {
				break
			}
			a.MemStorage.Set(m)
		}
		log.Println("metrics downloaded")
	}
}

// createMetricsDir - create directory for metrics
func createMetricsDir(fileStoragePath string) error {
	if _, err := os.Stat(fileStoragePath); os.IsNotExist(err) {
		if err = os.Mkdir(fileStoragePath, 0755); err != nil {
			log.Printf("the directory %s not created\n", fileStoragePath)
			return err
		}
	}
	log.Printf("the directory %s is done\n", fileStoragePath)
	return nil
}

// connectToDB - connection to database
func (a *ServerApp) connectToDB(ctx context.Context) error {
	if a.Config.DBDsn != "" {
		conn, err := pgx.Connect(ctx, a.Config.DBDsn)
		if err != nil {
			log.Printf("unable to connect to database: %s\n", err.Error())
			return err
		}
		a.DBConn = conn
	}
	return nil
}

func (a *ServerApp) Shutdown(ctx context.Context) {
	if a.Config.FileStoragePath != "" {
		if err := a.FileProducer.Close(); err != nil {
			log.Printf("%s\n", err.Error())
		}
		if err := a.FileConsumer.Close(); err != nil {
			log.Printf("%s\n", err.Error())
		}
	}

	if a.Config.DBDsn != "" {
		if a.DBConn != nil {
			if err := a.DBConn.Close(ctx); err != nil {
				log.Printf("%s\n", err.Error())
			}
		}
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
		r.Get(`/ping`, handlers.PingHandler(a.DBConn))

		// update
		r.Post(`/update/{MetricType}/{MetricName}/{MetricValue}`, handlers.SetMetricHandler(a.MemStorage))
		r.Post(`/update`, handlers.UpdateHandler(a.MemStorage, a.FileProducer, a.DBConn))

		// value
		r.Get(`/value/{MetricType}/{MetricName}`, handlers.GetMetricHandler(a.MemStorage))
		r.Post(`/value`, handlers.GetHandler(a.MemStorage))

		// main page
		r.Get(`/`, handlers.MainPageHandler(a.MemStorage))
	})
}

func (a *ServerApp) createTables(ctx context.Context) error {
	if a.DBConn != nil {

		queryMetricTable := `create table if not exists metric(
								id varchar(200) unique not null, 
								"type" varchar(50) not null, 
								delta integer, 
								"value" double precision);`
		_, err := a.DBConn.Exec(ctx, queryMetricTable)
		if err != nil {
			log.Printf("can't execute query because of %s\n", err.Error())
			return err
		}
	}
	return nil
}

func (a *ServerApp) initFileManagers() error {
	producer, err := file.NewProducer(a.Config.FileStoragePath)
	if err != nil {
		log.Printf("can't initialize FileProducer because of: %s\n", err.Error())
		return err
	}

	consumer, err := file.NewConsumer(a.Config.FileStoragePath)
	if err != nil {
		log.Printf("can't initialize FileConsumer because of: %s\n", err.Error())
		return err
	}

	a.FileProducer = producer
	a.FileConsumer = consumer
	return nil
}
