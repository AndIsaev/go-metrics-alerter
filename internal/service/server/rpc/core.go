package rpc

import (
	pb "github.com/AndIsaev/go-metrics-alerter/internal/service/proto"
	"github.com/AndIsaev/go-metrics-alerter/internal/storage"
)

type MetricServiceServer struct {
	pb.MetricServiceServer
	Storage storage.Storage
}
