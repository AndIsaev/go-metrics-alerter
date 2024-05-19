package main

import (
	"fmt"
	"github.com/AndIsaev/go-metrics-alerter/internal/service"
	"github.com/AndIsaev/go-metrics-alerter/internal/service/server/handlers"

	"github.com/go-chi/chi/middleware"
	"net/http"

	"github.com/go-chi/chi"
)

func run() error {
	r := chi.NewRouter()
	config := service.NewServerConfig()

	r.Use(middleware.Logger)

	r.Mount(`/`, handlers.MainPageRouter())
	r.Mount(`/update/`, handlers.UpdateMetricRouter())
	r.Mount(`/value/`, handlers.GetMetricRouter())

	fmt.Println("Running server on", config.Address)
	return http.ListenAndServe(config.Address, r)

}

func main() {

	if err := run(); err != nil {
		panic(err)
	}
}
