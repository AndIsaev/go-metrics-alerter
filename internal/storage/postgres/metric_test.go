package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/AndIsaev/go-metrics-alerter/internal/common"

	_ "github.com/lib/pq"
)

func TestPgStorage_List(t *testing.T) {
	sqlxDB, mock := setupMockDB(t)
	defer sqlxDB.Close()

	storage := &PgStorage{db: sqlxDB}
	ctx := context.Background()

	metrics := expectedMetrics()
	rows := setupMockRows(metrics)

	mock.ExpectQuery("^SELECT (.+) FROM metric*").
		WillReturnRows(rows)

	result, err := storage.List(ctx)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if len(result) != len(metrics) {
		t.Errorf("expected %v results, got %v", len(metrics), len(result))
	}

	for i, metric := range metrics {
		if result[i].ID != metric.ID || result[i].MType != metric.MType {
			t.Errorf("expected result[%d] = %v, got %v", i, metric, result[i])
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPgStorage_GetByName(t *testing.T) {
	sqlxDB, mock := setupMockDB(t)
	defer sqlxDB.Close()

	storage := &PgStorage{db: sqlxDB}
	ctx := context.Background()

	expectedMetric := common.Metrics{
		ID:    "metric1",
		MType: common.Gauge,
		Value: linkFloat64(24.5),
	}

	row := sqlmock.NewRows([]string{"id", "type", "value", "delta"}).
		AddRow(expectedMetric.ID, expectedMetric.MType, expectedMetric.Value, expectedMetric.Delta)

	mock.ExpectQuery("(?i)^select \\* from metric where id = \\$1;?").
		WithArgs(expectedMetric.ID).
		WillReturnRows(row)

	result, err := storage.GetByName(ctx, expectedMetric.ID)

	assert.NoError(t, err)
	assert.Equal(t, expectedMetric, result)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPgStorage_GetByNameType(t *testing.T) {
	sqlxDB, mock := setupMockDB(t)
	defer sqlxDB.Close()

	storage := &PgStorage{db: sqlxDB}
	ctx := context.Background()

	expectedMetric := common.Metrics{
		ID:    "metric1",
		MType: common.Gauge,
		Value: linkFloat64(20.5),
		Delta: nil,
	}

	row := sqlmock.NewRows([]string{"id", "type", "value", "delta"}).
		AddRow(expectedMetric.ID, expectedMetric.MType, expectedMetric.Value, expectedMetric.Delta)

	mock.ExpectQuery(`(?i)^select \* from metric where id = \$1 and "type" = \$2;?`).
		WithArgs(expectedMetric.ID, expectedMetric.MType).
		WillReturnRows(row)

	result, err := storage.GetByNameType(ctx, expectedMetric.ID, expectedMetric.MType)

	assert.NoError(t, err)
	assert.Equal(t, expectedMetric, result)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
func TestPgStorage_Insert(t *testing.T) {
	sqlxDB, mock := setupMockDB(t)
	defer sqlxDB.Close()

	storage := &PgStorage{db: sqlxDB}
	ctx := context.Background()

	inputMetric := common.Metrics{
		ID:    "metric1",
		MType: common.Gauge,
		Value: linkFloat64(42.0),
	}

	expectedMetric := common.Metrics{
		ID:    "metric1",
		MType: common.Gauge,
		Value: linkFloat64(42.0),
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`(?i)INSERT INTO metric \(id, type, delta, value\) VALUES \(\$1, \$2, \$3, \$4\) ON CONFLICT \(id\) DO UPDATE SET delta = metric\.delta \+ EXCLUDED\.delta, value = EXCLUDED\.value RETURNING id, type, delta, value;`).
		WithArgs(inputMetric.ID, inputMetric.MType, inputMetric.Delta, inputMetric.Value).
		WillReturnRows(sqlmock.NewRows([]string{"id", "type", "delta", "value"}).
			AddRow(expectedMetric.ID, expectedMetric.MType, expectedMetric.Delta, expectedMetric.Value))
	mock.ExpectCommit()

	result, err := storage.Insert(ctx, inputMetric)

	assert.NoError(t, err)
	assert.Equal(t, expectedMetric, result)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPgStorage_create(t *testing.T) {
	sqlxDB, mock := setupMockDB(t)
	defer sqlxDB.Close()

	storage := &PgStorage{db: sqlxDB}
	ctx := context.Background()

	inputMetric := common.Metrics{
		ID:    "metric1",
		MType: common.Gauge,
		Value: linkFloat64(42.0),
	}

	mock.ExpectExec(`(?i)insert into metric \(id, type, delta, value\) values \(\$1, \$2, \$3, \$4\) on conflict \(id\) do update set delta = metric\.delta \+ \$3, value = \$4;`).
		WithArgs(inputMetric.ID, inputMetric.MType, inputMetric.Delta, inputMetric.Value).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := storage.create(ctx, inputMetric)

	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPgStorage_createError(t *testing.T) {
	sqlxDB, mock := setupMockDB(t)
	defer sqlxDB.Close()

	storage := &PgStorage{db: sqlxDB}
	ctx := context.Background()

	inputMetric := common.Metrics{
		ID:    "metric1",
		MType: common.Gauge,
		Value: linkFloat64(42.0),
	}

	mock.ExpectExec(`(?i)insert into metric \(id, type, delta, value\) values \(\$1, \$2, \$3, \$4\) on conflict \(id\) do update set delta = metric\.delta \+ \$3, value = \$4;`).
		WithArgs(inputMetric.ID, inputMetric.MType, inputMetric.Delta, inputMetric.Value).
		WillReturnError(fmt.Errorf("error insert row to pg"))

	err := storage.create(ctx, inputMetric)

	assert.Error(t, err, "error insert row to pg")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPgStorage_InsertBatch(t *testing.T) {
	sqlxDB, mock := setupMockDB(t)
	defer sqlxDB.Close()

	storage := &PgStorage{db: sqlxDB}
	ctx := context.Background()

	metrics := []common.Metrics{
		{
			ID:    "metric1",
			MType: common.Gauge,
			Value: linkFloat64(42.0),
		},
		{
			ID:    "metric2",
			MType: common.Counter,
			Delta: linkInt64(10),
		},
	}

	mock.ExpectBegin()

	mock.ExpectPrepare(`(?i)insert into metric \(id, type, delta, value\) values \(\$1, \$2, \$3, \$4\) on conflict \(id\) do update set delta = metric\.delta \+ \$3, value = \$4`)

	for _, metric := range metrics {
		mock.ExpectExec(`(?i)insert into metric`).WithArgs(metric.ID, metric.MType, metric.Delta, metric.Value).
			WillReturnResult(sqlmock.NewResult(1, 1))
	}

	mock.ExpectCommit()

	err := storage.InsertBatch(ctx, metrics)

	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPgStorage_UpsertByValue(t *testing.T) {
	sqlxDB, mock := setupMockDB(t)
	defer sqlxDB.Close()

	storage := &PgStorage{db: sqlxDB}
	ctx := context.Background()

	metric := common.Metrics{
		ID:    "metric1",
		MType: common.Gauge,
	}

	metricValue := 42.0

	expectedMetric := metric
	expectedMetric.Value = &metricValue

	mock.ExpectQuery(`(?i)^select \* from metric where id = \$1;?`).
		WithArgs(metric.ID).
		WillReturnError(sql.ErrNoRows)

	mock.ExpectExec(`insert into metric`).WithArgs(expectedMetric.ID, expectedMetric.MType, expectedMetric.Delta, *expectedMetric.Value).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := storage.UpsertByValue(ctx, metric, metricValue)

	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
