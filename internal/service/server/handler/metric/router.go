package metric

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"net/http"
)

func MetricRouter() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger, middleware.SetHeader("Content-Type", "text/plain"))

	r.Post("/{MetricType}/{MetricName}/{MetricValue}", UpdateMetricHandler)

	return r
}
