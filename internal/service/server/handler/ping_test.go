package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

func ExampleHandler_PingHandler() {
	h := &Handler{
		MetricService: &MockMetricService{},
	}

	req, _ := http.NewRequest("GET", "/ping", nil)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler := h.PingHandler()
	handler.ServeHTTP(rr, req)

	fmt.Println("Status Code:", rr.Code)

	// Output:
	// Status Code: 200
}
