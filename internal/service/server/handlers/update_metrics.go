package handlers

import (
	"encoding/json"
	"github.com/AndIsaev/go-metrics-alerter/internal/common"
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

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	var body common.Metrics

	w.Header().Set("Content-Type", "application/json")

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !server.IsCorrectType(body.MType) {
		resp := common.Response{
			Status: http.StatusBadRequest,
			Text:   "An incorrect value is specified for the metric type",
		}

		answer, e := json.Marshal(resp)
		if e != nil {
			http.Error(w, e.Error(), http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(answer)

		return
	}

	storage.MS.Set(&body)

	result, err := json.Marshal(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(result)
}
