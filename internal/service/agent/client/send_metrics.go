package client

import (
	"bytes"
	"log"
	"net/http"
)

func SendMetricsClient(url, contentType string, body []byte) error {
	resp, err := http.Post(url, contentType, bytes.NewBuffer(body))
	if err != nil {
		log.Fatalf("an Error Occurred %v", err)
		return nil
	}

	defer resp.Body.Close()
	return nil
}
