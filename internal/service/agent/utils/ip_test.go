package utils

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLocalIP(t *testing.T) {
	resolver := NewDefaultIPResolver()

	tests := []struct {
		name           string
		address        string
		expectedIP     string
		expectingError bool
	}{
		{
			name:           "Valid address with open port",
			address:        "google.com:80",
			expectedIP:     "",
			expectingError: false,
		},
		{
			name:           "Invalid address with non-existant hostname",
			address:        "invalid.local:80",
			expectedIP:     "",
			expectingError: true,
		},
		{
			name:           "Invalid address without specified port",
			address:        "localhost",
			expectedIP:     "",
			expectingError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ip, err := resolver.GetLocalIP(tt.address)

			if tt.expectingError {
				assert.Error(t, err)
				assert.Empty(t, ip)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, ip)
				parsedIP := net.ParseIP(ip)
				assert.NotNil(t, parsedIP)
			}
		})
	}
}
