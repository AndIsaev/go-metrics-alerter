package api

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"net/http"
)

func UpdateMetricRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.SetHeader("Content-Type", "text/plain"))

	// set value for metric
	r.Post("/{MetricType}/{MetricName}/{MetricValue}", SetMetricHandler)

	return r
}

func GetMetricRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.SetHeader("Content-Type", "text/plain"))

	// get value of metric
	r.Get("/{MetricType}/{MetricName}", GetMetricHandler)

	return r
}

func MainPageRouter() http.Handler {
	r := chi.NewRouter()

	r.Get("/", MainPageHandler)

	return r
}
