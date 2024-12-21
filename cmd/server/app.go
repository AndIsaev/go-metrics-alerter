package main

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"net/http/pprof"
	"os"

	"github.com/AndIsaev/go-metrics-alerter/internal/service/server"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/mailru/easyjson"

	"github.com/AndIsaev/go-metrics-alerter/internal/common"
	"github.com/AndIsaev/go-metrics-alerter/internal/logger"
	"github.com/AndIsaev/go-metrics-alerter/internal/manager/file"
	"github.com/AndIsaev/go-metrics-alerter/internal/service/server/handlers"
	mid "github.com/AndIsaev/go-metrics-alerter/internal/service/server/middleware"
	"github.com/AndIsaev/go-metrics-alerter/internal/storage"
)

// ServerApp - structure of application
type ServerApp struct {
	Router       chi.Router
	MemStorage   *storage.MemStorage
	FileProducer *file.Producer
	FileConsumer *file.Consumer
	DBConn       storage.BaseStorage
	Config       *server.Config
	Server       *http.Server
}

// New - create new app
func New() *ServerApp {
	app := &ServerApp{}
	config := server.NewConfig()
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
		// connect to db
		conn, err := storage.NewPostgresStorage(a.Config.DBDsn)
		if err != nil {
			return err
		}

		// создаем таблицы
		if err := conn.Create(ctx); err != nil {
			return err
		}
		a.DBConn = conn
	}

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
	if err := a.downloadMetrics(); err != nil {
		return err
	}

	// init router
	a.initRouter()

	// init http server
	a.initHTTPServer()

	return a.startHTTPServer()
}

// startHTTPServer - start http server
func (a *ServerApp) startHTTPServer() error {
	log.Printf("start server on %s\n", a.Config.Address)
	return a.Server.ListenAndServe()
}

// initHTTPServer - init http server
func (a *ServerApp) initHTTPServer() {
	a.Server = &http.Server{Handler: a.Router, Addr: a.Config.Address}
}

// downloadMetrics - Read metrics from disk
func (a *ServerApp) downloadMetrics() error {
	if a.Config.Restore {
		log.Println("read metrics from disk")
		for {
			m, err := a.FileConsumer.ReadMetrics()
			if err != nil {
				break
			}
			if err := a.MemStorage.Set(m); err != nil {
				log.Printf("can't save metrics to local storage because of: %s\n", err.Error())
				return err
			}
		}
		log.Println("metrics downloaded")
	}
	return nil
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

func (a *ServerApp) Shutdown() {
	if err := a.FileProducer.Close(); err != nil {
		log.Printf("%s\n", err.Error())
	}
	if err := a.FileConsumer.Close(); err != nil {
		log.Printf("%s\n", err.Error())
	}

	if a.Config.DBDsn != "" {
		if a.DBConn != nil {
			if err := a.DBConn.Close(); err != nil {
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

	// Pprof routes
	r.HandleFunc("/debug/pprof", pprof.Index)
	r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	r.HandleFunc("/debug/pprof/trace", pprof.Trace)

	r.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	r.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	r.Handle("/debug/pprof/mutex", pprof.Handler("mutex"))
	r.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	r.Handle("/debug/pprof/block", pprof.Handler("block"))
	r.Handle("/debug/pprof/allocs", pprof.Handler("allocs"))

	// Routes
	r.Group(func(r chi.Router) {
		r.Use(mid.GzipMiddleware, a.secretMiddleware)
		r.Post(`/updates`, handlers.UpdateBatchHandler(a.MemStorage, a.FileProducer, a.DBConn))
	})

	r.Group(func(r chi.Router) {
		r.Use(mid.GzipMiddleware)

		// Ping db connection
		r.Get(`/ping`, handlers.PingHandler(a.DBConn))

		// update
		r.Post(`/update/{MetricType}/{MetricName}/{MetricValue}`, handlers.SetMetricHandler(a.MemStorage))
		r.Post(`/update`, handlers.UpdateHandler(a.MemStorage, a.FileProducer, a.DBConn))

		// value
		r.Get(`/value/{MetricType}/{MetricName}`, handlers.GetMetricHandler(a.MemStorage))
		r.Post(`/value`, handlers.GetHandler(a.MemStorage, a.DBConn))

		// main page
		r.Get(`/`, handlers.MainPageHandler(a.MemStorage))
	})
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

func (a *ServerApp) secretMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if a.Config.Key != "" {
			agentSha256sum := r.Header.Get("HashSHA256")

			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(rw, err.Error(), http.StatusBadRequest)
				return
			}

			r.Body = io.NopCloser(bytes.NewBuffer(body))

			defer r.Body.Close()

			serverSha256sum := common.Sha256sum(body, a.Config.Key)

			if agentSha256sum != serverSha256sum {
				log.Printf("compare hash is not success")
				rw.WriteHeader(http.StatusBadRequest)
				return
			}
			rw.Header().Set("HashSHA256", serverSha256sum)
		}

		next.ServeHTTP(rw, r)
	})
}
