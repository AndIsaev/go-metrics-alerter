package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/AndIsaev/go-metrics-alerter/internal/common"
	"github.com/AndIsaev/go-metrics-alerter/internal/common/models"
	"github.com/AndIsaev/go-metrics-alerter/internal/service/server"
	"github.com/AndIsaev/go-metrics-alerter/internal/storage"
	"github.com/go-chi/chi"

	"net/http"
)

func SetMetricHandler(w http.ResponseWriter, r *http.Request) {
	var MetricValue interface{}
	MetricType := chi.URLParam(r, "MetricType")
	MetricName := chi.URLParam(r, "MetricName")

	// check value is specified for the metric type
	if !server.IsCorrectType(MetricType) {
		http.Error(w, "An incorrect value is specified for the metric type", http.StatusBadRequest)
		return
	}

	if val, err := server.DefineMetricValue(MetricType, chi.URLParam(r, "MetricValue")); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else {
		MetricValue = val
	}

	if err := storage.MS.Add(MetricType, MetricName, MetricValue); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/plain")

}

func SetMetricsHandler(w http.ResponseWriter, r *http.Request) {

	var metric models.Metrics
	var MetricValue interface{}
	var response models.MetricsResponse

	if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// check value is specified for the metric type
	if !server.IsCorrectType(metric.MType) {
		http.Error(w, "An incorrect value is specified for the metric type", http.StatusBadRequest)
		return
	}

	switch metric.MType {
	case common.Gauge:
		MetricValue = *metric.Value
	case common.Counter:
		MetricValue = *metric.Delta
	}

	if err := storage.MS.Add(metric.MType, metric.ID, MetricValue); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if val, err := storage.MS.Get(metric.ID); err == nil {
		response.ID = metric.ID
		response.MType = metric.MType
		response.Value = val
	} else {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	resp, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)

}
