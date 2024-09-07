package storage

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/AndIsaev/go-metrics-alerter/internal/common"
)

type PostgresStorage struct {
	DB *sqlx.DB
}

func NewPostgresStorage(connString string) (*PostgresStorage, error) {
	conn, err := sqlx.Connect("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("%w\n", err)
	}
	return &PostgresStorage{DB: conn}, nil
}

func (s *PostgresStorage) Ping() error {
	err := s.DB.Ping()
	if err != nil {
		return fmt.Errorf("%w\n", err)
	}
	return nil
}

func (s *PostgresStorage) Insert(ctx context.Context, m common.Metrics) error {
	query := `insert into metric (id, type, delta, value) 
				values ($1, $2, $3, $4) on conflict (id) 
				do update set delta = metric.delta + $3, value = $4;`

	_, err := s.DB.ExecContext(ctx, query, m.ID, m.MType, m.Delta, m.Value)
	if err != nil {
		return fmt.Errorf("%w\n", err)
	}
	return nil
}

func (s *PostgresStorage) Get(ctx context.Context, m common.Metrics) (*common.Metrics, error) {
	result := common.Metrics{}

	query := `select * from metric where id = $1 and "type" = $2;`

	if err := s.DB.GetContext(ctx, &result, query, m.ID, m.MType); err != nil {
		return nil, fmt.Errorf("%w\n", err)
	}

	return &result, nil
}

func (s *PostgresStorage) Close() error {
	err := s.DB.Close()
	return fmt.Errorf("%w\n", err)
}

func (s *PostgresStorage) Create(ctx context.Context) error {
	queryMetricTable := `create table if not exists metric(
								id varchar(200) unique not null, 
								"type" varchar(50) not null, 
								delta bigint, 
								"value" double precision);`

	_, err := s.DB.ExecContext(ctx, queryMetricTable)
	if err != nil {
		return fmt.Errorf("%w\n", err)
	}
	return nil
}

func (s *PostgresStorage) InsertBatch(ctx context.Context, metrics []common.Metrics) error {
	if len(metrics) == 0 {
		return nil
	}
	tx, err := s.DB.Begin()
	if err != nil {
		return fmt.Errorf("%w\n", err)
	}

	query := `insert into metric (id, type, delta, value) 
				values ($1, $2, $3, $4) on conflict (id) 
				do update set delta = metric.delta + $3, value = $4;`

	for _, m := range metrics {
		if _, err := tx.ExecContext(ctx, query, m.ID, m.MType, m.Delta, m.Value); err != nil {
			tx.Rollback()
			return fmt.Errorf("%w\n", err)
		}
	}
	return tx.Commit()
}
