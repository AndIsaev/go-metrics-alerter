package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGzipMiddleware(t *testing.T) {
	// Функция-обработчик, которая возвращает "Hello, world!"
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, world!"))
	})

	// Создаем тестовый сервер
	ts := httptest.NewServer(GzipMiddleware(handler))
	defer ts.Close()

	// Создаем новый клиент, который поддерживает gzip
	client := &http.Client{}

	req, _ := http.NewRequest(http.MethodGet, ts.URL, nil)
	req.Header.Set("Accept-Encoding", "gzip")

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}
	defer resp.Body.Close()

	// Проверяем, что статус ответа 200 и установлен заголовок
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "gzip", resp.Header.Get("Content-Encoding"))

	gzr, err := gzip.NewReader(resp.Body)
	if err != nil {
		t.Fatalf("Failed to create new gzip reader: %v", err)
	}
	defer gzr.Close()

	// Читаем тело ответа после декомпрессии
	body, err := io.ReadAll(gzr)
	if err != nil {
		t.Fatalf("Failed to read gzip body: %v", err)
	}

	assert.Equal(t, "Hello, world!", string(body))
}
