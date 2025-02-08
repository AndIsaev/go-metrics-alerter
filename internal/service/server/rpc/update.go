package rpc

import (
	"context"
	"log"

	"google.golang.org/grpc/codes"

	"github.com/AndIsaev/go-metrics-alerter/internal/common"
	pb "github.com/AndIsaev/go-metrics-alerter/internal/service/proto"
)

func (p *MetricServiceServer) InsertBatch(ctx context.Context, req *pb.InsertBatchRequest) (*pb.InsertBatchResponse, error) {
	input := req.GetMetrics()
	metrics := make([]common.Metrics, req.Amount)
	for _, val := range input {
		metrics = append(
			metrics,
			common.Metrics{
				ID:    val.Id,
				MType: val.Type,
				Value: linkFloat64(val.Value),
				Delta: linkInt64(val.Delta)},
		)
	}
	err := p.Storage.Metric().InsertBatch(ctx, metrics)
	if err != nil {
		log.Printf("error update metrics: %v", err)
		return nil, err
	}
	return &pb.InsertBatchResponse{StatusCode: uint32(codes.OK), Message: "success"}, nil
}
