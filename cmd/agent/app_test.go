package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestInitHTTPClient(t *testing.T) {
	app := New()
	client := app.initHTTPClient()

	require.NotNil(t, client)
	require.Equal(t, time.Second*5, client.GetClient().Timeout)
}
