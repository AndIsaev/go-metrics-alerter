package main

import (
	"fmt"
	"github.com/AndIsaev/go-metrics-alerter/internal/service/server/handler/metric"
	"github.com/go-chi/chi/middleware"
	"net/http"

	"github.com/go-chi/chi"
)

func main() {
	parseFlags()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get(`/`, metric.MainPageHandler)
	r.Mount(`/update/`, metric.UpdateMetricRouter())
	r.Mount(`/value/`, metric.GetMetricRouter())

	fmt.Println("Running server on", flagRunAddr)
	err := http.ListenAndServe(flagRunAddr, r)
	if err != nil {
		panic(err)
	}
}
