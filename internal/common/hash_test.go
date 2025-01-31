package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSha256sum(t *testing.T) {
	tests := []struct {
		name     string
		value    []byte
		key      string
		expected string
	}{
		{
			name:     "empty value and key",
			value:    []byte{},
			key:      "",
			expected: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855", // SHA-256 от пустой строки
		},
		{
			name:     "non-empty value and empty key",
			value:    []byte("hello"),
			key:      "",
			expected: "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824", // SHA-256 от "hello"
		},
		{
			name:     "non-empty value and key",
			value:    []byte("hello"),
			key:      "world",
			expected: "936a185caaa266bb9cbe981e9e05cb78cd732b0b3280eb944412bb6f8f8f07af", // заранее вычисленное значение
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Sha256sum(tt.value, tt.key)
			assert.Equal(t, tt.expected, result, "the SHA-256 checksum should match the expected value")
		})
	}
}
