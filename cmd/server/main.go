package main

import (
	"github.com/AndIsaev/go-metrics-alerter/internal/service/server/handlers"
	"github.com/AndIsaev/go-metrics-alerter/internal/storage"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	ms := storage.NewMemStorage()

	mux.HandleFunc(`/update/`, handlers.UpdateMetricHandler(ms))

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
