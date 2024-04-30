package request

import (
	"bytes"
	"log"
	"net/http"
)

func SendMetricsHandler(url, contentType string, body []byte) {

	resp, err := http.Post(url, contentType, bytes.NewBuffer(body))
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
		return
	}

	defer resp.Body.Close()
}
