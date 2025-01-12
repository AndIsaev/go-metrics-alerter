package main

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"net/http/pprof"
	"sync"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/mailru/easyjson"

	"github.com/AndIsaev/go-metrics-alerter/internal/common"
	"github.com/AndIsaev/go-metrics-alerter/internal/logger"
	"github.com/AndIsaev/go-metrics-alerter/internal/service/server"
	"github.com/AndIsaev/go-metrics-alerter/internal/service/server/handler"
	mid "github.com/AndIsaev/go-metrics-alerter/internal/service/server/middleware"
	"github.com/AndIsaev/go-metrics-alerter/internal/storage"
	"github.com/AndIsaev/go-metrics-alerter/internal/storage/file"
	"github.com/AndIsaev/go-metrics-alerter/internal/storage/inmemory"
	"github.com/AndIsaev/go-metrics-alerter/internal/storage/postgres"
)

// ServerApp - structure of application
type ServerApp struct {
	Router  chi.Router
	Conn    storage.Storage
	Config  *Config
	Server  *http.Server
	Handler *handler.Handler
	fm      *file.FileManager
}

// New - create new app
func New() *ServerApp {
	app := &ServerApp{}
	app.Config = NewConfig()
	ctx := context.Background()

	// init file storage
	err := app.initStorage(ctx)
	if err != nil {
		log.Fatalf("error init storage%v\n", err.Error())
	}

	app.Router = chi.NewRouter()

	app.Handler = &handler.Handler{}
	app.Handler.MetricService = &server.Methods{Storage: app.Conn}

	return app
}

func (a *ServerApp) StartApp(ctx context.Context) error {
	if err := logger.Initialize(); err != nil {
		return err
	}
	var wg sync.WaitGroup
	var mu sync.RWMutex
	chMetrics := make(chan []common.Metrics)

	// init router
	a.initRouter()

	// init http server
	a.initHTTPServer()

	if a.fm != nil && a.Config.StoreInterval != 0 {
		wg.Add(2)

		go func(ctx context.Context, ch chan []common.Metrics) {
			defer wg.Done()
			for {
				time.Sleep(a.Config.StoreInterval)
				mu.RLock()
				metrics, _ := a.Conn.Metric().List(ctx)
				mu.RUnlock()
				ch <- metrics
			}
		}(ctx, chMetrics)

		go func() {
			defer wg.Done()
			for {
				time.Sleep(a.Config.StoreInterval)
				err := a.saveMetricsToDisc(chMetrics)
				if err != nil {
					log.Printf("error save metrics to disc")
				}
			}
		}()
	}
	wg.Add(1)

	go func() {
		_ = a.startHTTPServer()
	}()
	wg.Wait()
	return nil
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

func (a *ServerApp) Shutdown() {
	ctx := context.Background()
	if a.fm != nil {
		a.fm.Close()
	}
	err := a.Conn.System().Close(ctx)
	if err != nil {
		log.Printf("error close storage: %v\n", err.Error())
		return
	}
}

// initRouter - initialize new router
func (a *ServerApp) initRouter() {
	r := a.Router
	r.Use(logger.RequestLogger, logger.ResponseLogger)
	r.Use(middleware.Recoverer, middleware.StripSlashes)
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
		r.Post(`/updates`, a.Handler.UpdateBatchHandler())
	})

	r.Group(func(r chi.Router) {
		r.Use(mid.GzipMiddleware)

		// Ping db connection
		r.Get(`/ping`, a.Handler.PingHandler())

		//// update
		r.Post(`/update/{MetricType}/{MetricName}/{MetricValue}`, a.Handler.SetMetricHandler())
		r.Post(`/update`, a.Handler.UpdateRowHandler())

		// value
		r.Get(`/value/{MetricType}/{MetricName}`, a.Handler.GetByURLParamHandler())
		r.Post(`/value`, a.Handler.GetHandler())

		// main page
		r.Get(`/`, a.Handler.IndexHandler())
	})
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
				log.Println("compare hash is not success")
				rw.WriteHeader(http.StatusBadRequest)
				return
			}
			rw.Header().Set("HashSHA256", serverSha256sum)
		}

		next.ServeHTTP(rw, r)
	})
}

func (a *ServerApp) initStorage(ctx context.Context) error {
	if a.Config.Dsn != "" {
		conn, err := postgres.NewPgStorage(a.Config.Dsn)
		if err != nil {
			return err
		}

		a.Conn = conn
	} else {
		var syncFileManager *file.FileManager
		var syncSave = false
		if a.Config.FileStoragePath != "" {
			if err := a.fm.CreateDir(a.Config.FileStoragePath); err != nil {
				log.Printf("error create directory: %s\n", err.Error())
				return err
			}

			fileManager, err := file.NewFileManager(a.Config.FileStoragePath)
			if err != nil {
				log.Printf("error init file manager")
				return err
			}
			a.fm = fileManager

			if a.Config.StoreInterval == 0 {
				syncFileManager = fileManager
				syncSave = true
			}
			if a.Config.Restore {
				syncFileManager = fileManager
			}
		}
		a.Conn = inmemory.NewMemStorage(syncFileManager, syncSave)
	}

	if a.Config.Restore {
		if err := a.Conn.System().RunMigrations(ctx); err != nil {
			log.Printf("error run migrations for storage: %v\n", err.Error())
			return err
		}
	}

	return nil
}

func (a *ServerApp) saveMetricsToDisc(ch chan []common.Metrics) error {
	metrics := <-ch
	err := a.fm.Overwrite(metrics)
	if err != nil {
		return err
	}
	log.Printf("save %v metrics to disc\n", len(metrics))
	return nil
}
