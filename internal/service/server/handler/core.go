package handler

import "github.com/AndIsaev/go-metrics-alerter/internal/service/server"

// Handler use like provider service methods for handlers
type Handler struct {
	// MetricService service provide methods of storage
	MetricService server.Service
}
