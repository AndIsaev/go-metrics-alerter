package rpc

import (
	"context"
	"log"

	"google.golang.org/grpc"

	pb "github.com/AndIsaev/go-metrics-alerter/internal/service/proto"
)

// Client abstract structure over GRPCClient
type Client struct {
	GRPCClient *GRPCClient
}

func NewClient(grpcClient *GRPCClient) *Client {
	return &Client{GRPCClient: grpcClient}
}

type GRPCClient struct {
	pb.MetricServiceClient
}

func NewGRPCClient(conn *grpc.ClientConn) *GRPCClient {
	return &GRPCClient{pb.NewMetricServiceClient(conn)}
}

func (c *GRPCClient) InsertBatch(ctx context.Context, in *pb.InsertBatchRequest, _ ...grpc.CallOption) (*pb.InsertBatchResponse, error) {
	batch, err := c.MetricServiceClient.InsertBatch(ctx, in)
	if err != nil {
		log.Printf("error request in rpc method InsertBatch: %v\n", err)
		return nil, err
	}
	return batch, err
}
