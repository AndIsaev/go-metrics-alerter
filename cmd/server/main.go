package main

import (
	"net/http"
	"strconv"
	"strings"
)

type MemStorage struct {
	metrics map[string]interface{}
}

var gauge = "gauge"
var counter = "counter"

func NewMemStorage() *MemStorage {
	return &MemStorage{
		metrics: make(map[string]interface{}),
	}
}

func IsCorrectType(metricType string) bool {
	for _, v := range []string{counter, gauge} {
		if v == metricType {
			return true
		}
	}
	return false

}

func (ms *MemStorage) Update(metricType, metricName string, metricValue interface{}) {
	key := metricType + "/" + metricName
	if val, ok := ms.metrics[key]; ok {
		switch metricType {
		case "gauge":
			ms.metrics[key] = metricValue
		case "counter":
			ms.metrics[key] = val.(int64) + metricValue.(int64)
		}
	}
}

func handleUpdateMetric(ms *MemStorage) http.HandlerFunc {
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
		if !IsCorrectType(metricType) {
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

func main() {
	mux := http.NewServeMux()
	ms := NewMemStorage()

	mux.HandleFunc(`/update/`, handleUpdateMetric(ms))

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
