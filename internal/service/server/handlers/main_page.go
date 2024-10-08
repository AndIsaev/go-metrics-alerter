package handlers

import (
	"fmt"
	"net/http"

	"github.com/AndIsaev/go-metrics-alerter/internal/storage"
)

func MainPageHandler(mem *storage.MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Content-Encoding", "gzip")

		for key, val := range mem.Metrics {
			response := fmt.Sprintf("%s: %v", key, val)
			_, err := w.Write([]byte(response + "\n"))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	}
}
