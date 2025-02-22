package rpc

import (
	"context"
	"testing"

	pb "github.com/AndIsaev/go-metrics-alerter/internal/service/proto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

// MockMetricServiceClient is a mock implementation of the MetricServiceClient interface.
type MockMetricServiceClient struct {
	mock.Mock
}

func (m *MockMetricServiceClient) InsertBatch(ctx context.Context, in *pb.InsertBatchRequest, _ ...grpc.CallOption) (*pb.InsertBatchResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*pb.InsertBatchResponse), args.Error(1)
}

func TestInsertBatch_Success(t *testing.T) {
	mockClient := new(MockMetricServiceClient)
	response := &pb.InsertBatchResponse{}

	mockClient.On("InsertBatch", mock.Anything, mock.Anything).Return(response, nil)

	grpcClient := &GRPCClient{MetricServiceClient: mockClient}
	client := NewClient(grpcClient)

	req := &pb.InsertBatchRequest{}
	res, err := client.GRPCClient.InsertBatch(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, response, res)

	mockClient.AssertExpectations(t)
}
