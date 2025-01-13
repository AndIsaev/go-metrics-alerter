package postgres

import (
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/AndIsaev/go-metrics-alerter/internal/storage"
)

// PgStorage storage in postgresql
type PgStorage struct {
	db *sqlx.DB
}

// NewPgStorage init new storage
func NewPgStorage(connString string) (*PgStorage, error) {
	conn, err := sqlx.Connect("pgx", connString)
	if err != nil {
		log.Printf("unable to connect to database: %v\n", err)
		return nil, err
	}
	log.Println("init pg storage")

	return &PgStorage{db: conn}, nil
}

func (p *PgStorage) System() storage.SystemRepository {
	return p
}

func (p *PgStorage) Metric() storage.MetricRepository {
	return p
}
