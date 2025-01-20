package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

func ExampleHandler_IndexHandler() {
	h := &Handler{
		MetricService: &MockMetricService{},
	}

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler := h.IndexHandler()
	handler.ServeHTTP(rr, req)

	fmt.Println("Status Code:", rr.Code)
	fmt.Println("Response Body:", rr.Body)

	// Output:
	// Status Code: 200
	// Response Body: [{"id":"metric1","type":"counter","delta":1},{"id":"metric2","type":"gauge","value":10.4}]
}
