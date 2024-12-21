package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

func (h *Handler) IndexHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		metrics, err := h.MetricService.ListMetrics(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		response, err := json.Marshal(metrics)
		if err != nil {
			log.Println("error serialize metrics")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}
