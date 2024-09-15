package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/AndIsaev/go-metrics-alerter/internal/common"
	"github.com/AndIsaev/go-metrics-alerter/internal/service/agent"
	"github.com/AndIsaev/go-metrics-alerter/internal/service/agent/middleware"
	"github.com/AndIsaev/go-metrics-alerter/internal/service/agent/secure"
	"github.com/go-resty/resty/v2"
	"net/http"
	"time"
)

// AgentApp - structure of application
type AgentApp struct {
	Config *agent.Config
	Client *resty.Client
}

func New() *AgentApp {
	app := &AgentApp{}
	config := agent.NewConfig()
	app.Config = config

	return app
}

func (a *AgentApp) StartApp() {
	a.Client = a.initHTTPClient()
}

func (a *AgentApp) initHTTPClient() *resty.Client {
	cli := resty.New()
	cli.SetTimeout(time.Second * 5)
	cli.OnBeforeRequest(middleware.GzipRequestMiddleware)
	cli.OnBeforeRequest(a.HashMiddleware)
	return cli
}

func (a *AgentApp) SendMetrics() error {
	values := make([]common.Metrics, 0, 100)
	var result common.Metrics

	for _, v := range a.Config.StorageMetrics.Metrics {
		metric := common.Metrics{ID: v.ID, MType: v.MType, Value: v.Value, Delta: v.Delta}
		values = append(values, metric)
	}
	res, err := a.Client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(&values).
		SetResult(&result).
		Post(a.Config.UpdateMetricsAddress)

	if err != nil {
		return errors.Unwrap(err)
	}

	if res.StatusCode() != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode())
	}

	return nil
}

func (a *AgentApp) HashMiddleware(c *resty.Client, r *resty.Request) error {
	if a.Config.Key != "" {
		switch value := r.Body.(type) {
		case *[]common.Metrics:

			v, err := json.Marshal(value)
			if err != nil {
				return err
			}
			sha256sum := secure.Sha256sum(v, a.Config.Key)
			fmt.Println(sha256sum)
			c.Header.Set("HashSHA256", sha256sum)
		}
	}
	return nil
}
