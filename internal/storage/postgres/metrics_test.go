package postgres

import (
	"context"
	"database/sql"
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
	// 1. Создание мок базы данных и SQLx обертки
	sqlxDB, mock := setupMockDB(t)
	defer sqlxDB.Close()

	storage := &PgStorage{db: sqlxDB}
	ctx := context.Background()

	// 3. Ожидаемый результат
	expectedMetric := common.Metrics{
		ID:    "metric1",
		MType: common.Gauge,
		Value: linkFloat64(24.5),
	}

	row := sqlmock.NewRows([]string{"id", "type", "value", "delta"}).
		AddRow(expectedMetric.ID, expectedMetric.MType, expectedMetric.Value, expectedMetric.Delta)

	// 4. Настройка ожиданий мока
	mock.ExpectQuery("(?i)^select \\* from metric where id = \\$1;?").
		WithArgs(expectedMetric.ID).
		WillReturnRows(row)

	// 5. Вызов тестируемой функции
	result, err := storage.GetByName(ctx, expectedMetric.ID)

	// 6. Проверка результатов
	assert.NoError(t, err)
	assert.Equal(t, expectedMetric, result)

	// 7. Проверка, что все ожидания были выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPgStorage_GetByNameType(t *testing.T) {
	// 1. Создание мок базы данных и SQLx обертки
	sqlxDB, mock := setupMockDB(t)
	defer sqlxDB.Close()

	storage := &PgStorage{db: sqlxDB}
	ctx := context.Background()

	// 3. Ожидаемый результат
	expectedMetric := common.Metrics{
		ID:    "metric1",
		MType: common.Gauge,
		Value: linkFloat64(20.5),
		Delta: nil,
	}

	row := sqlmock.NewRows([]string{"id", "type", "value", "delta"}).
		AddRow(expectedMetric.ID, expectedMetric.MType, expectedMetric.Value, expectedMetric.Delta)

	// 4. Настройка ожиданий мока
	mock.ExpectQuery(`(?i)^select \* from metric where id = \$1 and "type" = \$2;?`).
		WithArgs(expectedMetric.ID, expectedMetric.MType).
		WillReturnRows(row)

	// 5. Вызов тестируемой функции
	result, err := storage.GetByNameType(ctx, expectedMetric.ID, expectedMetric.MType)

	// 6. Проверка результатов
	assert.NoError(t, err)
	assert.Equal(t, expectedMetric, result)

	// 7. Проверка, что все ожидания были выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
func TestPgStorage_Insert(t *testing.T) {
	// 1. Создание мок базы данных и SQLx обертки
	sqlxDB, mock := setupMockDB(t)
	defer sqlxDB.Close()

	storage := &PgStorage{db: sqlxDB}
	ctx := context.Background()

	// 2. Метрика, которую будем вставлять
	inputMetric := common.Metrics{
		ID:    "metric1",
		MType: common.Gauge,
		Value: linkFloat64(42.0),
	}

	// 3. Ожидаемый результат
	expectedMetric := common.Metrics{
		ID:    "metric1",
		MType: common.Gauge,
		Value: linkFloat64(42.0),
	}

	// 4. Настройка ожиданий мока
	mock.ExpectBegin()
	mock.ExpectQuery(`(?i)INSERT INTO metric \(id, type, delta, value\) VALUES \(\$1, \$2, \$3, \$4\) ON CONFLICT \(id\) DO UPDATE SET delta = metric\.delta \+ EXCLUDED\.delta, value = EXCLUDED\.value RETURNING id, type, delta, value;`).
		WithArgs(inputMetric.ID, inputMetric.MType, inputMetric.Delta, inputMetric.Value).
		WillReturnRows(sqlmock.NewRows([]string{"id", "type", "delta", "value"}).
			AddRow(expectedMetric.ID, expectedMetric.MType, expectedMetric.Delta, expectedMetric.Value))
	mock.ExpectCommit()

	// 5. Вызов тестируемой функции
	result, err := storage.Insert(ctx, inputMetric)

	// 6. Проверка результатов
	assert.NoError(t, err)
	assert.Equal(t, expectedMetric, result)

	// 7. Проверка, что все ожидания были выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPgStorage_create(t *testing.T) {
	// 1. Создание мок базы данных и SQLx обертки
	sqlxDB, mock := setupMockDB(t)
	defer sqlxDB.Close()

	storage := &PgStorage{db: sqlxDB}
	ctx := context.Background()

	// 2. Метрика, которую будем использовать
	inputMetric := common.Metrics{
		ID:    "metric1",
		MType: common.Gauge,
		Value: linkFloat64(42.0),
	}

	// 3. Настройка ожиданий мока
	mock.ExpectExec(`(?i)insert into metric \(id, type, delta, value\) values \(\$1, \$2, \$3, \$4\) on conflict \(id\) do update set delta = metric\.delta \+ \$3, value = \$4;`).
		WithArgs(inputMetric.ID, inputMetric.MType, inputMetric.Delta, inputMetric.Value).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// 4. Вызов тестируемой функции
	err := storage.create(ctx, inputMetric)

	// 5. Проверка отсутствия ошибок
	assert.NoError(t, err)

	// 6. Проверка, что все ожидания были выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPgStorage_InsertBatch(t *testing.T) {
	// Создание мок базы данных и SQLx обертки
	sqlxDB, mock := setupMockDB(t)
	defer sqlxDB.Close()

	storage := &PgStorage{db: sqlxDB}
	ctx := context.Background()

	// Создание метрик для теста
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

	// Настройка ожиданий мока
	mock.ExpectBegin()

	// Ожидание подготовки SQL-запроса
	mock.ExpectPrepare(`(?i)insert into metric \(id, type, delta, value\) values \(\$1, \$2, \$3, \$4\) on conflict \(id\) do update set delta = metric\.delta \+ \$3, value = \$4`)

	// Ожидание выполнения SQL-запроса для каждой метрики
	for _, metric := range metrics {
		mock.ExpectExec(`(?i)insert into metric`).WithArgs(metric.ID, metric.MType, metric.Delta, metric.Value).
			WillReturnResult(sqlmock.NewResult(1, 1))
	}

	// Ожидание коммита транзакции
	mock.ExpectCommit()

	// Вызов тестируемой функции
	err := storage.InsertBatch(ctx, metrics)

	// Проверка отсутствия ошибок
	assert.NoError(t, err)

	// Проверка, что все ожидания были выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPgStorage_UpsertByValue(t *testing.T) {
	// Создание мок базы данных и SQLx обертки
	sqlxDB, mock := setupMockDB(t)
	defer sqlxDB.Close()

	storage := &PgStorage{db: sqlxDB}
	ctx := context.Background()

	// Пример метрики для тестов
	metric := common.Metrics{
		ID:    "metric1",
		MType: common.Gauge,
	}

	metricValue := 42.0

	// Установка значений для тестирования установщика
	expectedMetric := metric
	expectedMetric.Value = &metricValue

	// Создание сценария для случая, когда метрика не существует
	mock.ExpectQuery(`(?i)^select \* from metric where id = \$1;?`).
		WithArgs(metric.ID).
		WillReturnError(sql.ErrNoRows)

	// Ожидание выполнения SQL-запроса на создание метрики
	mock.ExpectExec(`insert into metric`).WithArgs(expectedMetric.ID, expectedMetric.MType, expectedMetric.Delta, *expectedMetric.Value).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Вызов тестируемой функции
	err := storage.UpsertByValue(ctx, metric, metricValue)

	// Проверка отсутствия ошибок
	assert.NoError(t, err)

	// Проверка, что все ожидания были выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
