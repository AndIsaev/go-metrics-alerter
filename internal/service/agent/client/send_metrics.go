package client

import (
	"bytes"
	"errors"
	"log"
	"net/http"
)

func SendMetricsClient(url, contentType string, body []byte) error {
	ErrorSendMetricsHandler := errors.New("connection error")

	resp, err := http.Post(url, contentType, bytes.NewBuffer(body))
	if err != nil {
		log.Fatalf("an Error Occurred %v", err)
		return ErrorSendMetricsHandler
	}

	defer resp.Body.Close()
	return nil
}
