package handlers

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"net/http"
)

func UpdateMetricRouter() http.Handler {
	r := chi.NewRouter()

	// set value for metric
	r.Post("/{MetricType}/{MetricName}/{MetricValue}", SetMetricHandler)
	r.Post("/", UpdateHandler)

	return r
}

func GetMetricRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.SetHeader("Content-Type", "text/plain"))

	// get value of metric
	r.Get("/{MetricType}/{MetricName}", GetMetricHandler)
	r.Post("/", GetHandler)

	return r
}

func MainPageRouter() http.Handler {
	r := chi.NewRouter()

	r.Get("/", MainPageHandler)

	return r
}
