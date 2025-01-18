package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInitStorage(t *testing.T) {
	app := New()
	storage := app.initStorage

	require.NotNil(t, storage)
}
