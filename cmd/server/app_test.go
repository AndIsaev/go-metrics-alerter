package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInitStorage(t *testing.T) {
	app := New()
	storage := app.initStorage

	require.NotNil(t, storage)
}
