package postgres

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"

	"github.com/AndIsaev/go-metrics-alerter/internal/common"
)

// Функция для создания мок базы данных и Sqlx обертки
func setupMockDB(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	return sqlxDB, mock
}

// Функция для создания ожидаемых результатов
func expectedMetrics() []common.Metrics {
	return []common.Metrics{
		{ID: "metric1", MType: "gauge", Value: common.LinkFloat64(24.5)},
		{ID: "metric2", MType: "counter", Delta: common.LinkInt64(5)},
	}
}

func setupMockRows(metrics []common.Metrics) *sqlmock.Rows {
	return sqlmock.NewRows([]string{"id", "type", "value", "delta"}).
		AddRow(metrics[0].ID, metrics[0].MType, metrics[0].Value, metrics[0].Delta).
		AddRow(metrics[1].ID, metrics[1].MType, metrics[1].Value, metrics[1].Delta)
}
