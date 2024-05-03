package handlers

import (
	"github.com/AndIsaev/go-metrics-alerter/internal/service/server"
	"github.com/AndIsaev/go-metrics-alerter/internal/storage"
	"net/http"
	"strings"
)

func UpdateMetricHandler(ms *storage.MemStorage, w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// get params
	params := r.URL.Path[len("/update/"):]
	parts := strings.Split(params, "/")

	// check correct format URL
	if len(parts) != 3 {
		http.Error(w, "Not correct format URL", http.StatusNotFound)
		return
	}

	metricType := parts[0]
	metricName := parts[1]
	metricValue := parts[2]

	// check value is specified for the metric type
	if !server.IsCorrectType(metricType) {
		http.Error(w, "An incorrect value is specified for the metric type", http.StatusBadRequest)
		return
	}

	metricVal := server.DefineMetricValue(metricType, metricValue)
	if metricVal == nil {
		http.Error(w, "Incorrect type of value for type metric", http.StatusBadRequest)
		return
	}

	if err := ms.Add(metricType, metricName, metricVal); err != nil {
		http.Error(w, "Incorrect type of value for type metric", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/plain")

	w.WriteHeader(http.StatusOK)

}
