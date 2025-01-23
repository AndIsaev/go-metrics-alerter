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
	"github.com/AndIsaev/go-metrics-alerter/internal/service/agent/metrics"

	"github.com/stretchr/testify/require"
)

func linkInt64(num int64) *int64 {
	return &num
}

func TestInitHTTPClient(t *testing.T) {
	app := New()
	client := app.initHTTPClient()

	require.NotNil(t, client)
	require.Equal(t, time.Second*5, client.GetClient().Timeout)
}

func TestSendMetrics(t *testing.T) {
	defer httpmock.DeactivateAndReset()

	app := &AgentApp{
		Config: &Config{},
		Client: resty.New(),
	}
	httpmock.ActivateNonDefault(app.Client.GetClient())

	metric := metrics.StorageMetric{ID: "metric1", MType: common.Gauge, Delta: linkInt64(1)}
	storageMetrics := metrics.StorageMetrics{
		Metrics: make(map[string]metrics.StorageMetric),
	}
	storageMetrics.Metrics[metric.ID] = metric

	t.Run("success", func(t *testing.T) {
		httpmock.RegisterResponder("POST", app.Config.UpdateMetricsAddress,
			httpmock.NewBytesResponder(http.StatusOK, nil))

		err := app.sendMetrics(storageMetrics)
		assert.NoError(t, err)
	})

	t.Run("bad request", func(t *testing.T) {
		httpmock.RegisterResponder("POST", app.Config.UpdateMetricsAddress,
			httpmock.NewStringResponder(http.StatusBadRequest, `Bad Request`))

		err := app.sendMetrics(storageMetrics)
		assert.Error(t, err)
		assert.Equal(t, "unexpected status code: 400", err.Error())
	})

	t.Run("network error", func(t *testing.T) {
		httpmock.RegisterResponder("POST", app.Config.UpdateMetricsAddress,
			func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("some network error")
			})

		err := app.sendMetrics(storageMetrics)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "some network error")
	})
}
