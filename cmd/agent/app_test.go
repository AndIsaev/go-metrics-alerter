package main

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"

	"github.com/AndIsaev/go-metrics-alerter/internal/common"

	"github.com/stretchr/testify/require"
)

func linkInt64(num int64) *int64 {
	return &num
}

// MockIPResolver — мок-реализация интерфейса IPResolver для тестов
type MockIPResolver struct {
	IP  string
	Err error
}

// GetLocalIP возвращает заранее определенные значения
func (mock *MockIPResolver) GetLocalIP(address string) (string, error) {
	return mock.IP, mock.Err
}

func TestInitHTTPClient(t *testing.T) {
	app := New()
	client := app.initHTTPClient()

	require.NotNil(t, client)
	require.Equal(t, time.Second*5, client.GetClient().Timeout)
}

func TestSendMetrics(t *testing.T) {
	defer httpmock.DeactivateAndReset()
	mockResolver := &MockIPResolver{IP: "192.168.0.1", Err: nil}

	app := &AgentApp{
		Config:     &Config{},
		Client:     resty.New(),
		IPResolver: mockResolver,
	}
	httpmock.ActivateNonDefault(app.Client.GetClient())

	metrics := []common.Metrics{{ID: "metric1", MType: common.Gauge, Delta: linkInt64(1)}}

	t.Run("success", func(t *testing.T) {
		httpmock.RegisterResponder("POST", app.Config.UpdateMetricsAddress,
			httpmock.NewBytesResponder(http.StatusOK, nil))

		err := app.sendMetrics(metrics)
		assert.NoError(t, err)
	})

	t.Run("bad request", func(t *testing.T) {
		httpmock.RegisterResponder("POST", app.Config.UpdateMetricsAddress,
			httpmock.NewStringResponder(http.StatusBadRequest, `Bad Request`))

		err := app.sendMetrics(metrics)
		assert.Error(t, err)
	})

	t.Run("network error", func(t *testing.T) {
		httpmock.RegisterResponder("POST", app.Config.UpdateMetricsAddress,
			func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("some network error")
			})

		err := app.sendMetrics(metrics)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "some network error")
	})
}
