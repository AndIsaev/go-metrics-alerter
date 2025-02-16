package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/AndIsaev/go-metrics-alerter/internal/service/agent/client"
	"github.com/AndIsaev/go-metrics-alerter/internal/service/agent/client/http"
	"github.com/AndIsaev/go-metrics-alerter/internal/service/agent/client/rpc"
)

func TestNew(t *testing.T) {
	app := New()
	assert.NotNil(t, app.Config)
	assert.IsType(t, &Config{}, app.Config)
}

func TestInitRequestClient(t *testing.T) {
	app := New()
	app.Config = NewConfig()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tests := []struct {
		name   string
		http   bool
		client client.RequestClient
	}{
		{name: "http client", http: true, client: &http.Client{}},
		{name: "grpc client", http: false, client: &rpc.Client{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.http {
				app.Config.RPCClient = true
			}
			app.initRequestClient(ctx, cancel)
			assert.IsType(t, tt.client, app.Client)
		})
	}
}

func TestInitHTTPClient(t *testing.T) {
	app := New()
	app.Config = NewConfig()

	app.initHTTPClient()

	assert.IsType(t, &http.Client{}, app.Client)
}

func TestInitGRPCClient(t *testing.T) {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()
	app := New()
	app.Config = NewConfig()

	app.initGRPCClient(cancel)

	assert.IsType(t, &rpc.Client{}, app.Client)
}
