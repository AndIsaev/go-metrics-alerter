package http

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"

	"github.com/AndIsaev/go-metrics-alerter/internal/common"
)

func TestClient_SendMetricHandler(t *testing.T) {
	httpClient := NewClient("http://example.com/metrics", nil, "", nil)

	httpmock.ActivateNonDefault(httpClient.Client.GetClient())
	defer httpmock.DeactivateAndReset()

	tests := []struct {
		name        string
		mockSetup   func()
		expectedErr error
	}{
		{
			name: "successful request",
			mockSetup: func() {
				httpmock.RegisterResponder("POST", "http://example.com/metrics",
					httpmock.NewStringResponder(http.StatusOK, `{}`))
			},
			expectedErr: nil,
		},
		{
			name: "unexpected status code",
			mockSetup: func() {
				httpmock.RegisterResponder("POST", "http://example.com/metrics",
					httpmock.NewStringResponder(http.StatusBadRequest, `{}`))
			},
			expectedErr: fmt.Errorf("unexpected status code: 400"),
		},
		{
			name: "network error",
			mockSetup: func() {
				httpmock.RegisterNoResponder(
					httpmock.ConnectionFailure)
			},
			expectedErr: fmt.Errorf("unexpected status code: 400"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			body := common.Metrics{}

			err := httpClient.SendMetric(body)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

type MockIPResolver struct {
	Address string
}

func (r *MockIPResolver) GetLocalIP() (string, error) {
	return "", nil
}
func TestClient_SendMetricsHandler(t *testing.T) {
	ctx := context.Background()
	ipResolver := &MockIPResolver{}
	ipResolver.Address = ":8000"
	httpClient := NewClient("http://example.com/metrics", ipResolver, "", nil)

	httpmock.ActivateNonDefault(httpClient.Client.GetClient())
	defer httpmock.DeactivateAndReset()

	tests := []struct {
		name        string
		mockSetup   func()
		expectedErr error
	}{
		{
			name: "successful request",
			mockSetup: func() {
				httpmock.RegisterResponder("POST", "http://example.com/metrics",
					httpmock.NewStringResponder(http.StatusOK, `{}`))
			},
			expectedErr: nil,
		},
		{
			name: "unexpected status code",
			mockSetup: func() {
				httpmock.RegisterResponder("POST", "http://example.com/metrics",
					httpmock.NewStringResponder(http.StatusBadRequest, "unexpected status code: 400"))
			},
			expectedErr: fmt.Errorf("error sending request: unexpected status code: 400"),
		},
		{
			name: "network error",
			mockSetup: func() {
				httpmock.RegisterNoResponder(
					httpmock.ConnectionFailure)
			},
			expectedErr: fmt.Errorf("error sending request: unexpected status code: 400"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			body := []common.Metrics{
				{
					ID:    "metric1",
					MType: common.Gauge,
				},
				{
					ID:    "metric2",
					MType: common.Counter,
				},
			}

			err := httpClient.SendMetrics(ctx, body)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
