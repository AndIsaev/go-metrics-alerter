package server

import (
	"database/sql"
	"testing"

	"github.com/AndIsaev/go-metrics-alerter/internal/common"
)

func TestMethods_PingStorage(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(ts *testSuite)
		wantErr bool
	}{
		{
			name: "success ping",
			setup: func(ts *testSuite) {
				ts.mockSystemRepo.EXPECT().
					Ping(ts.ctx).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error ping",
			setup: func(ts *testSuite) {
				ts.mockSystemRepo.EXPECT().
					Ping(ts.ctx).
					Return(sql.ErrConnDone)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := setupTest(t)

			tt.setup(ts)

			s := &Methods{
				Storage: ts.mockStorage,
			}
			err := s.PingStorage(ts.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("PingStorage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestMethods_CloseStorage(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(ts *testSuite)
		wantErr bool
	}{
		{
			name: "success close",
			setup: func(ts *testSuite) {
				ts.mockSystemRepo.EXPECT().
					Close(ts.ctx).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error close",
			setup: func(ts *testSuite) {
				ts.mockSystemRepo.EXPECT().
					Close(ts.ctx).
					Return(sql.ErrConnDone)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := setupTest(t)

			tt.setup(ts)

			s := &Methods{
				Storage: ts.mockStorage,
			}
			err := s.CloseStorage(ts.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("CloseStorage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestMethods_RunMigrationsStorage(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(ts *testSuite)
		wantErr bool
	}{
		{
			name: "success migrations",
			setup: func(ts *testSuite) {
				ts.mockSystemRepo.EXPECT().
					RunMigrations(ts.ctx).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error migrations",
			setup: func(ts *testSuite) {
				ts.mockSystemRepo.EXPECT().
					RunMigrations(ts.ctx).
					Return(sql.ErrConnDone)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := setupTest(t)

			tt.setup(ts)

			s := &Methods{
				Storage: ts.mockStorage,
			}
			err := s.RunMigrationsStorage(ts.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("RunMigrationsStorage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestMethods_ListMetrics(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(ts *testSuite)
		wantErr bool
	}{
		{
			name: "success list",
			setup: func(ts *testSuite) {
				ts.mockMetricRepo.EXPECT().
					List(ts.ctx).
					Return(nil, nil)
			},
			wantErr: false,
		},
		{
			name: "error list",
			setup: func(ts *testSuite) {
				ts.mockMetricRepo.EXPECT().
					List(ts.ctx).
					Return(nil, sql.ErrConnDone)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := setupTest(t)

			tt.setup(ts)

			s := &Methods{
				Storage: ts.mockStorage,
			}
			_, err := s.ListMetrics(ts.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListMetrics() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestMethods_UpdateMetricByValue(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(ts *testSuite)
		wantErr bool
	}{
		{
			name: "success upsert by value",
			setup: func(ts *testSuite) {
				ts.mockMetricRepo.EXPECT().
					UpsertByValue(ts.ctx, common.Metrics{}, 1).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error upsert by value",
			setup: func(ts *testSuite) {
				ts.mockMetricRepo.EXPECT().
					UpsertByValue(ts.ctx, common.Metrics{}, 1).
					Return(sql.ErrConnDone)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := setupTest(t)

			tt.setup(ts)

			s := &Methods{
				Storage: ts.mockStorage,
			}
			err := s.UpdateMetricByValue(ts.ctx, common.Metrics{}, 1)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateMetricByValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestMethods_GetMetricByName(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(ts *testSuite)
		wantErr bool
	}{
		{
			name: "success get by name",
			setup: func(ts *testSuite) {
				ts.mockMetricRepo.EXPECT().
					GetByName(ts.ctx, "metric1").
					Return(common.Metrics{ID: "metric1", MType: common.Gauge}, nil)
			},
			wantErr: false,
		},
		{
			name: "error get by name",
			setup: func(ts *testSuite) {
				ts.mockMetricRepo.EXPECT().
					GetByName(ts.ctx, "metric1").
					Return(common.Metrics{}, sql.ErrConnDone)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := setupTest(t)

			tt.setup(ts)

			s := &Methods{
				Storage: ts.mockStorage,
			}
			_, err := s.GetMetricByName(ts.ctx, "metric1")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMetricByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestMethods_GetMetricByNameType(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(ts *testSuite)
		wantErr bool
	}{
		{
			name: "success get by name type",
			setup: func(ts *testSuite) {
				ts.mockMetricRepo.EXPECT().
					GetByNameType(ts.ctx, "metric1", common.Gauge).
					Return(common.Metrics{ID: "metric1", MType: common.Gauge}, nil)
			},
			wantErr: false,
		},
		{
			name: "error get by name type",
			setup: func(ts *testSuite) {
				ts.mockMetricRepo.EXPECT().
					GetByNameType(ts.ctx, "metric1", common.Gauge).
					Return(common.Metrics{}, sql.ErrConnDone)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := setupTest(t)

			tt.setup(ts)

			s := &Methods{
				Storage: ts.mockStorage,
			}
			_, err := s.GetMetricByNameType(ts.ctx, "metric1", common.Gauge)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMetricByNameType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestMethods_InsertMetric(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(ts *testSuite)
		wantErr bool
	}{
		{
			name: "success insert metric",
			setup: func(ts *testSuite) {
				ts.mockMetricRepo.EXPECT().
					Insert(ts.ctx, common.Metrics{ID: "metric1", MType: common.Gauge}).
					Return(common.Metrics{ID: "metric1", MType: common.Gauge}, nil)
			},
			wantErr: false,
		},
		{
			name: "error insert metric",
			setup: func(ts *testSuite) {
				ts.mockMetricRepo.EXPECT().
					Insert(ts.ctx, common.Metrics{ID: "metric1", MType: common.Gauge}).
					Return(common.Metrics{}, sql.ErrConnDone)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := setupTest(t)

			tt.setup(ts)

			s := &Methods{
				Storage: ts.mockStorage,
			}
			_, err := s.InsertMetric(ts.ctx, common.Metrics{ID: "metric1", MType: common.Gauge})
			if (err != nil) != tt.wantErr {
				t.Errorf("InsertMetric() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestMethods_InsertMetrics(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(ts *testSuite)
		wantErr bool
	}{
		{
			name: "success insert metrics",
			setup: func(ts *testSuite) {
				ts.mockMetricRepo.EXPECT().
					InsertBatch(ts.ctx, []common.Metrics{{ID: "metric1", MType: common.Gauge}}).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error insert metrics",
			setup: func(ts *testSuite) {
				ts.mockMetricRepo.EXPECT().
					InsertBatch(ts.ctx, []common.Metrics{{ID: "metric1", MType: common.Gauge}}).
					Return(sql.ErrConnDone)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := setupTest(t)

			tt.setup(ts)

			s := &Methods{
				Storage: ts.mockStorage,
			}
			err := s.InsertMetrics(ts.ctx, []common.Metrics{{ID: "metric1", MType: common.Gauge}})
			if (err != nil) != tt.wantErr {
				t.Errorf("InsertBatch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
