package storage

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"

	"github.com/AndIsaev/go-metrics-alerter/internal/common"
)

type PostgresStorage struct {
	Conn *pgx.Conn
}

func NewPostgresStorage(ctx context.Context, connString string) (*PostgresStorage, error) {
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		log.Printf("unable to connect to database: %s\n", err.Error())
		return nil, err
	}
	return &PostgresStorage{Conn: conn}, nil
}

func (s *PostgresStorage) Ping(ctx context.Context) error {
	err := s.Conn.Ping(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresStorage) Insert(ctx context.Context, m common.Metrics) error {
	query := `insert into metric (id, type, delta, value) 
				values ($1, $2, $3, $4) on CONFLICT (id) 
				do update set delta = metric.delta + $3, value = $4;`

	_, err := s.Conn.Exec(ctx, query, m.ID, m.MType, m.Delta, m.Value)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresStorage) Close(ctx context.Context) error {
	err := s.Conn.Close(ctx)
	return err
}

func (s *PostgresStorage) Create(ctx context.Context) error {
	queryMetricTable := `create table if not exists metric(
								id varchar(200) unique not null, 
								"type" varchar(50) not null, 
								delta integer, 
								"value" double precision);`

	_, err := s.Conn.Exec(ctx, queryMetricTable)
	if err != nil {
		log.Printf("can't execute query because of %s\n", err.Error())
		return err
	}
	return nil
}
