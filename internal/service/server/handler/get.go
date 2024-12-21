package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/AndIsaev/go-metrics-alerter/internal/common"
	"github.com/AndIsaev/go-metrics-alerter/internal/service/server"
	"github.com/go-chi/chi"
	"github.com/mailru/easyjson"
	"log"
	"net/http"
)

func (h *Handler) GetByURLParamHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		MetricType := chi.URLParam(r, "MetricType")
		MetricName := chi.URLParam(r, "MetricName")

		// check value is specified for the metric type
		if !server.IsCorrectType(MetricType) {
			http.Error(w, "An incorrect value is specified for the metric type", http.StatusBadRequest)
			return
		}

		metric, err := h.MetricService.GetMetricByName(r.Context(), MetricName)
		if err != nil {
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

func (h *Handler) GetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metric := common.Metrics{}
		w.Header().Set("Content-Type", "application/json")
		defer r.Body.Close()

		if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
			log.Println(err)
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
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}
