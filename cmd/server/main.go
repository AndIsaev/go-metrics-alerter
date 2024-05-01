package main

import (
	"github.com/AndIsaev/go-metrics-alerter/internal/service/server/handlers"
	"github.com/AndIsaev/go-metrics-alerter/internal/storage"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	ms := storage.NewMemStorage()

	mux.HandleFunc(`/update/`, func(w http.ResponseWriter, r *http.Request) {
		handlers.UpdateMetricHandler(ms, w, r)
	})

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
