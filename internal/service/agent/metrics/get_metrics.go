package metrics

import (
	"github.com/AndIsaev/go-metrics-alerter/internal/common"
	"math/rand"
	"runtime"
	"time"
)

var memStats runtime.MemStats
var pollCount int64

type List []common.Metrics

func float64Convert(i float64) *float64 {
	return &i
}

func GetMetrics(pollInterval time.Duration) List {
	time.Sleep(pollInterval)
	pollCount++

	runtime.ReadMemStats(&memStats)

	metrics := List{
		common.Metrics{ID: "Alloc", MType: common.Gauge, Value: float64Convert(float64(memStats.Alloc))},
		common.Metrics{ID: "BuckHashSys", MType: common.Gauge, Value: float64Convert(float64(memStats.BuckHashSys))},
		common.Metrics{ID: "Frees", MType: common.Gauge, Value: float64Convert(float64(memStats.Frees))},
		common.Metrics{ID: "GCCPUFraction", MType: common.Gauge, Value: &memStats.GCCPUFraction},
		common.Metrics{ID: "GCSys", MType: common.Gauge, Value: float64Convert(float64(memStats.GCSys))},
		common.Metrics{ID: "HeapAlloc", MType: common.Gauge, Value: float64Convert(float64(memStats.HeapAlloc))},
		common.Metrics{ID: "HeapIdle", MType: common.Gauge, Value: float64Convert(float64(memStats.HeapIdle))},
		common.Metrics{ID: "HeapInuse", MType: common.Gauge, Value: float64Convert(float64(memStats.HeapInuse))},
		common.Metrics{ID: "HeapObjects", MType: common.Gauge, Value: float64Convert(float64(memStats.HeapObjects))},
		common.Metrics{ID: "HeapReleased", MType: common.Gauge, Value: float64Convert(float64(memStats.HeapReleased))},
		common.Metrics{ID: "HeapSys", MType: common.Gauge, Value: float64Convert(float64(memStats.HeapSys))},
		common.Metrics{ID: "LastGC", MType: common.Gauge, Value: float64Convert(float64(memStats.LastGC))},
		common.Metrics{ID: "Lookups", MType: common.Gauge, Value: float64Convert(float64(memStats.Lookups))},
		common.Metrics{ID: "MCacheInuse", MType: common.Gauge, Value: float64Convert(float64(memStats.MCacheInuse))},
		common.Metrics{ID: "MSpanInuse", MType: common.Gauge, Value: float64Convert(float64(memStats.MSpanInuse))},
		common.Metrics{ID: "MSpanSys", MType: common.Gauge, Value: float64Convert(float64(memStats.MSpanSys))},
		common.Metrics{ID: "Mallocs", MType: common.Gauge, Value: float64Convert(float64(memStats.Mallocs))},
		common.Metrics{ID: "NextGC", MType: common.Gauge, Value: float64Convert(float64(memStats.NextGC))},
		common.Metrics{ID: "NumForcedGC", MType: common.Gauge, Value: float64Convert(float64(memStats.NumForcedGC))},
		common.Metrics{ID: "NumGC", MType: common.Gauge, Value: float64Convert(float64(memStats.NumGC))},
		common.Metrics{ID: "OtherSys", MType: common.Gauge, Value: float64Convert(float64(memStats.OtherSys))},
		common.Metrics{ID: "PauseTotalNs", MType: common.Gauge, Value: float64Convert(float64(memStats.PauseTotalNs))},
		common.Metrics{ID: "StackInuse", MType: common.Gauge, Value: float64Convert(float64(memStats.StackInuse))},
		common.Metrics{ID: "StackSys", MType: common.Gauge, Value: float64Convert(float64(memStats.StackSys))},
		common.Metrics{ID: "Sys", MType: common.Gauge, Value: float64Convert(float64(memStats.Sys))},
		common.Metrics{ID: "TotalAlloc", MType: common.Gauge, Value: float64Convert(float64(memStats.TotalAlloc))},
		common.Metrics{ID: "StackInuse", MType: common.Gauge, Value: float64Convert(float64(memStats.StackInuse))},
		common.Metrics{ID: "RandomValue", MType: common.Gauge, Value: float64Convert(float64(rand.Int()))},
		common.Metrics{ID: "pollCount", MType: common.Counter, Delta: &pollCount},
	}

	return metrics
}
