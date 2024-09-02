package handlers

import (
	"context"
	"net/http"

	"github.com/AndIsaev/go-metrics-alerter/internal/storage"
)

func PingHandler(DBConn storage.BaseStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Encoding", "gzip")

		err := DBConn.Ping(context.Background())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
