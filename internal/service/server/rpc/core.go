package rpc

import (
	pb "github.com/AndIsaev/go-metrics-alerter/internal/service/proto"
	"github.com/AndIsaev/go-metrics-alerter/internal/storage"
)

type MetricServiceServer struct {
	pb.MetricServiceServer
	Storage storage.Storage
}

func linkFloat64(num float64) *float64 {
	return &num
}
func linkInt64(num int64) *int64 {
	return &num
}
