package middleware

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"io"
	"log"
	"net/http"
)

// DecryptMiddleware - decrypt request body
func DecryptMiddleware(key *rsa.PrivateKey) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				log.Println("Error reading request body:", err)
				http.Error(w, "Invalid request", http.StatusBadRequest)
				return
			}
			defer r.Body.Close()

			plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, key, body)
			if err != nil {
				log.Println("Error decrypting data:", err)
				http.Error(w, "Could not decrypt data", http.StatusBadRequest)
				return
			}

			r.Body = io.NopCloser(bytes.NewBuffer(plaintext))
			next.ServeHTTP(w, r)
		})
	}
}
