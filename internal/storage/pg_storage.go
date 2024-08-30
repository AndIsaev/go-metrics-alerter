package storage

import (
	"context"
	"github.com/jackc/pgx/v5/pgconn"
)

type PgStorage interface {
	Ping(ctx context.Context) error
	Close(ctx context.Context) error
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}
