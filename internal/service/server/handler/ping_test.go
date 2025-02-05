package handler

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestPingHandler(t *testing.T) {
	testMock := setupTest(t)
	handler := &Handler{MetricService: testMock.mockService}

	tests := []struct {
		name               string
		expectedStatus     int
		expectedResultFunc error
	}{
		{
			name:               "success ping",
			expectedStatus:     http.StatusOK,
			expectedResultFunc: nil,
		},
		{
			name:               "error ping",
			expectedStatus:     http.StatusInternalServerError,
			expectedResultFunc: errors.New("ping error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/ping", nil)

			testMock.mockService.EXPECT().
				PingStorage(req.Context()).
				Return(tt.expectedResultFunc)

			handler.PingHandler().ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}
