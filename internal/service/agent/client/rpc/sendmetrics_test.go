package rpc

import (
	"context"
	"errors"
	"testing"

	"google.golang.org/grpc/status"

	"github.com/AndIsaev/go-metrics-alerter/internal/common"
	pb "github.com/AndIsaev/go-metrics-alerter/internal/service/proto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSendMetrics_Success(t *testing.T) {
	ctx := context.Background()
	metrics := []common.Metrics{
		{ID: "metric1", MType: "gauge", Value: ptrFloat64(123.456)},
		{ID: "metric2", MType: "counter", Delta: ptrInt64(10)},
	}

	mockClient := new(MockMetricServiceClient)
	mockClient.On("InsertBatch", ctx, mock.Anything).Return(&pb.InsertBatchResponse{}, nil)

	grpcClient := &GRPCClient{MetricServiceClient: mockClient}
	client := NewClient(grpcClient)

	err := client.SendMetrics(ctx, metrics)

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestSendMetrics_InsertBatchError(t *testing.T) {
	ctx := context.Background()
	metrics := []common.Metrics{
		{ID: "metric1", MType: "gauge", Value: ptrFloat64(123.456)},
		{ID: "metric2", MType: "counter", Delta: ptrInt64(10)},
	}

	mockClient := new(MockMetricServiceClient)
	mockClient.On("InsertBatch", ctx, mock.Anything).Return((*pb.InsertBatchResponse)(nil), status.Errorf(status.Code(errors.New("error")), "test insert batch error"))

	grpcClient := &GRPCClient{MetricServiceClient: mockClient}
	client := NewClient(grpcClient)

	err := client.SendMetrics(ctx, metrics)

	assert.Error(t, err)
	mockClient.AssertExpectations(t)
}

func TestSendMetrics_ContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	metrics := []common.Metrics{
		{ID: "metric1", MType: "gauge", Value: ptrFloat64(123.456)},
		{ID: "metric2", MType: "counter", Delta: ptrInt64(10)},
	}

	mockClient := new(MockMetricServiceClient)
	grpcClient := &GRPCClient{MetricServiceClient: mockClient}
	client := NewClient(grpcClient)

	err := client.SendMetrics(ctx, metrics)

	assert.NoError(t, err)
	mockClient.AssertNotCalled(t, "InsertBatch", ctx, mock.Anything)
}

func ptrInt64(v int64) *int64 {
	return &v
}

func ptrFloat64(v float64) *float64 {
	return &v
}
