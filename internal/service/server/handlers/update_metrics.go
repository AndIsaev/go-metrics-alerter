package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/mailru/easyjson"

	"github.com/AndIsaev/go-metrics-alerter/internal/common"
	"github.com/AndIsaev/go-metrics-alerter/internal/manager/file"
	"github.com/AndIsaev/go-metrics-alerter/internal/service/server"
	"github.com/AndIsaev/go-metrics-alerter/internal/storage"
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

		val, err := server.DefineMetricValue(MetricType, chi.URLParam(r, "MetricValue"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		MetricValue = val

		if err := mem.Add(MetricType, MetricName, MetricValue); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
	}
}

// UpdateHandler - saving metrics from agent
func UpdateHandler(mem *storage.MemStorage, producer *file.Producer, conn storage.BaseStorage) http.HandlerFunc {
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

		// save metrics to file
		if conn != nil {
			err := conn.Insert(context.Background(), metrics)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		if producer != nil {
			if err := producer.Insert(&metrics); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		// save new metrics to db
		if err := mem.Set(&metrics); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		result, _ := easyjson.Marshal(metrics)

		w.WriteHeader(http.StatusOK)
		w.Write(result)
	}
}

func UpdateBatchHandler(mem *storage.MemStorage, producer *file.Producer, conn storage.BaseStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metrics := make([]common.Metrics, 0, 100)

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewDecoder(r.Body).Decode(&metrics); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// save metrics to file
		if conn != nil {
			err := conn.InsertBatch(context.Background(), &metrics)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		if producer != nil {
			if err := producer.InsertBatch(&metrics); err != nil {
				log.Println(errors.Unwrap(err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		// save new metrics to db
		if err := mem.InsertBatch(&metrics); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		body := common.Response{Message: "success"}
		response, _ := easyjson.Marshal(body)

		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}
