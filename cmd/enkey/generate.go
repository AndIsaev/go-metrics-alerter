package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

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

func main() {
	privateKey, publicKey, err := GenerateRSAKeys(2048)
	if err != nil {
		fmt.Println("Error generating RSA keys:", err)
		return
	}

	privateKeyPEM := EncodePrivateKeyToPEM(privateKey)
	publicKeyPEM, err := EncodePublicKeyToPEM(publicKey)
	if err != nil {
		fmt.Println("Error encoding public key:", err)
		return
	}

	if err = WriteToFile("../server/private.pem", privateKeyPEM); err != nil {
		fmt.Println("Error writing private key to file:", err)
		return
	}

	if err = WriteToFile("../agent/public.pem", publicKeyPEM); err != nil {
		fmt.Println("Error writing public key to file:", err)
		return
	}
}
