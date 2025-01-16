package postgres

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/stretchr/testify/assert"
)

func TestPgStorage_Close(t *testing.T) {
	sqlxDB, mock := setupMockDB(t)
	defer sqlxDB.Close()

	storage := &PgStorage{db: sqlxDB}

	t.Run("successful close", func(t *testing.T) {
		mock.ExpectClose().WillReturnError(nil)

		err := storage.Close(context.Background())
		assert.NoError(t, err, "error should be nil on successful close")
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPgStorage_CloseError(t *testing.T) {
	sqlxDB, mock := setupMockDB(t)
	defer sqlxDB.Close()

	storage := &PgStorage{db: sqlxDB}
	assert.NotNil(t, sqlxDB, "sqlxDB should not be nil after setupMockDB")

	t.Run("error on close", func(t *testing.T) {
		mock.ExpectClose().WillReturnError(fmt.Errorf("close error"))

		err := storage.Close(context.Background())
		assert.Error(t, err, "an error should be returned on close failure")
		assert.Equal(t, "close error", err.Error(), "unexpected error message")
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPgStorage_Ping(t *testing.T) {
	sqlxDB, mock := setupMockDB(t)
	defer sqlxDB.Close()

	storage := &PgStorage{db: sqlxDB}
	ctx := context.Background()

	mock.ExpectPing().WillReturnError(nil)

	err := storage.Ping(ctx)
	assert.NoError(t, err, "error should be nil on successful ping")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPgStorage_PingError(t *testing.T) {
	sqlxDB, mock := setupMockDB(t)
	defer sqlxDB.Close()

	storage := &PgStorage{db: sqlxDB}
	ctx := context.Background()

	mock.ExpectPing().WillReturnError(fmt.Errorf("error close connection"))

	err := storage.Ping(ctx)
	assert.Error(t, err, "error close connection")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPgStorage_RunMigrations(t *testing.T) {
	t.Run("successful migration", func(t *testing.T) {
		sqlxDB, mock := setupMockDB(t)
		defer sqlxDB.Close()

		mock.ExpectExec(`create table if not exists metric`).WillReturnResult(sqlmock.NewResult(1, 1))

		storage := &PgStorage{db: sqlxDB}
		err := storage.RunMigrations(context.Background())

		assert.NoError(t, err, "error should be nil on successful migration")

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("migration error", func(t *testing.T) {
		sqlxDB, mock := setupMockDB(t)
		defer sqlxDB.Close()

		mock.ExpectExec(`create table if not exists metric`).WillReturnError(fmt.Errorf("migration failed"))

		storage := &PgStorage{db: sqlxDB}
		err := storage.RunMigrations(context.Background())

		assert.Error(t, err, "error should occur on migration failure")
		assert.Contains(t, err.Error(), "migration failed")

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}
