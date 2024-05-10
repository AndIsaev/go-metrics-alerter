package metric

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"net/http"
)

func UpdateMetricRouter() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger, middleware.SetHeader("Content-Type", "text/plain"))

	// set value for metric
	r.Post("/{MetricType}/{MetricName}/{MetricValue}", UpdateMetricHandler)

	return r
}

func GetMetricRouter() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger, middleware.SetHeader("Content-Type", "text/plain"))

	// get value of metric
	r.Get("/{MetricType}/{MetricName}", GetMetricHandler)

	return r
}
