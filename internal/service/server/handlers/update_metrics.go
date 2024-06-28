package handlers

import (
	"encoding/json"
	"github.com/AndIsaev/go-metrics-alerter/internal/common"
	"github.com/AndIsaev/go-metrics-alerter/internal/service/server"
	"github.com/AndIsaev/go-metrics-alerter/internal/storage"
	"github.com/go-chi/chi"
	"github.com/mailru/easyjson"

	"net/http"
)

func SetMetricHandler(mem *storage.MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		if err := mem.Add(MetricType, MetricName, MetricValue); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
	}
}

func UpdateHandler(mem *storage.MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metrics := common.Metrics{}
		w.Header().Set("Content-Type", "application/json")

		err := json.NewDecoder(r.Body).Decode(&metrics)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if ok := metrics.IsValid(); !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		mem.Set(&metrics)

		result, _ := easyjson.Marshal(metrics)

		w.WriteHeader(http.StatusOK)
		w.Write(result)
	}
}
