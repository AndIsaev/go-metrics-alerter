package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrustedSubnetMiddleware(t *testing.T) {
	// Функция для тестового обработчика, который находится за мидлвэйером.
	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Success"))
	})

	tests := []struct {
		name         string
		trustIP      string
		realIP       string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Trusted IP",
			trustIP:      "192.168.1.0/24",
			realIP:       "192.168.1.10",
			expectedCode: http.StatusOK,
			expectedBody: "Success",
		},
		{
			name:         "Untrusted IP",
			trustIP:      "192.168.1.0/24",
			realIP:       "192.168.2.10",
			expectedCode: http.StatusForbidden,
			expectedBody: "Forbidden\n",
		},
		{
			name:         "Missing X-Real-IP header",
			trustIP:      "192.168.1.0/24",
			realIP:       "",
			expectedCode: http.StatusForbidden,
			expectedBody: "Forbidden\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "http://example.com", nil)
			assert.NoError(t, err)

			// Установим заголовок X-Real-IP
			if tt.realIP != "" {
				req.Header.Set("X-Real-IP", tt.realIP)
			}

			rr := httptest.NewRecorder()

			// Создаём мидлвэйер с заданными параметрами и передаём финальный обработчик
			handler := TrustedSubnetMiddleware(tt.trustIP)(finalHandler)

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedCode, rr.Code)
			assert.Equal(t, tt.expectedBody, rr.Body.String())
		})
	}
}

func TestIsIPInSubnet(t *testing.T) {
	tests := []struct {
		ipStr string
		cidr  string
		want  bool
		name  string
	}{
		{
			ipStr: "192.168.1.1",
			cidr:  "192.168.1.0/24",
			want:  true,
			name:  "IP is in the subnet",
		},
		{
			ipStr: "192.168.2.1",
			cidr:  "192.168.1.0/24",
			want:  false,
			name:  "IP is not in the subnet",
		},
		{
			ipStr: "invalid-ip",
			cidr:  "192.168.1.0/24",
			want:  false,
			name:  "Invalid IP address",
		},
		{
			ipStr: "192.168.1.1",
			cidr:  "invalid-cidr",
			want:  false,
			name:  "Invalid CIDR",
		},
		{
			ipStr: "2607:f8b0:4005:0800:0000:0000:0000:200e", // Example of an IPv6 address
			cidr:  "2607:f8b0:4005:0800::/64",                // Corresponding IPv6 subnet
			want:  true,
			name:  "IPv6 IP is in the subnet",
		},
		{
			ipStr: "2607:f8b0:4005:0900:0000:0000:0000:200e",
			cidr:  "2607:f8b0:4005:0800::/64",
			want:  false,
			name:  "IPv6 IP is not in the subnet",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isIPInSubnet(tt.ipStr, tt.cidr)
			assert.Equal(t, tt.want, result)
		})
	}
}
