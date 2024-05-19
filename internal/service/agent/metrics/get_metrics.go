package metrics

import (
	"github.com/AndIsaev/go-metrics-alerter/internal/common"
	"math/rand"
	"runtime"
	"time"
)

var memStats runtime.MemStats
var pollCount int

type Metrics []struct {
	Name       string
	Value      interface{}
	MetricType string
}

func GetMetrics(pollInterval time.Duration) Metrics {
	time.Sleep(pollInterval)
	pollCount++

	runtime.ReadMemStats(&memStats)

	metrics := Metrics{
		{Name: "Alloc", Value: memStats.Alloc, MetricType: common.Gauge},
		{Name: "BuckHashSys", Value: memStats.BuckHashSys, MetricType: common.Gauge},
		{Name: "Frees", Value: memStats.Frees, MetricType: common.Gauge},
		{Name: "GCCPUFraction", Value: memStats.GCCPUFraction, MetricType: common.Gauge},
		{Name: "GCSys", Value: memStats.GCSys, MetricType: common.Gauge},
		{Name: "HeapAlloc", Value: memStats.HeapAlloc, MetricType: common.Gauge},
		{Name: "HeapIdle", Value: memStats.HeapIdle, MetricType: common.Gauge},
		{Name: "HeapInuse", Value: memStats.HeapInuse, MetricType: common.Gauge},
		{Name: "HeapObjects", Value: memStats.HeapObjects, MetricType: common.Gauge},
		{Name: "HeapReleased", Value: memStats.HeapReleased, MetricType: common.Gauge},
		{Name: "HeapSys", Value: memStats.HeapSys, MetricType: common.Gauge},
		{Name: "LastGC", Value: memStats.LastGC, MetricType: common.Gauge},
		{Name: "Lookups", Value: memStats.Lookups, MetricType: common.Gauge},
		{Name: "MCacheInuse", Value: memStats.MCacheInuse, MetricType: common.Gauge},
		{Name: "MCacheSys", Value: memStats.MCacheSys, MetricType: common.Gauge},
		{Name: "MSpanInuse", Value: memStats.MSpanInuse, MetricType: common.Gauge},
		{Name: "MSpanSys", Value: memStats.MSpanSys, MetricType: common.Gauge},
		{Name: "Mallocs", Value: memStats.Mallocs, MetricType: common.Gauge},
		{Name: "NextGC", Value: memStats.NextGC, MetricType: common.Gauge},
		{Name: "NumForcedGC", Value: memStats.NumForcedGC, MetricType: common.Gauge},
		{Name: "NumGC", Value: memStats.NumGC, MetricType: common.Gauge},
		{Name: "OtherSys", Value: memStats.OtherSys, MetricType: common.Gauge},
		{Name: "PauseTotalNs", Value: memStats.PauseTotalNs, MetricType: common.Gauge},
		{Name: "StackInuse", Value: memStats.StackInuse, MetricType: common.Gauge},
		{Name: "StackSys", Value: memStats.StackSys, MetricType: common.Gauge},
		{Name: "Sys", Value: memStats.Sys, MetricType: common.Gauge},
		{Name: "TotalAlloc", Value: memStats.TotalAlloc, MetricType: common.Gauge},
		{Name: "pollCount", Value: pollCount, MetricType: common.Counter},
		{Name: "RandomValue", Value: rand.Int(), MetricType: common.Counter},
	}

	return metrics
}
