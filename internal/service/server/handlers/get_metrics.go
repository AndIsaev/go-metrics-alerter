package handlers

import (
	"fmt"
	"github.com/AndIsaev/go-metrics-alerter/internal/service/server"
	"github.com/AndIsaev/go-metrics-alerter/internal/storage"
	"github.com/go-chi/chi"
	"net/http"
)

func GetMetricHandler(w http.ResponseWriter, r *http.Request) {
	MetricType := chi.URLParam(r, "MetricType")
	MetricName := chi.URLParam(r, "MetricName")

	// check value is specified for the metric type
	if !server.IsCorrectType(MetricType) {
		http.Error(w, "An incorrect value is specified for the metric type", http.StatusBadRequest)
		return
	}
	if val, err := storage.MS.Get(MetricName); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else {
		w.Write([]byte(fmt.Sprintf("%v", val)))
	}

}
