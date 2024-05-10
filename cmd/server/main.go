package main

import (
	"github.com/AndIsaev/go-metrics-alerter/internal/service/server/handler/metric"
	"net/http"

	"github.com/go-chi/chi"
)

func main() {
	r := chi.NewRouter()
	r.Mount(`/update/`, metric.MetricRouter())

	err := http.ListenAndServe(`:8080`, r)
	if err != nil {
		panic(err)
	}
}
