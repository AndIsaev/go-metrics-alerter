package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/mailru/easyjson"

	"github.com/AndIsaev/go-metrics-alerter/internal/common"
	"github.com/AndIsaev/go-metrics-alerter/internal/logger"
	"github.com/AndIsaev/go-metrics-alerter/internal/manager/file"
	mid "github.com/AndIsaev/go-metrics-alerter/internal/service/server/middleware"
	"github.com/AndIsaev/go-metrics-alerter/internal/storage"
)

func ServerRouter(memory *storage.MemStorage, fileProducer *file.Producer, DbConn storage.PgStorage) chi.Router {
	r := chi.NewRouter()
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
		r.Get(`/ping`, PingHandler(DbConn))

		// update
		r.Post(`/update/{MetricType}/{MetricName}/{MetricValue}`, SetMetricHandler(memory))
		r.Post(`/update`, UpdateHandler(memory, fileProducer))

		// value
		r.Get(`/value/{MetricType}/{MetricName}`, GetMetricHandler(memory))
		r.Post(`/value`, GetHandler(memory))

		// main page
		r.Get(`/`, MainPageHandler(memory))
	})

	return r
}
