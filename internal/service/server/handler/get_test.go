package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/AndIsaev/go-metrics-alerter/internal/common"
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi"
)

func ExampleHandler_GetByURLParamHandler() {
	// Создаем новый хендлер с макетным MetricService
	h := &Handler{
		MetricService: &MockMetricService{},
	}

	// Создаем пример HTTP-запроса с заданными параметрами URL
	req := httptest.NewRequest("GET", "/value/gauge/metric1", nil)
	rr := httptest.NewRecorder()

	// Используем chi.RouteContext параметры в запросе
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("MetricType", "counter")
	chiCtx.URLParams.Add("MetricName", "metric1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

	// Вызываем хендлер с искусственным запросом.
	handler := h.GetByURLParamHandler()
	handler.ServeHTTP(rr, req)

	// Выводим результат (ожидаемое значение "23.5")
	fmt.Println("Response Body:", rr.Body.String())

	// Output:
	// Response Body: 23.5
}

func ExampleHandler_GetHandler() {
	// Создаем новый хендлер с макетным MetricService
	h := &Handler{
		MetricService: &MockMetricService{},
	}
	metric := common.Metrics{ID: "metric1", MType: "counter", Value: new(float64)}
	*metric.Value = 23.5

	// Кодирование метрик в JSON
	requestBody, _ := json.Marshal(metric)

	req, _ := http.NewRequest("GET", "/value", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// Создание ResponseRecorder для получения HTTP-ответа
	rr := httptest.NewRecorder()

	// Вызов хендлера с созданным запросом
	handler := h.GetHandler()
	handler.ServeHTTP(rr, req)

	fmt.Println("Status Code:", rr.Code)
	fmt.Println("Response Body:", rr.Body.String())

	// Output:
	// Status Code: 200
	// Response Body: {"id":"metric1","type":"counter","value":23.5}
}
