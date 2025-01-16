package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/AndIsaev/go-metrics-alerter/internal/common"

	"github.com/go-chi/chi"
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
	// Response Body: 23.5
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
	// Response Body: {"id":"metric1","type":"counter","value":23.5}
}
