package client

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"

	"github.com/AndIsaev/go-metrics-alerter/internal/common"
)

func StartMockServer(t *testing.T, responses map[string][]byte) *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if response, ok := responses[r.URL.Path]; ok {
			w.WriteHeader(http.StatusOK)
			w.Write(response)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	t.Cleanup(func() { ts.Close() })
	return ts
}

func TestSendMetricsClient(t *testing.T) {
	responses := make(map[string][]byte)
	responses["/update/counter/pollCount/1"] = []byte{}
	c := resty.New()

	mockServer := StartMockServer(t, responses)
	defer mockServer.Close()

	type want struct {
		url         string
		contentType string
		body        common.Metrics
		status      int
		method      string
	}
	tests := []struct {
		name string
		want want
	}{
		{name: "success test #1", want: want{url: mockServer.URL + `/update/counter/pollCount/1`, contentType: "text/plain", body: common.Metrics{}, status: http.StatusOK, method: http.MethodPost}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := httptest.NewRecorder()
			err := SendMetricsClient(c, tt.want.url, tt.want.body)

			assert.NoError(t, err)
			assert.Equal(t, tt.want.status, resp.Code)
		})
	}
}
