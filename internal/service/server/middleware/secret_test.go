package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AndIsaev/go-metrics-alerter/internal/common"
)

func TestSecretMiddleware(t *testing.T) {
	middleware := SecretMiddleware("testkey")

	t.Run("Successful hash validation", func(t *testing.T) {
		body := "test message"
		expectedHash := common.Sha256sum([]byte(body), "testkey")

		req := httptest.NewRequest("POST", "/", bytes.NewBuffer([]byte(body)))
		req.Header.Set("HashSHA256", expectedHash)

		w := httptest.NewRecorder()
		handlerToTest := middleware(http.HandlerFunc(echoHandler))

		handlerToTest.ServeHTTP(w, req)

		resp := w.Result()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200 OK, got %v", resp.StatusCode)
		}

		if w.Header().Get("HashSHA256") != expectedHash {
			t.Errorf("Expected header HashSHA256 to be %s, got %s", expectedHash, w.Header().Get("HashSHA256"))
		}
	})

	t.Run("Invalid hash", func(t *testing.T) {
		body := "test message"
		invalidHash := "invalid_sha256_hash"

		req := httptest.NewRequest("POST", "/", bytes.NewBuffer([]byte(body)))
		req.Header.Set("HashSHA256", invalidHash)

		w := httptest.NewRecorder()
		handlerToTest := middleware(http.HandlerFunc(echoHandler))

		handlerToTest.ServeHTTP(w, req)

		resp := w.Result()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400 Bad Request, got %v", resp.StatusCode)
		}
	})

	t.Run("Empty key bypass", func(t *testing.T) {
		middlewareNoKey := SecretMiddleware("")

		body := "test message"
		req := httptest.NewRequest("POST", "/", bytes.NewBuffer([]byte(body)))

		w := httptest.NewRecorder()
		handlerToTest := middlewareNoKey(http.HandlerFunc(echoHandler))

		handlerToTest.ServeHTTP(w, req)

		resp := w.Result()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200 OK, got %v", resp.StatusCode)
		}
	})
}
