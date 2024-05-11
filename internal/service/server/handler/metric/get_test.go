package metric

import (
	"fmt"
	"github.com/AndIsaev/go-metrics-alerter/internal/storage"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetMetricHandler(t *testing.T) {
	r := chi.NewRouter()
	r.Mount(`/value/`, GetMetricRouter())

	ts := httptest.NewServer(r)
	defer ts.Close()

	type want struct {
		code        int
		response    string
		contentType string
		address     string
		key         storage.MetricKey
		value       interface{}
		method      string
	}

	tests := []struct {
		name string
		want want
	}{
		{
			name: "test #1 - if key not found",
			want: want{
				code:     http.StatusNotFound,
				response: fmt.Sprintf("%v\n", storage.ErrKeyErrorStorage.Error()),
				address:  "/value/gauge/Alloc",
				key:      "gauge/Alloc",
				method:   http.MethodGet,
			},
		},
		{
			name: "test #2 - if incorrect metric type",
			want: want{
				code:     http.StatusBadRequest,
				response: "An incorrect value is specified for the metric type\n",
				address:  "/value/error/Alloc",
				key:      "gauge/Alloc",
				method:   http.MethodGet,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ClearStorage()

			resp, body := testRequest(t, ts, tt.want.method, tt.want.address)
			defer resp.Body.Close()

			assert.Equal(t, tt.want.code, resp.StatusCode)
			assert.Equal(t, tt.want.response, body)
			assert.Nil(t, storage.MS.Metrics[tt.want.key])
		})
	}
}
