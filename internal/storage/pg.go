package storage

import (
	"context"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/AndIsaev/go-metrics-alerter/internal/common"
)

type PostgresStorage struct {
	Db *sqlx.DB
}

func NewPostgresStorage(connString string) (*PostgresStorage, error) {
	conn, err := sqlx.Connect("postgres", connString)
	if err != nil {
		log.Printf("unable to connect to database: %s\n", err.Error())
		return nil, err
	}
	return &PostgresStorage{Db: conn}, nil
}

func (s *PostgresStorage) Ping() error {
	err := s.Db.Ping()
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresStorage) Insert(ctx context.Context, m common.Metrics) error {
	query := `insert into metric (id, type, delta, value) 
				values ($1, $2, $3, $4) on conflict (id) 
				do update set delta = metric.delta + $3, value = $4;`

	_, err := s.Db.ExecContext(ctx, query, m.ID, m.MType, m.Delta, m.Value)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresStorage) Get(ctx context.Context, m common.Metrics) (*common.Metrics, error) {
	result := common.Metrics{}

	query := `select * from metric where id = $1 and "type" = $2;`

	if err := s.Db.GetContext(ctx, &result, query, m.ID, m.MType); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &result, nil
}

func (s *PostgresStorage) Close() error {
	err := s.Db.Close()
	return err
}

func (s *PostgresStorage) Create(ctx context.Context) error {
	queryMetricTable := `create table if not exists metric(
								id varchar(200) unique not null, 
								"type" varchar(50) not null, 
								delta bigint, 
								"value" double precision);`

	_, err := s.Db.ExecContext(ctx, queryMetricTable)
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
	tx, err := s.Db.Begin()
	if err != nil {
		return err
	}

	query := `insert into metric (id, type, delta, value) 
				values ($1, $2, $3, $4) on conflict (id) 
				do update set delta = metric.delta + $3, value = $4;`

	for _, m := range metrics {
		// все изменения записываются в транзакцию
		if _, err := tx.ExecContext(ctx, query, m.ID, m.MType, m.Delta, m.Value); err != nil {
			// если ошибка, то откатываем изменения
			log.Println("+++++++")
			log.Println(err)
			log.Println("+++++++")

			tx.Rollback()
			return err
		}
	}
	// завершаем транзакцию
	return tx.Commit()
}
