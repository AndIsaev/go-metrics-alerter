package client

import (
	"fmt"
	"github.com/go-resty/resty/v2"
)

func SendMetricsClient(client *resty.Client, url string, body []byte) error {
	_, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(url)

	if err != nil {
		fmt.Println("+++++++++++++++++++")
		fmt.Println(err)
		fmt.Println("+++++++++++++++++++")
		return err
	}

	return nil
}
