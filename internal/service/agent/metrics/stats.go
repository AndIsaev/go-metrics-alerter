package metrics

import (
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"sync"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"

	"github.com/AndIsaev/go-metrics-alerter/internal/common"
)

var memStats runtime.MemStats
var pollCount int64

// StorageMetrics stores map of metrics
type StorageMetrics struct {
	Metrics map[string]common.Metrics
	mu      sync.Mutex
}

// NewListMetrics init storage metrics
func NewListMetrics() *StorageMetrics {
	return &StorageMetrics{Metrics: make(map[string]common.Metrics)}
}

// Pull get metrics
func (sm *StorageMetrics) Pull() {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		log.Printf("error getting memory stats: %v", err)
	}
	totalMemory := float64(vmStat.Total)
	freeMemory := float64(vmStat.Free)

	cpuUtilization, err := cpu.Percent(0, true)
	if err != nil {
		log.Printf("error getting CPU stats: %v", err)
	}

	pollCount++
	runtime.ReadMemStats(&memStats)
	sm.Metrics["TotalMemory"] = common.Metrics{ID: "TotalMemory", MType: common.Gauge, Value: &totalMemory}
	sm.Metrics["FreeMemory"] = common.Metrics{ID: "FreeMemory", MType: common.Gauge, Value: &freeMemory}
	for i, utilization := range cpuUtilization {
		name := fmt.Sprintf("CPUutilization%d", i+1)
		sm.Metrics[name] = common.Metrics{ID: name, MType: common.Gauge, Value: &utilization}
	}

	sm.Metrics["Alloc"] = common.Metrics{ID: "Alloc", MType: common.Gauge, Value: common.LinkFloat64(float64(memStats.Alloc))}
	sm.Metrics["BuckHashSys"] = common.Metrics{ID: "BuckHashSys", MType: common.Gauge, Value: common.LinkFloat64(float64(memStats.BuckHashSys))}
	sm.Metrics["Frees"] = common.Metrics{ID: "Frees", MType: common.Gauge, Value: common.LinkFloat64(float64(memStats.Frees))}
	sm.Metrics["GCCPUFraction"] = common.Metrics{ID: "GCCPUFraction", MType: common.Gauge, Value: &memStats.GCCPUFraction}
	sm.Metrics["GCSys"] = common.Metrics{ID: "GCSys", MType: common.Gauge, Value: common.LinkFloat64(float64(memStats.GCSys))}
	sm.Metrics["HeapAlloc"] = common.Metrics{ID: "HeapAlloc", MType: common.Gauge, Value: common.LinkFloat64(float64(memStats.HeapAlloc))}
	sm.Metrics["HeapIdle"] = common.Metrics{ID: "HeapIdle", MType: common.Gauge, Value: common.LinkFloat64(float64(memStats.HeapIdle))}
	sm.Metrics["HeapInuse"] = common.Metrics{ID: "HeapInuse", MType: common.Gauge, Value: common.LinkFloat64(float64(memStats.HeapInuse))}
	sm.Metrics["HeapObjects"] = common.Metrics{ID: "HeapObjects", MType: common.Gauge, Value: common.LinkFloat64(float64(memStats.HeapObjects))}
	sm.Metrics["HeapReleased"] = common.Metrics{ID: "HeapReleased", MType: common.Gauge, Value: common.LinkFloat64(float64(memStats.HeapReleased))}
	sm.Metrics["HeapSys"] = common.Metrics{ID: "HeapSys", MType: common.Gauge, Value: common.LinkFloat64(float64(memStats.HeapSys))}
	sm.Metrics["LastGC"] = common.Metrics{ID: "LastGC", MType: common.Gauge, Value: common.LinkFloat64(float64(memStats.LastGC))}
	sm.Metrics["Lookups"] = common.Metrics{ID: "Lookups", MType: common.Gauge, Value: common.LinkFloat64(float64(memStats.Lookups))}
	sm.Metrics["MCacheInuse"] = common.Metrics{ID: "MCacheInuse", MType: common.Gauge, Value: common.LinkFloat64(float64(memStats.MCacheInuse))}
	sm.Metrics["MSpanInuse"] = common.Metrics{ID: "MSpanInuse", MType: common.Gauge, Value: common.LinkFloat64(float64(memStats.MSpanInuse))}
	sm.Metrics["MSpanSys"] = common.Metrics{ID: "MSpanSys", MType: common.Gauge, Value: common.LinkFloat64(float64(memStats.MSpanSys))}
	sm.Metrics["Mallocs"] = common.Metrics{ID: "Mallocs", MType: common.Gauge, Value: common.LinkFloat64(float64(memStats.Mallocs))}
	sm.Metrics["NextGC"] = common.Metrics{ID: "NextGC", MType: common.Gauge, Value: common.LinkFloat64(float64(memStats.NextGC))}
	sm.Metrics["NumForcedGC"] = common.Metrics{ID: "NumForcedGC", MType: common.Gauge, Value: common.LinkFloat64(float64(memStats.NumForcedGC))}
	sm.Metrics["NumGC"] = common.Metrics{ID: "NumGC", MType: common.Gauge, Value: common.LinkFloat64(float64(memStats.NumGC))}
	sm.Metrics["OtherSys"] = common.Metrics{ID: "OtherSys", MType: common.Gauge, Value: common.LinkFloat64(float64(memStats.OtherSys))}
	sm.Metrics["PauseTotalNs"] = common.Metrics{ID: "PauseTotalNs", MType: common.Gauge, Value: common.LinkFloat64(float64(memStats.PauseTotalNs))}
	sm.Metrics["StackInuse"] = common.Metrics{ID: "StackInuse", MType: common.Gauge, Value: common.LinkFloat64(float64(memStats.StackInuse))}
	sm.Metrics["StackSys"] = common.Metrics{ID: "StackSys", MType: common.Gauge, Value: common.LinkFloat64(float64(memStats.StackSys))}
	sm.Metrics["Sys"] = common.Metrics{ID: "Sys", MType: common.Gauge, Value: common.LinkFloat64(float64(memStats.Sys))}
	sm.Metrics["TotalAlloc"] = common.Metrics{ID: "TotalAlloc", MType: common.Gauge, Value: common.LinkFloat64(float64(memStats.TotalAlloc))}
	sm.Metrics["MCacheSys"] = common.Metrics{ID: "MCacheSys", MType: common.Gauge, Value: common.LinkFloat64(float64(memStats.MCacheSys))}
	sm.Metrics["RandomValue"] = common.Metrics{ID: "RandomValue", MType: common.Gauge, Value: common.LinkFloat64(float64(rand.Int()))}
	sm.Metrics["PollCount"] = common.Metrics{ID: "PollCount", MType: common.Counter, Delta: &pollCount}
}
