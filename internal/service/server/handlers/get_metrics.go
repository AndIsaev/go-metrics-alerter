package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/AndIsaev/go-metrics-alerter/internal/common/models"
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

func GetMetricsHandler(w http.ResponseWriter, r *http.Request) {

	var response models.MetricsResponse
	var request models.MetricsRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !server.IsCorrectType(request.MType) {
		http.Error(w, "An incorrect value is specified for the metric type", http.StatusBadRequest)
		return
	}

	if val, err := storage.MS.Get(request.ID); err == nil {
		response.ID = request.ID
		response.MType = request.MType
		response.Value = val
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	resp, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)

}
