package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/mailru/easyjson"

	"github.com/AndIsaev/go-metrics-alerter/internal/common"
	"github.com/AndIsaev/go-metrics-alerter/internal/service/server"
	"github.com/AndIsaev/go-metrics-alerter/internal/storage"
)

func GetMetricHandler(mem *storage.MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		MetricType := chi.URLParam(r, "MetricType")
		MetricName := chi.URLParam(r, "MetricName")

		// check value is specified for the metric type
		if !server.IsCorrectType(MetricType) {
			http.Error(w, "An incorrect value is specified for the metric type", http.StatusBadRequest)
			return
		}

		val, err := mem.GetMetricByName(MetricName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.Write([]byte(fmt.Sprintf("%v", val)))
	}
}

func GetHandler(m *storage.MemStorage, conn storage.BaseStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metric := common.Metrics{}
		w.Header().Set("Content-Type", "application/json")

		if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if !metric.IsValidType() {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// get metric
		if conn != nil {
			val, err := conn.Get(context.Background(), metric)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			metric = *val
		} else {
			val, err := m.GetMetric(metric.MType, metric.ID)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			metric = val
		}

		response, _ := easyjson.Marshal(metric)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}
