package storage

import (
	"context"

	"github.com/AndIsaev/go-metrics-alerter/internal/common"
)

type BaseStorage interface {
	Ping() error
	Close() error
	Insert(ctx context.Context, metric common.Metrics) error
	Create(ctx context.Context) error
	InsertBatch(ctx context.Context, metrics *[]common.Metrics) error
	Get(ctx context.Context, metric common.Metrics) (*common.Metrics, error)
}
