package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/AndIsaev/go-metrics-alerter/internal/common"
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
	if val, err := storage.MS.Get(MetricType + "-" + MetricName); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else {
		w.Write([]byte(fmt.Sprintf("%v", val)))
	}

}
func GetHandler(w http.ResponseWriter, r *http.Request) {
	metrics := common.Metrics{}
	w.Header().Set("Content-Type", "application/json")

	err := json.NewDecoder(r.Body).Decode(&metrics)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !metrics.IsValidType() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	e := storage.MS.GetV1(&metrics)
	if e != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	resp, _ := json.Marshal(metrics)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)

}
