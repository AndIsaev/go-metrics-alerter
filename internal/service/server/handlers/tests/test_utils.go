package tests

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AndIsaev/go-metrics-alerter/internal/service/server/handlers"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/mailru/easyjson"

	"github.com/AndIsaev/go-metrics-alerter/internal/common"
	"github.com/AndIsaev/go-metrics-alerter/internal/logger"
	"github.com/AndIsaev/go-metrics-alerter/internal/manager/file"
	"github.com/AndIsaev/go-metrics-alerter/internal/service"
	mid "github.com/AndIsaev/go-metrics-alerter/internal/service/server/middleware"
	"github.com/AndIsaev/go-metrics-alerter/internal/storage"

	"github.com/stretchr/testify/require"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

type TestServerApp struct {
	Router       chi.Router
	MemStorage   *storage.MemStorage
	FileProducer *file.Producer
	FileConsumer *file.Consumer
	DBConn       storage.BaseStorage
	Config       *service.ServerConfig
	Server       *httptest.Server
}

func NewTestServerApp() *TestServerApp {
	testApp := &TestServerApp{}
	config := service.ServerConfig{}
	config.FileStoragePath = "./test_metrics"
	testApp.Config = &config

	fileProducer, _ := file.NewProducer(config.FileStoragePath)
	fileConsumer, _ := file.NewConsumer(config.FileStoragePath)

	testApp.FileProducer = fileProducer
	testApp.FileConsumer = fileConsumer

	testApp.MemStorage = storage.NewMemStorage()

	testApp.Router = chi.NewRouter()
	testApp.Server = httptest.NewServer(testApp.Router)

	testApp.initRouter()
	return testApp
}

func (a *TestServerApp) initRouter() {
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
		r.Post(`/value`, handlers.GetHandler(a.MemStorage, a.DBConn))

		// main page
		r.Get(`/`, handlers.MainPageHandler(a.MemStorage))
	})
}
