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
	if val, err := storage.MS.GetV1(MetricType + "-" + MetricName); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else {
		w.Write([]byte(fmt.Sprintf("%v", val)))
	}
	w.Header().Set("Content-Type", "text/plain")

}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	var metrics common.Metrics
	w.Header().Set("Content-Type", "application/json")

	err := json.NewDecoder(r.Body).Decode(&metrics)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !server.IsCorrectType(metrics.MType) {
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

	e := storage.MS.Get(&metrics)
	if e != nil {
		resp := common.Response{
			Status: http.StatusNotFound,
			Text:   e.Error(),
		}
		answer, ee := json.Marshal(resp)
		if ee != nil {
			http.Error(w, ee.Error(), http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusNotFound)
		w.Write(answer)
		return
	}

	resp, err := json.Marshal(metrics)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)

}
