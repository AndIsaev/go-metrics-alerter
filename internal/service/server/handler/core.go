package handler

import "github.com/AndIsaev/go-metrics-alerter/internal/service/server"

type Handler struct {
	MetricService server.Service
}
