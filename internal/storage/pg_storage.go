package storage

import "context"

type PgStorage interface {
	Ping(ctx context.Context) error
	Close(ctx context.Context) error
}
