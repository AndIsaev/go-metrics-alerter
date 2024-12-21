package handler

import (
	"encoding/json"
	"github.com/AndIsaev/go-metrics-alerter/internal/common"
	"github.com/AndIsaev/go-metrics-alerter/internal/service/server"
	"github.com/go-chi/chi"
	"github.com/mailru/easyjson"
	"net/http"
)

func (h *Handler) SetMetricHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		MetricType := chi.URLParam(r, "MetricType")
		MetricName := chi.URLParam(r, "MetricName")

		// check value is specified for the metric type
		if !server.IsCorrectType(MetricType) {
			http.Error(w, "An incorrect value is specified for the metric type", http.StatusBadRequest)
			return
		}
		MetricValue, err := server.DefineMetricValue(MetricType, chi.URLParam(r, "MetricValue"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		metric := common.Metrics{ID: MetricName, MType: MetricType}

		if err := h.MetricService.UpdateMetricByValue(r.Context(), metric, MetricValue); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
	}
}

// UpdateRowHandler - upsert metric
func (h *Handler) UpdateRowHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metric := common.Metrics{}
		w.Header().Set("Content-Type", "application/json")
		defer r.Body.Close()

		err := json.NewDecoder(r.Body).Decode(&metric)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if ok := metric.IsValid(); !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// save metrics to file
		updatedMetric, err := h.MetricService.InsertMetric(r.Context(), metric)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		result, _ := easyjson.Marshal(updatedMetric)

		w.WriteHeader(http.StatusOK)
		w.Write(result)
	}
}

func (h *Handler) UpdateBatchHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metrics := make([]common.Metrics, 0, 100)
		defer r.Body.Close()

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewDecoder(r.Body).Decode(&metrics); err != nil {
			w.WriteHeader(http.StatusBadRequest)

			body := common.Response{Message: err.Error()}
			response, _ := easyjson.Marshal(body)
			w.Write(response)

			return
		}

		err := h.MetricService.InsertMetrics(r.Context(), metrics)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		body := common.Response{Message: "success"}
		response, _ := easyjson.Marshal(body)

		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}
