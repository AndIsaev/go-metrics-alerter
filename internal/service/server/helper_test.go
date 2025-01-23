package server

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/AndIsaev/go-metrics-alerter/internal/storage/mocks"
)

type testSuite struct {
	ctrl           *gomock.Controller
	mockStorage    *mocks.MockStorage
	mockSystemRepo *mocks.MockSystemRepository
	mockMetricRepo *mocks.MockMetricRepository
	ctx            context.Context
}

// setupTest - use it to prepare the context for testing service.Service
func setupTest(t *testing.T) *testSuite {
	ctrl := gomock.NewController(t)

	mockStorage := mocks.NewMockStorage(ctrl)
	mockSystemRepo := mocks.NewMockSystemRepository(ctrl)
	mockMetricRepo := mocks.NewMockMetricRepository(ctrl)

	ctx := context.Background()

	mockStorage.EXPECT().System().Return(mockSystemRepo).AnyTimes()
	mockStorage.EXPECT().Metric().Return(mockMetricRepo).AnyTimes()

	return &testSuite{
		ctrl:           ctrl,
		mockStorage:    mockStorage,
		mockSystemRepo: mockSystemRepo,
		mockMetricRepo: mockMetricRepo,
		ctx:            ctx,
	}
}
