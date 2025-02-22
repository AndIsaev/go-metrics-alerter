package middleware

import (
	"net"
	"net/http"
)

func TrustedSubnetMiddleware(trustIP string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			realIP := r.Header.Get("X-Real-IP")
			if realIP == "" {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			if !isIPInSubnet(realIP, trustIP) {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// isIPInSubnet check trust and input ips
func isIPInSubnet(ipStr, cidr string) bool {
	ip := net.ParseIP(ipStr)
	_, subnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return false
	}

	return subnet.Contains(ip)
}
