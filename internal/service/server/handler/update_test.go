package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi"

	"github.com/AndIsaev/go-metrics-alerter/internal/common"
)

func ExampleHandler_UpdateBatchHandler() {
	// Создание инстанса вашего хендлера
	h := &Handler{
		MetricService: &MockMetricService{},
	}

	// Массив метрик для отправки
	metrics := []common.Metrics{
		{ID: "metric1", MType: "gauge", Value: new(float64)},
		{ID: "metric2", MType: "counter", Delta: new(int64)},
	}
	*metrics[0].Value = 123.45
	*metrics[1].Delta = 678

	// Кодирование метрик в JSON
	requestBody, _ := json.Marshal(metrics)

	// Создание нового HTTP-запроса
	req, _ := http.NewRequest("POST", "/updates", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// Создание ResponseRecorder для получения HTTP-ответа
	rr := httptest.NewRecorder()

	// Вызов хендлера с созданным запросом
	handler := h.UpdateBatchHandler()
	handler.ServeHTTP(rr, req)

	fmt.Println("Status Code:", rr.Code)
	fmt.Println("Response Body:", rr.Body.String())

	// Output:
	// Status Code: 200
	// Response Body: {"message":"success"}
}

func ExampleHandler_UpdateRowHandler() {
	// Создание инстанса вашего хендлера
	h := &Handler{
		MetricService: &MockMetricService{},
	}

	// Массив метрик для отправки
	metric := common.Metrics{ID: "metric1", MType: "counter", Value: new(float64)}
	*metric.Value = 123.45

	// Кодирование метрик в JSON
	requestBody, _ := json.Marshal(metric)

	// Создание нового HTTP-запроса
	req, _ := http.NewRequest("POST", "/update", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// Создание ResponseRecorder для получения HTTP-ответа
	rr := httptest.NewRecorder()

	// Вызов хендлера с созданным запросом
	handler := h.UpdateRowHandler()
	handler.ServeHTTP(rr, req)

	fmt.Println("Status Code:", rr.Code)
	fmt.Println("Response Body:", rr.Body.String())

	// Output:
	// Status Code: 200
	// Response Body: {"id":"metric1","type":"counter","value":123.45}
}

func ExampleHandler_SetMetricHandler() {
	// Создаем новый хендлер с макетным MetricService
	h := &Handler{
		MetricService: &MockMetricService{},
	}

	// Создаем пример HTTP-запроса с заданными параметрами URL
	req := httptest.NewRequest("POST", "/update/gauge/temperature/23.5", nil)
	rr := httptest.NewRecorder()

	// Используем chi.Router параметры в запросе
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("MetricType", "gauge")
	chiCtx.URLParams.Add("MetricName", "metric1")
	chiCtx.URLParams.Add("MetricValue", "23.5")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

	// Вызываем хендлер с искусственным запросом.
	handler := h.SetMetricHandler()
	handler.ServeHTTP(rr, req)

	// Выводим результирующий HTTP статус-код
	fmt.Println("Status Code:", rr.Code)

	// Output:
	// Status Code: 200
}
