package handlers

import (
	"github.com/AndIsaev/go-metrics-alerter/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateMetricHandler(t *testing.T) {
	var mockStorage = storage.MemStorage{
		Metrics: make(map[storage.MetricKey]interface{}),
	}

	type want struct {
		code        int
		response    string
		contentType string
		address     string
		key         storage.MetricKey
		value       interface{}
	}

	tests := []struct {
		name string
		ms   *storage.MemStorage
		want want
	}{
		{
			name: "success test #1",
			ms:   &mockStorage,
			want: want{code: 200, response: ``, contentType: "text/plain", address: "/update/gauge/Alloc/20.4", key: "gauge/Alloc", value: 20.4},
		},
		{
			name: "success test #2",
			ms:   &mockStorage,
			want: want{code: 200, response: ``, contentType: "text/plain", address: "/update/counter/pollCount/1", key: "counter/pollCount", value: int64(1)},
		},
		{
			name: "success test #3",
			ms:   &mockStorage,
			want: want{code: 200, response: ``, contentType: "text/plain", address: "/update/counter/pollCount/1", key: "counter/pollCount", value: int64(1)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.want.address, nil)

			// создаём новый Recorder
			w := httptest.NewRecorder()
			UpdateMetricHandler(tt.ms, w, request)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, tt.want.code, res.StatusCode)
			// получаем и проверяем тело запроса
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {

				}
			}(res.Body)
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, tt.want.response, string(resBody))
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
			assert.NotNil(t, mockStorage.Metrics[tt.want.key])

			switch tt.name {
			case "success test #1":
				assert.Equal(t, mockStorage.Metrics[tt.want.key], tt.want.value.(float64))
			case "success test #2":
				assert.Equal(t, mockStorage.Metrics[tt.want.key], tt.want.value.(int64))
			case "success test #3":
				assert.Equal(t, mockStorage.Metrics[tt.want.key], tt.want.value.(int64)*2)
			}

		})
	}
}

func TestUpdateMetricHandlerError(t *testing.T) {
	var mockStorage = storage.MemStorage{
		Metrics: make(map[storage.MetricKey]interface{}),
	}

	type want struct {
		code        int
		response    string
		contentType string
		address     string
		key         storage.MetricKey
		value       interface{}
	}

	tests := []struct {
		name string
		ms   *storage.MemStorage
		want want
	}{
		{
			name: "unsuccessful test #1",
			ms:   &mockStorage,
			want: want{code: http.StatusBadRequest, address: "/update/gauge/Alloc/error", key: "gauge/Alloc", value: "error"},
		},
		{
			name: "unsuccessful test #2",
			ms:   &mockStorage,
			want: want{code: http.StatusNotFound, address: "/update/counter/pollCount", key: "counter/pollCount"},
		},
		{
			name: "unsuccessful test #3",
			ms:   &mockStorage,
			want: want{code: http.StatusBadRequest, address: "/update/error/pollCount/1", key: "error/pollCount"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.want.address, nil)

			// создаём новый Recorder
			w := httptest.NewRecorder()
			UpdateMetricHandler(tt.ms, w, request)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, tt.want.code, res.StatusCode)
			// получаем и проверяем тело запроса
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {

				}
			}(res.Body)
			_, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Error(t, storage.ErrIncorrectMetricValue)
			assert.Nil(t, mockStorage.Metrics[tt.want.key])

			switch tt.name {
			case "unsuccessful test #1":
				assert.Equal(t, tt.want.code, http.StatusBadRequest)
			case "unsuccessful test #2":
				assert.Equal(t, tt.want.code, http.StatusNotFound)
			case "unsuccessful test #3":
				assert.Equal(t, tt.want.code, http.StatusBadRequest)
			}

		})
	}
}
