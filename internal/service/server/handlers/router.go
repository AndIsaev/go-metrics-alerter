package handlers

import (
	"encoding/json"
	"github.com/AndIsaev/go-metrics-alerter/internal/common"
	"github.com/AndIsaev/go-metrics-alerter/internal/logger"
	"github.com/AndIsaev/go-metrics-alerter/internal/storage"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"net/http"
)

func ServerRouter(memory *storage.MemStorage) chi.Router {
	r := chi.NewRouter()
	r.Use(logger.RequestLogger, logger.ResponseLogger)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.Recoverer)
	r.Use(middleware.SetHeader("Content-Type", "application/json"))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		body := common.Response{Message: "route does not exist"}
		response, _ := json.Marshal(body)
		w.Write(response)
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		body := common.Response{Message: "method is not valid"}
		response, _ := json.Marshal(body)
		w.Write(response)
	})

	// Routes

	r.Group(func(r chi.Router) {
		// update
		r.Post(`/update/{MetricType}/{MetricName}/{MetricValue}`, SetMetricHandler(memory))
		r.Post(`/update`, UpdateHandler(memory))

		// value
		r.Get(`/value/{MetricType}/{MetricName}`, GetMetricHandler(memory))
		r.Post(`/value`, GetHandler(memory))

		// main page
		r.Get(`/`, MainPageHandler(memory))
	})

	return r
}
