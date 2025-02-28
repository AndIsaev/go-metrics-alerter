package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPublicKey(t *testing.T) {
	tempFile, err := os.CreateTemp("", "key.pem")
	if err != nil {
		t.Fatalf("Не удалось создать временный файл: %v", err)
	}
	defer os.Remove(tempFile.Name())
	cfg := &Config{PublicKeyPath: tempFile.Name()}

	_, err = tempFile.Write(
		[]byte(
			"-----BEGIN AGENT PUBLIC KEY-----\n" +
				"MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAx21q1jsVbHBUyvNoyTHj\n" +
				"ayuJgtANaCqSAs3PYaI4fM3ewaYXfb9osqvb0KS4drfSNij1M5Gf9W603e4ScNmy\n" +
				"UOEdZXuTUpSkuhxxbjgAxR4xShFS6N1ctpgRq1ZUWV6DWhGoOLH6+tnxEVWJp2ag\n" +
				"AByInipjThe5ZU2UQ6Zj9n7g3XpZZvZWJedFQDokr7GoWilUSgBwwFbLk7bx5C30\n" +
				"cxV0WUHPMnA5imjm4Bh17UFBp9rnk405JRbeMDbPeHp4GlSvCrqphBomyX1jq7Yf\n" +
				"zbQmQVde5AtL5+aAqtKMY/5sLvptMdi8JaUS5xqXVts65RVXtCVGzgCvteNR8Vqy\n" +
				"IQIDAQAB\n" +
				"-----END AGENT PUBLIC KEY-----",
		),
	)
	assert.NoError(t, err)

	err = tempFile.Close()
	assert.NoError(t, err)

	key, err := cfg.getPublicKey()
	assert.NoError(t, err)
	assert.NotNil(t, key)
}

func TestGetPublicKeyErrorFormat(t *testing.T) {
	tempFile, err := os.CreateTemp("", "public.lol")
	if err != nil {
		t.Fatalf("Не удалось создать временный файл: %v", err)
	}
	defer os.Remove(tempFile.Name())
	cfg := &Config{PublicKeyPath: tempFile.Name()}

	assert.NoError(t, err)

	err = tempFile.Close()
	assert.NoError(t, err)

	key, err := cfg.getPublicKey()
	assert.Errorf(t, err, "failed to decode PEM block containing public key")
	assert.Nil(t, key)
}

func TestGetPublicKeyDoesntExistsFile(t *testing.T) {
	cfg := &Config{PublicKeyPath: "key.pem"}

	key, err := cfg.getPublicKey()
	assert.Errorf(t, err, "error reading public key file")
	assert.Nil(t, key)
}
