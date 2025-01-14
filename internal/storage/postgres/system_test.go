package postgres

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPgStorage_Ping(t *testing.T) {
	tests := []struct {
		name  string
		setup func(suite *testSuite)
		want  error
	}{
		{
			name:  "success ping",
			setup: func(ts *testSuite) { ts.mockSystemRepo.EXPECT().Ping(ts.ctx).Return(nil) },
			want:  nil,
		},
		{
			name:  "error ping",
			setup: func(ts *testSuite) { ts.mockSystemRepo.EXPECT().Ping(ts.ctx).Return(sql.ErrConnDone) },
			want:  sql.ErrConnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := setupTest(t)

			tt.setup(ts)

			s := ts.mockStorage

			got := s.System().Ping(ts.ctx)

			assert.Equal(t, tt.want, got)

		})
	}
}

func TestPgStorage_Close(t *testing.T) {
	tests := []struct {
		name  string
		setup func(suite *testSuite)
		want  error
	}{
		{
			name:  "success close",
			setup: func(ts *testSuite) { ts.mockSystemRepo.EXPECT().Close(ts.ctx).Return(nil) },
			want:  nil,
		},
		{
			name:  "error close",
			setup: func(ts *testSuite) { ts.mockSystemRepo.EXPECT().Close(ts.ctx).Return(sql.ErrConnDone) },
			want:  sql.ErrConnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := setupTest(t)

			tt.setup(ts)

			s := ts.mockStorage

			got := s.System().Close(ts.ctx)

			assert.Equal(t, tt.want, got)

		})
	}
}

func TestPgStorage_RunMigrations(t *testing.T) {
	tests := []struct {
		name  string
		setup func(suite *testSuite)
		want  error
	}{
		{
			name:  "success close",
			setup: func(ts *testSuite) { ts.mockSystemRepo.EXPECT().RunMigrations(ts.ctx).Return(nil) },
			want:  nil,
		},
		{
			name:  "error close",
			setup: func(ts *testSuite) { ts.mockSystemRepo.EXPECT().RunMigrations(ts.ctx).Return(sql.ErrConnDone) },
			want:  sql.ErrConnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := setupTest(t)

			tt.setup(ts)

			s := ts.mockStorage

			got := s.System().RunMigrations(ts.ctx)

			assert.Equal(t, tt.want, got)

		})
	}
}
