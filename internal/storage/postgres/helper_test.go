package postgres

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/AndIsaev/go-metrics-alerter/internal/storage/postgres/mocks"
)

type testSuite struct {
	ctrl           *gomock.Controller
	mockStorage    *mocks.MockStorage
	mockMetricRepo *mocks.MockMetricRepository
	mockSystemRepo *mocks.MockSystemRepository
	ctx            context.Context
}

func setupTest(t *testing.T) *testSuite {
	ctrl := gomock.NewController(t)

	mockStorage := mocks.NewMockStorage(ctrl)
	mockMetricRepo := mocks.NewMockMetricRepository(ctrl)
	mockSystemRepo := mocks.NewMockSystemRepository(ctrl)

	ctx := context.Background()

	mockStorage.EXPECT().Metric().Return(mockMetricRepo).AnyTimes()
	mockStorage.EXPECT().System().Return(mockSystemRepo).AnyTimes()

	return &testSuite{
		ctrl:           ctrl,
		mockStorage:    mockStorage,
		mockMetricRepo: mockMetricRepo,
		mockSystemRepo: mockSystemRepo,
		ctx:            ctx,
	}
}
