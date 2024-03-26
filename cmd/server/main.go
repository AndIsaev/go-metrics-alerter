package main

import (
	"fmt"
	"net/http"
	"strings"
)

type Middleware func(http.Handler) http.Handler

var gauge string = "gauge"
var counter string = "counter"

type MemStorage struct {
	metrics map[string]interface{}
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		metrics: make(map[string]interface{}),
	}
}

func (ms *MemStorage) Update(metricType, metricName string, metricValue interface{}) {
	fmt.Println(metricType, metricName, metricValue)
	key := metricType + "/" + metricName
	if val, ok := ms.metrics[key]; ok {
		switch metricType {
		case "gauge":
			ms.metrics[key] = metricValue
		case "counter":
			ms.metrics[key] = val.(int64) + metricValue.(int64)
		}
	} else {
		ms.metrics[key] = metricValue
	}
}

func handleUpdateMetric(ms *MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		fmt.Println(w, r)
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		// получаем параметры
		params := r.URL.Path[len("/update/"):]
		parts := strings.Split(params, "/")
		fmt.Println(parts)

		// чекаем корректность урла
		if len(parts) != 3 {
			http.Error(w, "Not correct format URL", http.StatusBadRequest)
			return
		}

		metricType := parts[0]
		metricName := parts[1]
		metricValue := parts[2]

		// чекаем корректность названия метрики
		if metricName == "" {
			http.Error(w, "Имя метрики не указано", http.StatusNotFound)
			return
		}

		ms.Update(metricType, metricName, metricValue)
		w.Header().Set("Content-Type", "text/plain")

		w.WriteHeader(http.StatusOK)
		fmt.Println(ms.metrics)

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
