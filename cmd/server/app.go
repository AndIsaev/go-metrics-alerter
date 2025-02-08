package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"google.golang.org/grpc"

	pb "github.com/AndIsaev/go-metrics-alerter/internal/service/proto"
	"github.com/AndIsaev/go-metrics-alerter/internal/service/server"
	"github.com/AndIsaev/go-metrics-alerter/internal/service/server/rpc"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/mailru/easyjson"

	"github.com/AndIsaev/go-metrics-alerter/internal/common"
	"github.com/AndIsaev/go-metrics-alerter/internal/logger"
	"github.com/AndIsaev/go-metrics-alerter/internal/service/server/handler"
	mid "github.com/AndIsaev/go-metrics-alerter/internal/service/server/middleware"
	"github.com/AndIsaev/go-metrics-alerter/internal/storage"
	"github.com/AndIsaev/go-metrics-alerter/internal/storage/file"
	"github.com/AndIsaev/go-metrics-alerter/internal/storage/inmemory"
	"github.com/AndIsaev/go-metrics-alerter/internal/storage/postgres"
)

// ServerApp - structure of application
type ServerApp struct {
	Router     chi.Router
	Conn       storage.Storage
	Config     *Config
	Server     *http.Server
	GRPCServer *grpc.Server
	Handler    *handler.Handler
	fm         *file.Manager
	wg         sync.WaitGroup
	chMetrics  chan []common.Metrics
}

// New - create new app
func New() *ServerApp {
	app := &ServerApp{}
	app.Config = NewConfig()
	app.chMetrics = make(chan []common.Metrics)
	app.Router = chi.NewRouter()
	app.Handler = &handler.Handler{}

	return app
}

// StartApp start application
func (a *ServerApp) StartApp() error {
	if err := logger.Initialize(); err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer a.shutdown(ctx) // закрываем соединения

	err := a.initStorage(ctx)
	if err != nil {
		log.Printf("error init storage %v\n", err)
		return err
	}

	if a.fm != nil && a.Config.StoreInterval != 0 {
		a.wg.Add(2)
		go a.pullMetrics(ctx)
		go a.runFileWorker(ctx)
	}

	a.wg.Add(1)
	go a.runServer(ctx)

	// обработка сигналов завершения
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	// горутина для прослушивания сигналов
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		sig := <-sigs
		log.Printf("Received signal: %s", sig)
		cancel() // отменяем контекст
	}()

	a.wg.Wait()
	close(a.chMetrics)

	return nil
}

// startHTTPServer - start http server
func (a *ServerApp) startHTTPServer() error {
	log.Printf("start http server on %s\n", a.Config.Address)
	return a.Server.ListenAndServe()
}

// startGRPCServer - start grpc server
func (a *ServerApp) startGRPCServer() error {
	log.Printf("start grpc server on %s\n", a.Config.Address)
	listen, err := net.Listen("tcp", a.Config.Address)
	if err != nil {
		return err
	}

	return a.GRPCServer.Serve(listen)
}

// Shutdown close connections
func (a *ServerApp) shutdown(ctx context.Context) {
	if a.fm != nil {
		a.fm.Close()
	}
	if a.Conn != nil {
		err := a.Conn.System().Close(ctx)
		if err != nil {
			log.Printf("error close storage: %v\n", err.Error())
			return
		}
	}
}

// initRouter - initialize new router
func (a *ServerApp) initRouter() {
	r := a.Router
	r.Use(logger.RequestLogger, logger.ResponseLogger)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.SetHeader("Content-Type", "application/json"))
	if a.Config.PrivateKey != "" {
		privateKey, err := a.Config.GetPrivateKey()
		if err != nil {
			log.Printf("Error get private key: %v\n", err)
		}
		r.Use(mid.DecryptMiddleware(privateKey))
	}
	if a.Config.TrustedSubnet != "" {
		r.Use(mid.TrustedSubnetMiddleware(a.Config.TrustedSubnet))
	}

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
		r.Use(mid.GzipMiddleware, mid.SecretMiddleware(a.Config.Key))
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

