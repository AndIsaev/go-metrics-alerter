package inmemory

import (
	"context"
	"github.com/AndIsaev/go-metrics-alerter/internal/common"
	"github.com/AndIsaev/go-metrics-alerter/internal/storage"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMemStorage_Ping(t *testing.T) {

	t.Run("success ping", func(t *testing.T) {
		inMemStorage := &MemStorage{Metrics: map[string]common.Metrics{}}
		err := inMemStorage.Ping(context.Background())
		assert.NoError(t, err)
	})

	t.Run(storage.ErrMapNotAvailable.Error(), func(t *testing.T) {
		inMemStorage := &MemStorage{}
		err := inMemStorage.Ping(context.Background())
		assert.ErrorIs(t, err, storage.ErrMapNotAvailable)
	})
}

func TestMemStorage_Close(t *testing.T) {

	t.Run("success close", func(t *testing.T) {
		inMemStorage := &MemStorage{Metrics: map[string]common.Metrics{}}
		err := inMemStorage.Close(context.Background())
		assert.NoError(t, err)
	})
}
