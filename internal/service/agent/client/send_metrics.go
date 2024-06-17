package client

import (
	"fmt"
	"github.com/AndIsaev/go-metrics-alerter/internal/common"
	"github.com/go-resty/resty/v2"
	"net/http"
)

func SendMetricsClient(client *resty.Client, url string, body common.Metrics) error {
	var result common.Metrics

	res, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		SetResult(&result).
		Post(url)

	if err != nil {
		return err
	}

	if res.StatusCode() != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode())
	}

	return nil
}