func (a *ServerApp) initStorage(ctx context.Context) error {
	select {
	case <-ctx.Done():
		log.Println("context done -> exit from initStorage")
		return nil
	default:
		if a.Config.Dsn != "" {
			conn, err := postgres.NewPgStorage(a.Config.Dsn)
			if err != nil {
				return err
			}

			a.Conn = conn
		} else {
			var syncFileManager *file.Manager
			var syncSave = false
			if a.Config.FileStoragePath != "" {
				if err := a.fm.CreateDir(a.Config.FileStoragePath); err != nil {
					log.Printf("error create directory: %s\n", err.Error())
					return err
				}

				fileManager, err := file.NewManager(a.Config.FileStoragePath)
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
	}

	return nil
}

func (a *ServerApp) saveMetricsToDisc(ctx context.Context, metrics []common.Metrics) error {
	select {
	case <-ctx.Done():
		log.Println("context done -> exit from saveMetricsToDisc")
		return nil
	default:
		err := a.fm.Overwrite(ctx, metrics)
		if err != nil {
			return err
		}
		log.Printf("save %v metrics to disc\n", len(metrics))
		return nil
	}
}

func (a *ServerApp) runServer(ctx context.Context) {
	defer a.wg.Done()

	// Создаем канал для ошибок сервера
	serverErrChan := make(chan error, 1)

	// Запускаем сервер в отдельной горутине
	go func() {
		if !a.Config.IsRPC {
			a.Handler.MetricService = &server.Methods{Storage: a.Conn}
			// init http router
			a.initRouter()
			// params for http server
			a.Server = &http.Server{Handler: a.Router, Addr: a.Config.Address}

			if err := a.startHTTPServer(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				serverErrChan <- fmt.Errorf("HTTP server error: %w", err)
			}
		} else {
			// init GRPC server
			a.GRPCServer = grpc.NewServer()
			// set params for GRPC interface
			pb.RegisterMetricServiceServer(a.GRPCServer, &rpc.MetricServiceServer{Storage: a.Conn})

			if err := a.startGRPCServer(); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
				serverErrChan <- fmt.Errorf("GRPC server error: %w", err)
			}
		}

		close(serverErrChan) // Закрываем канал после завершения работы сервера
	}()

	select {
	case err := <-serverErrChan:
		if err != nil {
			log.Printf("got error from server chan: %v\n", err)
		}
	case <-ctx.Done():
		// Если контекст отменен, корректно завершаем работу сервера
		shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelShutdown()

		if a.Config.IsRPC {
			a.GRPCServer.GracefulStop()
			log.Println("GRPC server shutdown gracefully")
			return
		}
		if err := a.Server.Shutdown(shutdownCtx); err != nil {
			log.Printf("HTTP server shutdown error: %v\n", err)
		} else {
			log.Println("HTTP server shutdown gracefully")
		}
	}
}

func (a *ServerApp) pullMetrics(ctx context.Context) {
	defer a.wg.Done()

	ticker := time.NewTicker(a.Config.StoreInterval) // Используем ticker для интервальных операций
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Shutting down metrics collection goroutine")
			return
		case <-ticker.C:
			metrics, err := a.Conn.Metric().List(ctx)
			if err != nil {
				log.Printf("failed пуе list metrics: %v", err)
				continue
			}

			select {
			case a.chMetrics <- metrics:
			case <-ctx.Done():
				return
			}
		}
	}
}

func (a *ServerApp) runFileWorker(ctx context.Context) {
	defer a.wg.Done()
	for {
		select {
		case <-ctx.Done():
			log.Println("Shutting down metrics saving goroutine")
			return
		case metrics := <-a.chMetrics:
			if err := a.saveMetricsToDisc(ctx, metrics); err != nil {
				log.Printf("error saving metrics to disk: %v", err)
			}
		}
	}
}
