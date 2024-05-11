package main

import (
	"github.com/AndIsaev/go-metrics-alerter/internal/service/server/handler/metric"
	"github.com/go-chi/chi/middleware"
	"net/http"

	"github.com/go-chi/chi"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get(`/`, metric.MainPageHandler)
	r.Mount(`/update/`, metric.UpdateMetricRouter())
	r.Mount(`/value/`, metric.GetMetricRouter())

	err := http.ListenAndServe(`:8080`, r)
	if err != nil {
		panic(err)
	}
}
