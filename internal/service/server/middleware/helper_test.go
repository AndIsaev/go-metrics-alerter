package middleware

import (
	"io"
	"net/http"
)

// Обработчик, который просто читает тело текста и возвращает его.
func echoHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	w.Write(body)
}
