package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/AndIsaev/go-metrics-alerter/internal/service/agent/utils"

	"github.com/AndIsaev/go-metrics-alerter/internal/common"
)

// SendMetric using the method to establish indicators
func (c *Client) SendMetric(url string, body common.Metrics) error {
	var result common.Metrics

	res, err := c.Client.R().
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

func (c *Client) SendMetrics(ctx context.Context, metrics []common.Metrics) error {
	select {
	case <-ctx.Done():
		log.Println("context is done -> exit from SendMetrics")
		return nil
	default:
		ip, err := c.IPResolver.GetLocalIP()
		if err != nil {
			log.Printf("Error getting local IP: %v\n", err)
			return err
		}

		body, err := json.Marshal(metrics)
		if err != nil {
			return errors.Unwrap(fmt.Errorf("error encoding metric: %w", err))
		}

		client := c.Client.R()
		if c.SecretKey != "" {
			sha256sum := common.Sha256sum(body, c.SecretKey)
			client.SetHeader("HashSHA256", sha256sum)
		}

		if c.PublicKey != nil {
			body, err = utils.Encrypt(c.PublicKey, body)
			if err != nil {
				return fmt.Errorf("error encrypting metrics: %w", errors.Unwrap(err))
			}
		}

		log.Println("send metrics")

		res, err := client.
			SetHeader("Accept-Encoding", "gzip").
			SetHeader("Content-Type", "application/octet-stream").
			SetHeader("X-Real-IP", ip).
			SetBody(body).
			Post(c.URL)

		if err != nil {
			log.Printf("error sending request: %v\n", err)
			return fmt.Errorf("error sending request: %w", err)
		}

		if res.StatusCode() != http.StatusOK {
			log.Printf("error sending request: status: %v, response: %v\n", res.StatusCode(), res)
			return fmt.Errorf("error sending request: response: %v", res)
		}
	}

	return nil
}
