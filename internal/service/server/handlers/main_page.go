package handlers

import (
	"fmt"
	"github.com/AndIsaev/go-metrics-alerter/internal/storage"
	"net/http"
)

func MainPageHandler(mem *storage.MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		for key, val := range mem.Metrics {
			response := fmt.Sprintf("%s: %v", key, val)
			_, err := w.Write([]byte(response + "\n"))
			if err != nil {
				return
			}
		}
	}
}
