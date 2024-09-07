package client

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-resty/resty/v2"

	"github.com/AndIsaev/go-metrics-alerter/internal/common"
)

func SendMetricHandler(client *resty.Client, url string, body common.Metrics) error {
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

func SendMetricsHandler(client *resty.Client, url string, body *[]common.Metrics) error {
	var result common.Metrics

	log.Println("send metrics")
	res, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		SetResult(&result).
		Post(url)

	if err != nil {
		return fmt.Errorf("%w", err)
	}

	if res.StatusCode() != http.StatusOK {
		return fmt.Errorf("%w", err)
	}

	return nil
}
