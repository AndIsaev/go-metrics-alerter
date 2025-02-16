package handler

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/AndIsaev/go-metrics-alerter/internal/service/server/handler/mocks"

	"github.com/AndIsaev/go-metrics-alerter/internal/common"
)

// MockMetricService - это макетный сервис для тестирования
type MockMetricService struct{}

func (m *MockMetricService) InsertMetrics(_ context.Context, _ []common.Metrics) error {
	return nil
}

func (m *MockMetricService) PingStorage(_ context.Context) error {
	return nil
}
func (m *MockMetricService) CloseStorage(_ context.Context) error {
	return nil
}
func (m *MockMetricService) RunMigrationsStorage(_ context.Context) error {
	return nil
}

func (m *MockMetricService) ListMetrics(_ context.Context) ([]common.Metrics, error) {
	return []common.Metrics{
		{ID: "metric1", MType: common.Counter, Delta: common.LinkInt64(1)},
		{ID: "metric2", MType: common.Gauge, Value: common.LinkFloat64(10.4)},
	}, nil
}

func (m *MockMetricService) UpdateMetricByValue(_ context.Context, _ common.Metrics, _ any) error {
	return nil
}

func (m *MockMetricService) GetMetricByName(_ context.Context, _ string) (common.Metrics, error) {
	return common.Metrics{ID: "metric1", MType: common.Counter, Delta: common.LinkInt64(23)}, nil
}

func (m *MockMetricService) GetMetricByNameType(_ context.Context, _ string, _ string) (common.Metrics, error) {
	return common.Metrics{}, nil
}

func (m *MockMetricService) InsertMetric(_ context.Context, _ common.Metrics) (common.Metrics, error) {
	return common.Metrics{
		ID:    "metric1",
		MType: common.Counter,
		Delta: common.LinkInt64(123),
	}, nil
}

type testSuite struct {
	ctrl        *gomock.Controller
	mockService *mocks.MockService
	ctx         context.Context
}

func setupTest(t *testing.T) *testSuite {
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	mockService := mocks.NewMockService(ctrl)
	return &testSuite{
		ctrl:        ctrl,
		mockService: mockService,
		ctx:         ctx,
	}
}
