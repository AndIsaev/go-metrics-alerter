package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/mailru/easyjson"

	"github.com/AndIsaev/go-metrics-alerter/internal/common"
)

// GetByURLParamHandler get value of metric by type and name
func (h *Handler) GetByURLParamHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		MetricType := chi.URLParam(r, "MetricType")
		MetricName := chi.URLParam(r, "MetricName")

		// check value is specified for the metric type
		if !IsCorrectType(MetricType) {
			http.Error(w, "an incorrect value is specified for the metric type", http.StatusBadRequest)
			return
		}

		metric, err := h.MetricService.GetMetricByName(r.Context(), MetricName)
		if err != nil {
			log.Println(errors.Unwrap(err))
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		var val string
		if metric.Delta != nil {
			val = fmt.Sprintf("%v", *metric.Delta)
		} else {
			val = fmt.Sprintf("%v", *metric.Value)
		}

		w.Write([]byte(val))
	}
}

// GetHandler get metric by json
func (h *Handler) GetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metric := common.Metrics{}
		w.Header().Set("Content-Type", "application/json")
		defer r.Body.Close()

		if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
			log.Println(errors.Unwrap(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if !metric.IsValidType() {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		value, err := h.MetricService.GetMetricByName(r.Context(), metric.ID)
		if err != nil {
			log.Println(errors.Unwrap(err))
			w.WriteHeader(http.StatusNotFound)
			return
		}

		response, err := easyjson.Marshal(value)
		if err != nil {
			log.Println(errors.Unwrap(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}
