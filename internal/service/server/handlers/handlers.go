package handlers

import (
	"github.com/AndIsaev/go-metrics-alerter/internal/service/server"
	"github.com/AndIsaev/go-metrics-alerter/internal/storage"

	"net/http"
	"strconv"
	"strings"
)

func UpdateMetricHandler(ms *storage.MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		// получаем параметры
		params := r.URL.Path[len("/update/"):]
		parts := strings.Split(params, "/")

		// чекаем корректность урла
		if len(parts) != 3 {
			http.Error(w, "Not correct format URL", http.StatusNotFound)
			return
		}

		metricType := parts[0]
		metricName := parts[1]
		metricValue := parts[2]

		// чекаем корректность названия метрики
		if !server.IsCorrectType(metricType) {
			http.Error(w, "Указано не корректное значение для типа метрики", http.StatusBadRequest)
			return
		}

		if metricName == "" {
			http.Error(w, "Имя метрики не указано", http.StatusNotFound)
			return
		}
		val, err := strconv.Atoi(metricValue)
		if err != nil {
			http.Error(w, "Указано не корректное значение", http.StatusBadRequest)
			return
		}
		metricValue = strconv.Itoa(val)

		if metricValue == "" {
			http.Error(w, "Не указано значение", http.StatusBadRequest)
			return
		}

		ms.Update(metricType, metricName, metricValue)
		w.Header().Set("Content-Type", "text/plain")

		w.WriteHeader(http.StatusOK)

	}
}
