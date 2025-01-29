package middleware

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDecryptMiddleware(t *testing.T) {
	// Генерация пары ключей для теста
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate RSA keys: %v", err)
	}
	publicKey := &privateKey.PublicKey

	// Исходные данные для теста
	plaintext := "Hello, World!"

	// Шифруем исходные данные для их передачи в middleware
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, []byte(plaintext))
	if err != nil {
		t.Fatalf("Failed to encrypt data: %v", err)
	}

	// Создаём middleware
	middleware := DecryptMiddleware(privateKey)

	// Используем httptest для создания тестового сервера
	nextHandler := http.HandlerFunc(echoHandler)
	handlerToTest := middleware(nextHandler)

	req := httptest.NewRequest("POST", "/", bytes.NewReader(ciphertext))
	w := httptest.NewRecorder()

	handlerToTest.ServeHTTP(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	// Проверяем статус и содержание ответа
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Unexpected status code: got %v want %v", resp.StatusCode, http.StatusOK)
	}

	if string(body) != plaintext {
		t.Errorf("Unexpected body content: got %s want %s", string(body), plaintext)
	}
}

func TestDecryptMiddleware_InvalidData(t *testing.T) {
	// Генерация пары ключей для теста
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate RSA keys: %v", err)
	}

	// Создаём middleware
	middleware := DecryptMiddleware(privateKey)

	// Используем httptest для создания тестового сервера
	nextHandler := http.HandlerFunc(echoHandler)
	handlerToTest := middleware(nextHandler)

	// Неправильные (незашифрованные) данные
	badData := "invalid encrypted data"

	req := httptest.NewRequest("POST", "/", strings.NewReader(badData))
	w := httptest.NewRecorder()

	handlerToTest.ServeHTTP(w, req)

	resp := w.Result()

	// Проверяем статус код
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Unexpected status code: got %v want %v", resp.StatusCode, http.StatusBadRequest)
	}
}
