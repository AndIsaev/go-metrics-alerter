package metrics

import (
	"github.com/AndIsaev/go-metrics-alerter/internal/common"
	"math/rand"
	"runtime"
)

var memStats runtime.MemStats
var pollCount int64

type StorageMetric struct {
	ID    string
	MType string
	Value float64
	Delta int64
}

type StorageMetrics struct {
	Metrics map[string]StorageMetric
}

func NewListMetrics() *StorageMetrics {
	return &StorageMetrics{make(map[string]StorageMetric)}
}

func (LM *StorageMetrics) Pull() {

	pollCount++
	runtime.ReadMemStats(&memStats)

	LM.Metrics["Alloc"] = StorageMetric{ID: "Alloc", MType: common.Gauge, Value: float64(memStats.Alloc)}
	LM.Metrics["BuckHashSys"] = StorageMetric{ID: "BuckHashSys", MType: common.Gauge, Value: float64(memStats.BuckHashSys)}
	LM.Metrics["Frees"] = StorageMetric{ID: "Frees", MType: common.Gauge, Value: float64(memStats.Frees)}
	LM.Metrics["GCCPUFraction"] = StorageMetric{ID: "GCCPUFraction", MType: common.Gauge, Value: memStats.GCCPUFraction}
	LM.Metrics["GCSys"] = StorageMetric{ID: "GCSys", MType: common.Gauge, Value: float64(memStats.GCSys)}
	LM.Metrics["HeapAlloc"] = StorageMetric{ID: "HeapAlloc", MType: common.Gauge, Value: float64(memStats.HeapAlloc)}
	LM.Metrics["HeapIdle"] = StorageMetric{ID: "HeapIdle", MType: common.Gauge, Value: float64(memStats.HeapIdle)}
	LM.Metrics["HeapInuse"] = StorageMetric{ID: "HeapInuse", MType: common.Gauge, Value: float64(memStats.HeapInuse)}
	LM.Metrics["HeapObjects"] = StorageMetric{ID: "HeapObjects", MType: common.Gauge, Value: float64(memStats.HeapObjects)}
	LM.Metrics["HeapReleased"] = StorageMetric{ID: "HeapReleased", MType: common.Gauge, Value: float64(memStats.HeapReleased)}
	LM.Metrics["HeapSys"] = StorageMetric{ID: "HeapSys", MType: common.Gauge, Value: float64(memStats.HeapSys)}
	LM.Metrics["LastGC"] = StorageMetric{ID: "LastGC", MType: common.Gauge, Value: float64(memStats.LastGC)}
	LM.Metrics["Lookups"] = StorageMetric{ID: "Lookups", MType: common.Gauge, Value: float64(memStats.Lookups)}
	LM.Metrics["MCacheInuse"] = StorageMetric{ID: "MCacheInuse", MType: common.Gauge, Value: float64(memStats.MCacheInuse)}
	LM.Metrics["MSpanInuse"] = StorageMetric{ID: "MSpanInuse", MType: common.Gauge, Value: float64(memStats.MSpanInuse)}
	LM.Metrics["MSpanSys"] = StorageMetric{ID: "MSpanSys", MType: common.Gauge, Value: float64(memStats.MSpanSys)}
	LM.Metrics["Mallocs"] = StorageMetric{ID: "Mallocs", MType: common.Gauge, Value: float64(memStats.Mallocs)}
	LM.Metrics["NextGC"] = StorageMetric{ID: "NextGC", MType: common.Gauge, Value: float64(memStats.NextGC)}
	LM.Metrics["NumForcedGC"] = StorageMetric{ID: "NumForcedGC", MType: common.Gauge, Value: float64(memStats.NumForcedGC)}
	LM.Metrics["NumGC"] = StorageMetric{ID: "NumGC", MType: common.Gauge, Value: float64(memStats.NumGC)}
	LM.Metrics["OtherSys"] = StorageMetric{ID: "OtherSys", MType: common.Gauge, Value: float64(memStats.OtherSys)}
	LM.Metrics["PauseTotalNs"] = StorageMetric{ID: "PauseTotalNs", MType: common.Gauge, Value: float64(memStats.PauseTotalNs)}
	LM.Metrics["StackInuse"] = StorageMetric{ID: "StackInuse", MType: common.Gauge, Value: float64(memStats.StackInuse)}
	LM.Metrics["StackSys"] = StorageMetric{ID: "StackSys", MType: common.Gauge, Value: float64(memStats.StackSys)}
	LM.Metrics["Sys"] = StorageMetric{ID: "Sys", MType: common.Gauge, Value: float64(memStats.Sys)}
	LM.Metrics["TotalAlloc"] = StorageMetric{ID: "TotalAlloc", MType: common.Gauge, Value: float64(memStats.TotalAlloc)}
	LM.Metrics["MCacheSys"] = StorageMetric{ID: "MCacheSys", MType: common.Gauge, Value: float64(memStats.MCacheSys)}
	LM.Metrics["RandomValue"] = StorageMetric{ID: "RandomValue", MType: common.Gauge, Value: float64(rand.Int())}
	LM.Metrics["PollCount"] = StorageMetric{ID: "PollCount", MType: common.Counter, Delta: pollCount}

}
