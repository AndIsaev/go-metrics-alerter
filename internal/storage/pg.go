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
				values ($1, $2, $3, $4) on conflict (id) 
				do update set delta = metric.delta + $3, value = $4;`

	_, err := s.Conn.Exec(ctx, query, m.ID, m.MType, m.Delta, m.Value)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresStorage) Get(ctx context.Context, m common.Metrics) (*common.Metrics, error) {
	result := new(common.Metrics)

	query := `select * from metric where id = $1 and "type" = $2;`

	row := s.Conn.QueryRow(ctx, query, m.ID, m.MType)

	if err := row.Scan(&result.ID, &result.MType, &result.Delta, &result.Value); err != nil {
		return nil, err
	}

	return result, nil
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

func (s *PostgresStorage) InsertBatch(ctx context.Context, metrics []common.Metrics) error {
	if len(metrics) == 0 {
		return nil
	}
	tx, err := s.Conn.Begin(ctx)
	if err != nil {
		return err
	}

	query := `insert into metric (id, type, delta, value) 
				values ($1, $2, $3, $4) on conflict (id) 
				do update set delta = metric.delta + $3, value = $4;`

	for _, m := range metrics {
		// все изменения записываются в транзакцию
		if _, err := tx.Exec(ctx, query, m.ID, m.MType, m.Delta, m.Value); err != nil {
			// если ошибка, то откатываем изменения
			tx.Rollback(ctx)
			return err
		}
	}
	// завершаем транзакцию
	return tx.Commit(ctx)
}
