package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/mailru/easyjson"

	"github.com/AndIsaev/go-metrics-alerter/internal/common"
)

// SetMetricHandler set value for metric by params
func (h *Handler) SetMetricHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		MetricType := chi.URLParam(r, "MetricType")
		MetricName := chi.URLParam(r, "MetricName")

		// check value is specified for the metric type
		if !IsCorrectType(MetricType) {
			http.Error(w, "an incorrect value is specified for the metric type", http.StatusBadRequest)
			return
		}
		MetricValue, err := DefineMetricValue(MetricType, chi.URLParam(r, "MetricValue"))
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

// UpdateRowHandler - set value for one metric by json
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
		if !metric.IsValidType() || !metric.IsValidValue() {
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

// UpdateBatchHandler set metrics by batch
func (h *Handler) UpdateBatchHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metrics := make([]common.Metrics, 0)
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
