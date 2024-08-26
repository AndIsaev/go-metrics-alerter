package handlers

import (
	"net/http"

	"github.com/golang/mock/gomock"

	"github.com/AndIsaev/go-metrics-alerter/internal/storage/mock"

	"github.com/stretchr/testify/assert"

	"github.com/AndIsaev/go-metrics-alerter/internal/storage"

	"testing"
)

func TestUpdateMetricHandler(t *testing.T) {
	testApp := NewTestServerApp()
	ctrl := gomock.NewController(t)
	testApp.DBConn = mock.NewMockPgStorage(ctrl)

	ts := testApp.Server

	type want struct {
		code        int
		response    string
		contentType string
		address     string
		key         string
		value       interface{}
		method      string
	}

	tests := []struct {
		name string
		want want
	}{
		{
			name: "success test #1",
			want: want{
				code:        http.StatusOK,
				response:    ``,
				contentType: "text/plain",
				address:     "/update/gauge/Alloc/20.4",
				key:         "Alloc",
				value:       20.4,
				method:      http.MethodPost,
			},
		},
		{
			name: "success test #2",
			want: want{
				code:        http.StatusOK,
				response:    ``,
				contentType: "text/plain",
				address:     "/update/counter/pollCount/1",
				key:         "pollCount",
				value:       int64(1),
				method:      http.MethodPost,
			},
		},
		{
			name: "success test #3",
			want: want{
				code:        http.StatusOK,
				response:    ``,
				contentType: "text/plain",
				address:     "/update/counter/pollCount/1",
				key:         "pollCount",
				value:       int64(1),
				method:      http.MethodPost,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, _ := testRequest(t, ts, tt.want.method, tt.want.address)
			resp.Body.Close()

			assert.Equal(t, tt.want.code, resp.StatusCode)
			//assert.Equal(t, tt.want.response, body)
			assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))
			assert.NotNil(t, testApp.MemStorage.Metrics[tt.want.key])

			switch tt.name {
			case "success test #1":
				assert.Equal(t, testApp.MemStorage.Metrics[tt.want.key], tt.want.value)
			case "success test #2":
				assert.Equal(t, testApp.MemStorage.Metrics[tt.want.key], tt.want.value.(int64))
			case "success test #3":
				assert.Equal(t, testApp.MemStorage.Metrics[tt.want.key], tt.want.value.(int64)*2)
			}
		})
	}
}

func TestUpdateMetricHandlerError(t *testing.T) {
	testApp := NewTestServerApp()
	ctrl := gomock.NewController(t)
	testApp.DBConn = mock.NewMockPgStorage(ctrl)

	ts := testApp.Server

	type want struct {
		code     int
		response string
		address  string
		key      string
		method   string
		json     bool
	}

	tests := []struct {
		name string
		ms   *storage.MemStorage
		want want
	}{
		{
			name: "unsuccessful test #1 - incorrect value",
			want: want{
				code:     http.StatusBadRequest,
				response: "incorrect value for gauge type\n",
				address:  "/update/gauge/Alloc/error",
				key:      "Alloc",
				method:   http.MethodPost,
			},
		},
		{
			name: "unsuccessful test #3 - incorrect metric type",
			want: want{
				code:     http.StatusBadRequest,
				address:  "/update/error/pollCount/1",
				key:      "pollCount",
				response: "An incorrect value is specified for the metric type\n",
				method:   http.MethodPost,
			},
		},
		{
			name: "unsuccessful test #4 - incorrect method put",
			want: want{
				code:     http.StatusMethodNotAllowed,
				address:  "/update/counter/pollCount/1",
				key:      "pollCount",
				response: "{\"message\":\"method is not valid\"}",
				method:   http.MethodPut,
				json:     true,
			},
		},
		{
			name: "unsuccessful test #5 - incorrect method patch",
			want: want{
				code:     http.StatusMethodNotAllowed,
				address:  "/update/counter/pollCount/1",
				key:      "pollCount",
				response: "{\"message\":\"method is not valid\"}",
				method:   http.MethodPatch,
				json:     true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, body := testRequest(t, ts, tt.want.method, tt.want.address)
			resp.Body.Close()

			// создаём новый Recorder
			assert.Equal(t, tt.want.code, resp.StatusCode)
			assert.Equal(t, tt.want.response, body)
			assert.Error(t, storage.ErrIncorrectMetricValue)
			assert.Nil(t, testApp.MemStorage.Metrics[tt.want.key])

			switch tt.name {
			case "unsuccessful test #1":
				assert.Equal(t, tt.want.code, http.StatusBadRequest)
			case "unsuccessful test #2":
				assert.Equal(t, tt.want.code, http.StatusNotFound)
			case "unsuccessful test #3":
				assert.Equal(t, tt.want.code, http.StatusBadRequest)
			case "unsuccessful test #4":
				assert.Equal(t, tt.want.code, http.StatusMethodNotAllowed)
			case "unsuccessful test #5":
				assert.Equal(t, tt.want.code, http.StatusMethodNotAllowed)
			}
		})
	}
}
