package storage

import (
	"context"
	"github.com/AndIsaev/go-metrics-alerter/internal/common"
)

type BaseStorage interface {
	Ping(ctx context.Context) error
	Close(ctx context.Context) error
	Insert(ctx context.Context, metrics common.Metrics) error
	Create(ctx context.Context) error
}
