package middleware

import (
	"bytes"
	"io"
	"log"
	"net/http"

	"github.com/AndIsaev/go-metrics-alerter/internal/common"
)

func SecretMiddleware(key string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if key != "" {
				agentSha256sum := r.Header.Get("HashSHA256")

				body, err := io.ReadAll(r.Body)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}

				r.Body = io.NopCloser(bytes.NewBuffer(body))

				defer r.Body.Close()

				serverSha256sum := common.Sha256sum(body, key)

				if agentSha256sum != serverSha256sum {
					log.Println("compare hash is not success")
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				w.Header().Set("HashSHA256", serverSha256sum)
			}

			next.ServeHTTP(w, r)
		})
	}
}
