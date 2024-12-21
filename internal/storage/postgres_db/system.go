package postgres_db

import (
	"context"
	"errors"
	"fmt"
	"log"
)

func (p *PgStorage) Close(_ context.Context) error {
	err := p.db.Close()
	if err != nil {
		log.Println(errors.Unwrap(err))
	}
	return err
}

func (p *PgStorage) Ping(_ context.Context) error {
	err := p.db.Ping()
	if err != nil {
		log.Println(errors.Unwrap(err))
	}
	return err
}

func (p *PgStorage) RunMigrations(ctx context.Context) error {
	log.Println("run migrations")

	query := `create table if not exists metric(
								id varchar(200) unique not null, 
								"type" varchar(50) not null, 
								delta bigint, 
								"value" double precision);`

	_, err := p.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	log.Println("migrations completed")

	return nil
}
