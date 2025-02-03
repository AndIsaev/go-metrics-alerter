package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateRSAKeys(t *testing.T) {
	privateKey, publicKey, err := GenerateRSAKeys(2048)
	assert.NoError(t, err)
	assert.NotNil(t, privateKey)
	assert.NotNil(t, publicKey)
}

func TestEncodePrivateKeyToPEM(t *testing.T) {
	privateKey, _, err := GenerateRSAKeys(2048)
	assert.NoError(t, err)

	privateKeyPEM := EncodePrivateKeyToPEM(privateKey)
	assert.NotEmpty(t, privateKeyPEM)
}

func TestEncodePublicKeyToPEM(t *testing.T) {
	_, publicKey, err := GenerateRSAKeys(2048)
	assert.NoError(t, err)

	publicKeyPEM, err := EncodePublicKeyToPEM(publicKey)
	assert.NoError(t, err)
	assert.NotEmpty(t, publicKeyPEM)
}

func TestWriteToFile(t *testing.T) {
	data := []byte("test data")
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "testfile.txt")

	err := WriteToFile(testFile, data)
	assert.NoError(t, err)

	readData, err := os.ReadFile(testFile)
	assert.NoError(t, err)
	assert.Equal(t, data, readData)
}

func TestRunGenerate(t *testing.T) {
	// Создаем временную директорию для тестов
	tempDir := t.TempDir()

	// Определяем пути для временных файлов
	privateKeyPath := filepath.Join(tempDir, "private.pem")
	publicKeyPath := filepath.Join(tempDir, "public.pem")

	// Запускаем функцию с временными файлами
	err := RunGenerate(privateKeyPath, publicKeyPath)
	assert.NoError(t, err, "Функция RunGenerate должна завершиться без ошибок")

	// Проверяем, что оба файла были созданы и содержат данные
	privateFileInfo, err := os.Stat(privateKeyPath)
	assert.NoError(t, err, "Файл приватного ключа должен существовать")
	assert.NotZero(t, privateFileInfo.Size(), "Файл приватного ключа не должен быть пустым")

	publicFileInfo, err := os.Stat(publicKeyPath)
	assert.NoError(t, err, "Файл публичного ключа должен существовать")
	assert.NotZero(t, publicFileInfo.Size(), "Файл публичного ключа не должен быть пустым")
}
