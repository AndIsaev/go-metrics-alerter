package handlers

import (
	"log"
	"net/http"

	"github.com/AndIsaev/go-metrics-alerter/internal/storage"
)

func PingHandler(DBConn storage.BaseStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Encoding", "gzip")

		err := DBConn.Ping()
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
