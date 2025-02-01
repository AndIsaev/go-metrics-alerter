package middleware

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
