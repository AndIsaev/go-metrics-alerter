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
