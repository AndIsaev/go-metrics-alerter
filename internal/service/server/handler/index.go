package handler

import (
	"encoding/json"
	"net/http"
)

// IndexHandler main page
func (h *Handler) IndexHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		metrics, err := h.MetricService.ListMetrics(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		response, _ := json.Marshal(metrics)

		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}
