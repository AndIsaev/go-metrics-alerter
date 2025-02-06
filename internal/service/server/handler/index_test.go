package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/AndIsaev/go-metrics-alerter/internal/common"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
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

func TestIndexHandler(t *testing.T) {
	testMock := setupTest(t)
	handler := &Handler{MetricService: testMock.mockService}

	tests := []struct {
		name           string
		expectedStatus int
		result         string
		callFunc       bool
		setup          func(ts *testSuite)
	}{
		{
			name:           "success index",
			expectedStatus: http.StatusOK,
			result:         `[{"id":"metric1","type":"counter","delta":10}]`,
			callFunc:       true,
			setup: func(ts *testSuite) {
				testMock.mockService.EXPECT().
					ListMetrics(context.Background()).
					Return([]common.Metrics{{ID: "metric1", MType: common.Counter, Delta: linkInt64(10)}}, nil)
			},
		},
		{
			name:           "error getting list metric",
			expectedStatus: http.StatusInternalServerError,
			result:         "",
			callFunc:       true,
			setup: func(ts *testSuite) {
				testMock.mockService.EXPECT().
					ListMetrics(context.Background()).
					Return([]common.Metrics{}, errors.New("error getting list metric"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()

			if tt.callFunc {
				tt.setup(testMock)
			}

			handler.IndexHandler().ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			assert.Equal(t, tt.result, rec.Body.String())
		})
	}
}
