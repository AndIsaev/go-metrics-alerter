package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
)

func TestEncrypt(t *testing.T) {
	// Генерация пар ключей для тестирования.
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}
	publicKey := &privateKey.PublicKey

	// Тестовые данные
	data := []byte("test data")

	// Позитивный тест - успешно шифруем данные
	ciphertext, err := Encrypt(publicKey, data)
	if err != nil {
		t.Errorf("Failed to encrypt data: %v", err)
	}

	// Декодируем данные обратно для проверки правильности.
	plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, ciphertext)
	if err != nil {
		t.Errorf("Failed to decrypt ciphertext: %v", err)
	}

	// Проверяем, что расшифрованный текст совпадает с оригинальными данными
	if string(plaintext) != string(data) {
		t.Errorf("Decrypted text does not match the original data. Got %s, want %s", plaintext, data)
	}
}
