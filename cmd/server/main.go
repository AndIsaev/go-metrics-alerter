package main

import (
	"fmt"
	"github.com/AndIsaev/go-metrics-alerter/internal/service/server/api"
	"github.com/go-chi/chi/middleware"
	"net/http"

	"github.com/go-chi/chi"
)

func run() error {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Mount(`/`, api.MainPageRouter())
	r.Mount(`/update/`, api.UpdateMetricRouter())
	r.Mount(`/value/`, api.GetMetricRouter())

	fmt.Println("Running server on", flagRunAddr)
	return http.ListenAndServe(flagRunAddr, r)
}

func main() {
	parseFlags()

	if err := run(); err != nil {
		panic(err)
	}
}
