package client

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
)

func SendMetricsClient(client *resty.Client, url string, body []byte) error {
	res, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(url)

	if err != nil {
		return err
	}
	fmt.Println(res)

	if res.StatusCode() != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode())
	}

	return nil
}
