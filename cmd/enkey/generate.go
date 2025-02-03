package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"
)

const privatePath = "../server/private.pem"
const publicPath = "../agent/public.pem"

// GenerateRSAKeys генерирует пару RSA ключей
func GenerateRSAKeys(bits int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}
	return privateKey, &privateKey.PublicKey, nil
}

// EncodePrivateKeyToPEM кодирует приватный ключ в формате PEM
func EncodePrivateKeyToPEM(privateKey *rsa.PrivateKey) []byte {
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	return pem.EncodeToMemory(&pem.Block{
		Type:  "SERVER PRIVATE KEY",
		Bytes: privateKeyBytes,
	})
}

// EncodePublicKeyToPEM кодирует публичный ключ в формате PEM
func EncodePublicKeyToPEM(publicKey *rsa.PublicKey) ([]byte, error) {
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, err
	}
	return pem.EncodeToMemory(&pem.Block{
		Type:  "AGENT PUBLIC KEY",
		Bytes: publicKeyBytes,
	}), nil
}

// WriteToFile записывает данные в файл
func WriteToFile(filename string, data []byte) error {
	return os.WriteFile(filename, data, 0644)
}

func RunGenerate(privatePath, publicPath string) error {
	privateKey, publicKey, err := GenerateRSAKeys(2048)
	if err != nil {
		return err
	}

	privateKeyPEM := EncodePrivateKeyToPEM(privateKey)
	publicKeyPEM, err := EncodePublicKeyToPEM(publicKey)
	if err != nil {
		return err
	}

	if err = WriteToFile(privatePath, privateKeyPEM); err != nil {
		return err
	}

	if err = WriteToFile(publicPath, publicKeyPEM); err != nil {
		return err
	}
	return nil
}

func main() {
	err := RunGenerate(privatePath, publicPath)
	if err != nil {
		log.Fatalf("error generate keys: %v", err)
	}
}
