package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	_ "github.com/AndIsaev/go-metrics-alerter/internal/service/agent/middleware"
	_ "github.com/AndIsaev/go-metrics-alerter/internal/service/agent/utils"
)

func TestInitHTTPClient(t *testing.T) {
	app := New()
	client := app.initHTTPClient()

	require.NotNil(t, client)
	require.Equal(t, time.Second*5, client.GetClient().Timeout)
}
