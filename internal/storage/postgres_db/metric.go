package postgres_db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/AndIsaev/go-metrics-alerter/internal/common"
	"github.com/AndIsaev/go-metrics-alerter/internal/storage"
)

func (p *PgStorage) List(ctx context.Context) ([]common.Metrics, error) {
	var metrics []common.Metrics
	query := "SELECT * FROM metric"

	err := p.db.SelectContext(ctx, &metrics, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	return metrics, nil
}

func (p *PgStorage) UpsertByValue(ctx context.Context, metric common.Metrics, metricValue any) error {
	newValue := storage.MetricValue{}

	if err := newValue.Set(metric.MType, metricValue); err != nil {
		return err
	}

	existsMetric, err := p.GetByName(ctx, metric.ID)
	if !errors.Is(err, sql.ErrNoRows) && err != nil {
		return err
	}
	if errors.Is(err, sql.ErrNoRows) {
		existsMetric = metric
	}

	switch metric.MType {
	case common.Gauge:
		existsMetric.Value = &newValue.FloatValue
		return p.Create(ctx, existsMetric)

	case common.Counter:
		existsMetric.Delta = &newValue.IntValue
		return p.Create(ctx, existsMetric)
	}

	return nil
}

func (p *PgStorage) GetByName(ctx context.Context, name string) (common.Metrics, error) {
	result := common.Metrics{}

	query := `select * from metric where id = $1;`

	if err := p.db.GetContext(ctx, &result, query, name); err != nil {
		return result, err
	}

	return result, nil
}

func (p *PgStorage) Create(ctx context.Context, metric common.Metrics) error {
	query := `insert into metric (id, type, delta, value)
					values ($1, $2, $3, $4) on conflict (id)
					do update set delta = metric.delta + $3, value = $4;`

	_, err := p.db.ExecContext(ctx, query, metric.ID, metric.MType, metric.Delta, metric.Value)
	if err != nil {
		log.Println("error insert row to pg")
		return err
	}
	return nil
}

func (p *PgStorage) InsertBatch(ctx context.Context, metrics []common.Metrics) error {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	stmt, err := tx.PrepareContext(ctx, `
        insert into metric (id, type, delta, value) 
        values ($1, $2, $3, $4) 
        on conflict (id) 
        do update set delta = metric.delta + $3, value = $4;
    `)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, m := range metrics {
		if _, err := stmt.ExecContext(ctx, m.ID, m.MType, m.Delta, m.Value); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("failed to execute statement: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (p *PgStorage) GetByNameType(ctx context.Context, name, mType string) (common.Metrics, error) {
	result := common.Metrics{}

	query := `select * from metric where id = $1 and "type" = $2;`

	if err := p.db.GetContext(ctx, &result, query, name, mType); err != nil {
		return result, err
	}

	return result, nil
}

func (p *PgStorage) Insert(ctx context.Context, metric common.Metrics) (common.Metrics, error) {
	tx, err := p.db.Begin()
	if err != nil {
		return metric, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	query := `
			  INSERT INTO metric (id, type, delta, value)
			  VALUES ($1, $2, $3, $4)
			  ON CONFLICT (id)
			  DO UPDATE SET
			   delta = metric.delta + EXCLUDED.delta,
			   value = EXCLUDED.value
			  RETURNING id, type, delta, value; `

	row := tx.QueryRowContext(
		ctx,
		query,
		metric.ID,
		metric.MType,
		metric.Delta,
		metric.Value,
	)

	err = row.Scan(&metric.ID, &metric.MType, &metric.Delta, &metric.Value)
	if err != nil {
		return metric, err
	}

	return metric, nil
}
