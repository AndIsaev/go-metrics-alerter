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

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"

	"github.com/AndIsaev/go-metrics-alerter/internal/common"
)

func ExampleHandler_GetByURLParamHandler() {
	h := &Handler{
		MetricService: &MockMetricService{},
	}

	req := httptest.NewRequest("GET", "/value/gauge/metric1", nil)
	rr := httptest.NewRecorder()

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("MetricType", "counter")
	chiCtx.URLParams.Add("MetricName", "metric1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

	handler := h.GetByURLParamHandler()
	handler.ServeHTTP(rr, req)

	fmt.Println("Response Body:", rr.Body.String())

	// Output:
	// Response Body: 23
}

func ExampleHandler_GetHandler() {
	h := &Handler{
		MetricService: &MockMetricService{},
	}
	metric := common.Metrics{ID: "metric1", MType: "counter", Value: new(float64)}
	*metric.Value = 23.5

	requestBody, _ := json.Marshal(metric)

	req, _ := http.NewRequest("GET", "/value", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler := h.GetHandler()
	handler.ServeHTTP(rr, req)

	fmt.Println("Status Code:", rr.Code)
	fmt.Println("Response Body:", rr.Body.String())

	// Output:
	// Status Code: 200
	// Response Body: {"id":"metric1","type":"counter","delta":23}
}

func TestGetHandler(t *testing.T) {
	testMock := setupTest(t)
	handler := &Handler{MetricService: testMock.mockService}

	tests := []struct {
		name           string
		setup          func(ts *testSuite)
		serviceError   error
		expectedStatus int
		inputBody      string
	}{
		{
			name:           "Valid metric",
			expectedStatus: http.StatusOK,
			inputBody:      `{"id":"metric1","type":"gauge"}`,
			setup: func(ts *testSuite) {
				ts.mockService.EXPECT().
					GetMetricByName(ts.ctx, "metric1").
					Return(common.Metrics{ID: "metric1", MType: common.Gauge, Delta: linkInt64(25)}, nil)
			},
		},
		{
			name:           "incorrect json",
			expectedStatus: http.StatusBadRequest,
			inputBody:      `{"id":incorrect,"type":"gauge"`,
			setup:          nil,
		},
		{
			name:           "incorrect type metric",
			expectedStatus: http.StatusBadRequest,
			inputBody:      `{"id":"metric1","type":"incorrect"}`,
			setup:          nil,
		},
		{
			name:           "not found metric",
			expectedStatus: http.StatusNotFound,
			inputBody:      `{"id":"metric1","type":"gauge"}`,
			setup: func(ts *testSuite) {
				ts.mockService.EXPECT().
					GetMetricByName(ts.ctx, "metric1").
					Return(common.Metrics{}, errors.New("metric not found"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(testMock)
			}

			req := httptest.NewRequest(http.MethodPost, "/value", bytes.NewBuffer([]byte(tt.inputBody)))
			rec := httptest.NewRecorder()

			handler.GetHandler().ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}

func TestGetByURLParamHandler(t *testing.T) {
	testMock := setupTest(t)
	handler := &Handler{MetricService: testMock.mockService}

	tests := []struct {
		name               string
		expectedErrFunc    error
		expectedStatus     int
		result             string
		mType              string
		mName              string
		expectedResultFunc common.Metrics
		callFunc           bool
	}{
		{
			name:               "valid metric counter",
			expectedErrFunc:    nil,
			expectedStatus:     http.StatusOK,
			mType:              common.Counter,
			mName:              "metric1",
			result:             "10",
			expectedResultFunc: common.Metrics{ID: "metric1", MType: common.Counter, Delta: linkInt64(10)},
			callFunc:           true,
		},
		{
			name:               "valid metric gauge",
			expectedErrFunc:    nil,
			expectedStatus:     http.StatusOK,
			mType:              common.Gauge,
			mName:              "metric1",
			result:             "10.54",
			expectedResultFunc: common.Metrics{ID: "metric1", MType: common.Gauge, Value: linkFloat64(10.54)},
			callFunc:           true,
		},
		{
			name:               "incorrect type",
			expectedErrFunc:    nil,
			expectedStatus:     http.StatusBadRequest,
			mType:              "incorrect",
			mName:              "metric1",
			result:             "an incorrect value is specified for the metric type\n",
			expectedResultFunc: common.Metrics{},
			callFunc:           false,
		},
		{
			name:               "metric not found",
			expectedErrFunc:    errors.New("metric not found"),
			expectedStatus:     http.StatusNotFound,
			mType:              common.Gauge,
			mName:              "metric1",
			result:             "metric not found\n",
			expectedResultFunc: common.Metrics{},
			callFunc:           true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/value/{MetricType}/{MetricName}", nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("MetricType", tt.mType)
			rctx.URLParams.Add("MetricName", tt.mName)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			if tt.callFunc {
				testMock.mockService.EXPECT().
					GetMetricByName(req.Context(), tt.mName).
					Return(tt.expectedResultFunc, tt.expectedErrFunc)
			}

			handler.GetByURLParamHandler().ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			assert.Equal(t, tt.result, rec.Body.String())
		})
	}
}
