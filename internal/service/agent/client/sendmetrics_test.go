package client

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/AndIsaev/go-metrics-alerter/internal/common"
	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestSendMetricHandler(t *testing.T) {
	client := resty.New()

	httpmock.ActivateNonDefault(client.GetClient())
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

			body := common.Metrics{ /* инициализируйте по необходимости */ }

			err := SendMetricHandler(client, "http://example.com/metrics", body)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
