package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"

	"github.com/AndIsaev/go-metrics-alerter/internal/service/server/handler"
	"github.com/AndIsaev/go-metrics-alerter/internal/storage/inmemory"

	"github.com/stretchr/testify/require"
)

func TestInitStorage(t *testing.T) {
	app := New()
	storage := app.initStorage

	require.NotNil(t, storage)
}
func TestInitRouter(t *testing.T) {
	config := &Config{}

	app := &ServerApp{
		Router:  chi.NewRouter(),
		Config:  config,
		Handler: &handler.Handler{},
		Conn:    inmemory.NewMemStorage(nil, false),
	}

	app.initRouter()

	// Проверка существования маршрута
	testCases := []struct {
		method       string
		url          string
		expectedCode int
	}{
		{"GET", "/", http.StatusOK},
		{"POST", "/update/counter/usage/1", http.StatusOK},
		{"POST", "/update", http.StatusOK},
		{"GET", "/ping", http.StatusOK},
		{"POST", "/value", http.StatusOK},
		{"GET", "/value/counter/usage", http.StatusOK},
		{"GET", "/debug/pprof", http.StatusOK},
	}

	for _, tc := range testCases {
		t.Run(tc.method+" "+tc.url, func(t *testing.T) {
			_, err := http.NewRequest(tc.method, tc.url, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()

			assert.Equal(t, tc.expectedCode, rr.Code, "Response code mismatch for method: %s, url: %s", tc.method, tc.url)
		})
	}

	// Проверка на Method Not Allowed
	req, err := http.NewRequest("PUT", "/", nil)
	assert.NoError(t, err, "Error while creating new request")
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code, "Expected Method Not Allowed")
}
