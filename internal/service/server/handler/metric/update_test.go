package metric

import (
	//"github.com/AndIsaev/go-metrics-alerter/internal/service/server/handler"
	"github.com/AndIsaev/go-metrics-alerter/internal/storage"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"net/http"

	//"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateMetricHandler(t *testing.T) {
	r := chi.NewRouter()
	r.Mount(`/update/`, UpdateMetricRouter())

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
			name: "success test #1",
			want: want{
				code:        http.StatusOK,
				response:    ``,
				contentType: "text/plain",
				address:     "/update/gauge/Alloc/20.4",
				key:         "gauge/Alloc", value: 20.4,
				method: http.MethodPost,
			},
		},
		{
			name: "success test #2",
			want: want{
				code:        http.StatusOK,
				response:    ``,
				contentType: "text/plain",
				address:     "/update/counter/pollCount/1",
				key:         "counter/pollCount",
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
				key:         "counter/pollCount",
				value:       int64(1),
				method:      http.MethodPost,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, body := testRequest(t, ts, tt.want.method, tt.want.address)
			defer resp.Body.Close()

			assert.Equal(t, tt.want.code, resp.StatusCode)
			assert.Equal(t, tt.want.response, body)
			assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))
			assert.NotNil(t, storage.MS.Metrics[tt.want.key])

			switch tt.name {
			case "success test #1":
				assert.Equal(t, storage.MS.Metrics[tt.want.key], tt.want.value)
			case "success test #2":
				assert.Equal(t, storage.MS.Metrics[tt.want.key], tt.want.value.(int64))
			case "success test #3":
				assert.Equal(t, storage.MS.Metrics[tt.want.key], tt.want.value.(int64)*2)
			}

		})
	}
}

func TestUpdateMetricHandlerError(t *testing.T) {
	// clear storage before tests
	ClearStorage()

	r := chi.NewRouter()
	r.Mount(`/update/`, UpdateMetricRouter())

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
		ms   *storage.MemStorage
		want want
	}{
		{
			name: "unsuccessful test #1 - incorrect value",
			want: want{
				code:     http.StatusBadRequest,
				response: "incorrect value for gauge type\n",
				address:  "/update/gauge/Alloc/error",
				key:      "gauge/Alloc",
				value:    "error",
				method:   http.MethodPost,
			},
		},
		{
			name: "unsuccessful test #2 - incorrect address",
			want: want{
				code:     http.StatusNotFound,
				address:  "/update/counter/pollCount",
				key:      "counter/pollCount",
				method:   http.MethodPost,
				response: "404 page not found\n",
			},
		},
		{
			name: "unsuccessful test #3 - incorrect metric type",
			want: want{
				code:     http.StatusBadRequest,
				address:  "/update/error/pollCount/1",
				key:      "error/pollCount",
				response: "An incorrect value is specified for the metric type\n",
				method:   http.MethodPost,
			},
		},
		{
			name: "unsuccessful test #4 - incorrect method put",
			want: want{
				code:     http.StatusMethodNotAllowed,
				address:  "/update/counter/pollCount/1",
				key:      "counter/pollCount",
				response: "",
				method:   http.MethodPut,
			},
		},
		{
			name: "unsuccessful test #5 - incorrect method patch",
			want: want{
				code:     http.StatusMethodNotAllowed,
				address:  "/update/counter/pollCount/1",
				key:      "counter/pollCount",
				response: "",
				method:   http.MethodPatch,
			},
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			resp, body := testRequest(t, ts, tt.want.method, tt.want.address)
			defer resp.Body.Close()

			// создаём новый Recorder
			assert.Equal(t, tt.want.code, resp.StatusCode)
			assert.Equal(t, tt.want.response, body)
			assert.Error(t, storage.ErrIncorrectMetricValue)
			assert.Nil(t, storage.MS.Metrics[tt.want.key])

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
