package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

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
	h := &Handler{
		MetricService: &MockMetricService{},
	}

	metric := common.Metrics{ID: "metric1", MType: "counter", Delta: linkInt64(123)}

	requestBody, _ := json.Marshal(metric)

	req, _ := http.NewRequest("POST", "/update", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler := h.UpdateRowHandler()
	handler.ServeHTTP(rr, req)

	fmt.Println("Status Code:", rr.Code)
	fmt.Println("Response Body:", rr.Body.String())

	// Output:
	// Status Code: 200
	// Response Body: {"id":"metric1","type":"counter","delta":123}
}

func ExampleHandler_SetMetricHandler() {
	h := &Handler{
		MetricService: &MockMetricService{},
	}

	req := httptest.NewRequest("POST", "/update/gauge/temperature/23.5", nil)
	rr := httptest.NewRecorder()

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("MetricType", "gauge")
	chiCtx.URLParams.Add("MetricName", "metric1")
	chiCtx.URLParams.Add("MetricValue", "23.5")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

	handler := h.SetMetricHandler()
	handler.ServeHTTP(rr, req)

	fmt.Println("Status Code:", rr.Code)

	// Output:
	// Status Code: 200
}

func TestSetMetricHandler(t *testing.T) {
	testMock := setupTest(t)
	handler := &Handler{MetricService: testMock.mockService}

	tests := []struct {
		name            string
		expectedErrFunc error
		expectedStatus  int
		result          string
		mType           string
		mName           string
		callFunc        bool
		mVal            string
	}{
		{
			name:            "valid metric counter",
			expectedErrFunc: nil,
			expectedStatus:  http.StatusOK,
			mType:           common.Counter,
			mName:           "metric1",
			mVal:            "10",
			callFunc:        true,
		},
		{
			name:            "error update metric",
			expectedErrFunc: errors.New("error update metric"),
			expectedStatus:  http.StatusBadRequest,
			mType:           common.Gauge,
			mName:           "metric1",
			mVal:            "1.5",
			result:          "error update metric\n",
			callFunc:        true,
		},
		{
			name:            "incorrect type",
			expectedErrFunc: nil,
			expectedStatus:  http.StatusBadRequest,
			mType:           "incorrect",
			mName:           "metric1",
			mVal:            "1",
			result:          "an incorrect value is specified for the metric type\n",
			callFunc:        false,
		},
		{
			name:            "incorrect value for metric type",
			expectedErrFunc: nil,
			expectedStatus:  http.StatusBadRequest,
			mType:           common.Counter,
			mName:           "metric1",
			mVal:            "19.9",
			result:          "incorrect value for counter type\n",
			callFunc:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/update/{MetricType}/{MetricName}/{MetricValue}", nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("MetricType", tt.mType)
			rctx.URLParams.Add("MetricName", tt.mName)
			rctx.URLParams.Add("MetricValue", tt.mVal)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			if tt.callFunc {
				metric := common.Metrics{ID: tt.mName, MType: tt.mType}
				MetricValue, _ := DefineMetricValue(tt.mType, tt.mVal)
				testMock.mockService.EXPECT().
					UpdateMetricByValue(req.Context(), metric, MetricValue).
					Return(tt.expectedErrFunc)
			}

			handler.SetMetricHandler().ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			assert.Equal(t, tt.result, rec.Body.String())
		})
	}
}

func TestUpdateRowHandler(t *testing.T) {
	testMock := setupTest(t)
	handler := &Handler{MetricService: testMock.mockService}

	tests := []struct {
		name               string
		expectedStatus     int
		result             string
		expectedResultFunc common.Metrics
		callFunc           bool
		body               string
		setup              func(ts *testSuite)
	}{
		{
			name:           "success update",
			expectedStatus: http.StatusOK,
			result:         `{"id":"metric1","type":"counter","delta":10}`,
			callFunc:       true,
			body:           `{"id":"metric1","type":"counter","delta":10}`,
			setup: func(ts *testSuite) {
				testMock.mockService.EXPECT().
					InsertMetric(context.Background(), common.Metrics{ID: "metric1", MType: common.Counter, Delta: linkInt64(10)}).
					Return(common.Metrics{ID: "metric1", MType: common.Counter, Delta: linkInt64(10)}, nil)
			},
		},
		{
			name:           "error update metric",
			expectedStatus: http.StatusInternalServerError,
			result:         "",
			callFunc:       true,
			body:           `{"id":"metric1","type":"counter","delta":10}`,
			setup: func(ts *testSuite) {
				testMock.mockService.EXPECT().
					InsertMetric(context.Background(), common.Metrics{ID: "metric1", MType: common.Counter, Delta: linkInt64(10)}).
					Return(common.Metrics{}, errors.New("error update metric"))
			},
		},
		{
			name:           "incorrect type",
			expectedStatus: http.StatusBadRequest,
			result:         "",
			callFunc:       false,
			body:           `{"id": "metric1", "type": "test"}`,
		},
		{
			name:           "empty body",
			expectedStatus: http.StatusInternalServerError,
			result:         "EOF\n",
			callFunc:       false,
			body:           "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/update", bytes.NewBufferString(tt.body))
			rec := httptest.NewRecorder()

			if tt.callFunc {
				tt.setup(testMock)
			}

			handler.UpdateRowHandler().ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			assert.Equal(t, tt.result, rec.Body.String())
		})
	}
}

func TestUpdateBatchHandler(t *testing.T) {
	testMock := setupTest(t)
	handler := &Handler{MetricService: testMock.mockService}

	tests := []struct {
		name               string
		expectedStatus     int
		result             string
		expectedResultFunc []common.Metrics
		callFunc           bool
		body               string
		setup              func(ts *testSuite)
	}{
		{
			name:               "success update",
			expectedStatus:     http.StatusOK,
			result:             `{"message":"success"}`,
			callFunc:           true,
			expectedResultFunc: []common.Metrics{{ID: "metric1", MType: common.Counter, Delta: linkInt64(10)}},
			body:               `[{"id":"metric1","type":"counter","delta":10}]`,
			setup: func(ts *testSuite) {
				testMock.mockService.EXPECT().
					InsertMetrics(context.Background(), []common.Metrics{{ID: "metric1", MType: common.Counter, Delta: linkInt64(10)}}).
					Return(nil)
			},
		},
		{
			name:               "error update metric",
			expectedStatus:     http.StatusInternalServerError,
			result:             "",
			callFunc:           true,
			body:               `[{"id":"metric1","type":"counter","delta":10}]`,
			expectedResultFunc: []common.Metrics{{ID: "metric1", MType: common.Counter, Delta: linkInt64(10)}},
			setup: func(ts *testSuite) {
				testMock.mockService.EXPECT().
					InsertMetrics(context.Background(), []common.Metrics{{ID: "metric1", MType: common.Counter, Delta: linkInt64(10)}}).
					Return(errors.New("error update metric"))
			},
		},

		{
			name:           "empty body",
			expectedStatus: http.StatusBadRequest,
			result:         "{\"message\":\"EOF\"}",
			callFunc:       false,
			body:           "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/updates", bytes.NewBufferString(tt.body))
			rec := httptest.NewRecorder()

			if tt.callFunc {
				tt.setup(testMock)
			}

			handler.UpdateBatchHandler().ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			assert.Equal(t, tt.result, rec.Body.String())
		})
	}
}
