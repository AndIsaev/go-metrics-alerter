package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-chi/chi"
	"github.com/mailru/easyjson"

	"github.com/AndIsaev/go-metrics-alerter/internal/common"
	"github.com/AndIsaev/go-metrics-alerter/internal/manager/file"
	"github.com/AndIsaev/go-metrics-alerter/internal/service/server"
	"github.com/AndIsaev/go-metrics-alerter/internal/storage"

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
func UpdateHandler(mem *storage.MemStorage, producer *file.Producer, conn storage.PgStorage) http.HandlerFunc {
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
			exec, err := conn.Exec(
				context.Background(),
				`insert into metric (id, type, delta, value)
								  values ($1, $2, $3, $4)	`, metrics.ID, metrics.MType, metrics.Delta, metrics.Value)
			fmt.Println(exec, err)
			if err != nil {
				return
			}
		}

		if producer != nil {
			if err := server.SaveMetricsOnFile(producer, metrics); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		// save new metrics to DB
		mem.Set(&metrics)

		result, _ := easyjson.Marshal(metrics)

		w.WriteHeader(http.StatusOK)
		w.Write(result)
	}
}
