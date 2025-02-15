package rpc

import (
	"context"
	"log"

	"github.com/AndIsaev/go-metrics-alerter/internal/common"
	pb "github.com/AndIsaev/go-metrics-alerter/internal/service/proto"
)

func (c *Client) SendMetrics(metrics []common.Metrics) error {
	ctx := context.Background()
	var delta int64
	var value float64
	var body []*pb.Metric

	for _, m := range metrics {
		if m.Delta != nil {
			delta = *m.Delta
		}
		if m.Value != nil {
			value = *m.Value
		}
		body = append(body, &pb.Metric{Id: m.ID, Type: m.MType, Delta: delta, Value: value})
	}
	request := &pb.InsertBatchRequest{
		Metrics: body,
		Amount:  int32(len(body)),
	}

	resp, err := c.GRPCClient.InsertBatch(ctx, request)
	if err != nil {
		log.Printf("error insert batch by grpc client: %v\n", err)
		return err
	}
	log.Printf("response: %v\n", resp)

	return nil
}
