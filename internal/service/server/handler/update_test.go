package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/AndIsaev/go-metrics-alerter/internal/common"
)

// MockMetricService - это макетный сервис для тестирования
type MockMetricService struct{}

func linkFloat64(num float64) *float64 {
	return &num
}
func (m *MockMetricService) InsertMetrics(_ context.Context, _ []common.Metrics) error {
	return nil
}

func (m *MockMetricService) PingStorage(_ context.Context) error {
	return nil
}
func (m *MockMetricService) CloseStorage(_ context.Context) error {
	return nil
}
func (m *MockMetricService) RunMigrationsStorage(_ context.Context) error {
	return nil
}

func (m *MockMetricService) ListMetrics(_ context.Context) ([]common.Metrics, error) {
	return []common.Metrics{}, nil
}

func (m *MockMetricService) UpdateMetricByValue(_ context.Context, _ common.Metrics, _ any) error {
	return nil
}

func (m *MockMetricService) GetMetricByName(_ context.Context, _ string) (common.Metrics, error) {
	return common.Metrics{}, nil
}

func (m *MockMetricService) GetMetricByNameType(_ context.Context, _ string, _ string) (common.Metrics, error) {
	return common.Metrics{}, nil
}

func (m *MockMetricService) InsertMetric(_ context.Context, _ common.Metrics) (common.Metrics, error) {
	return common.Metrics{
		ID:    "metric1",
		MType: "counter",
		Value: linkFloat64(123.45),
	}, nil
}

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
