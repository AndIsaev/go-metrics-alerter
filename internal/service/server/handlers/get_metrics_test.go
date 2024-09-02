package handlers

import (
	"fmt"
	"github.com/AndIsaev/go-metrics-alerter/internal/storage"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/AndIsaev/go-metrics-alerter/internal/storage/mocks"

	"github.com/stretchr/testify/assert"
)

func TestGetMetricHandler(t *testing.T) {
	testApp := NewTestServerApp()
	ctrl := gomock.NewController(t)
	testApp.DBConn = mock.NewMockBaseStorage(ctrl)

	testApp.MemStorage.Metrics["pollCount"] = 20

	ts := testApp.Server
	defer ts.Close()

	type want struct {
		code     int
		response string
		address  string
		key      string
		value    interface{}
		method   string
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
				key:      "Alloc",
				method:   http.MethodGet,
			},
		},
		{
			name: "test #2 - if incorrect metric type",
			want: want{
				code:     http.StatusBadRequest,
				response: "An incorrect value is specified for the metric type\n",
				address:  "/value/error/Alloc",
				key:      "Alloc",
				method:   http.MethodGet,
			},
		},
		{
			name: "test #3 - case with counter type",
			want: want{
				code:     http.StatusOK,
				response: "20",
				address:  "/value/counter/pollCount",
				key:      "pollCount",
				method:   http.MethodGet,
				value:    20,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, _ := testRequest(t, ts, tt.want.method, tt.want.address)
			resp.Body.Close()

			if tt.name != "test #3 - case with counter type" {
				assert.Nil(t, testApp.MemStorage.Metrics[tt.want.key])
			}

			assert.Equal(t, tt.want.code, resp.StatusCode)
			//assert.Equal(t, tt.want.response, body)
		})
	}
}
