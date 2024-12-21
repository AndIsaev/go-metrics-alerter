package metrics

import (
	"fmt"
	"log"
	"math/rand"
	"runtime"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"

	"github.com/AndIsaev/go-metrics-alerter/internal/common"
)

var memStats runtime.MemStats
var pollCount int64

type StorageMetric struct {
	ID    string
	MType string
	Value *float64
	Delta *int64
}

type StorageMetrics struct {
	Metrics map[string]StorageMetric
}

func NewListMetrics() *StorageMetrics {
	return &StorageMetrics{make(map[string]StorageMetric)}
}

func getAddress(f float64) *float64 {
	return &f
}

func (listMetrics *StorageMetrics) Pull() {
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
	listMetrics.Metrics["TotalMemory"] = StorageMetric{ID: "TotalMemory", MType: common.Gauge, Value: &totalMemory}
	listMetrics.Metrics["FreeMemory"] = StorageMetric{ID: "FreeMemory", MType: common.Gauge, Value: &freeMemory}
	for i, utilization := range cpuUtilization {
		name := fmt.Sprintf("CPUutilization%d", i+1)
		listMetrics.Metrics[name] = StorageMetric{ID: name, MType: common.Gauge, Value: &utilization}
	}

	listMetrics.Metrics["Alloc"] = StorageMetric{ID: "Alloc", MType: common.Gauge, Value: getAddress(float64(memStats.Alloc))}
	listMetrics.Metrics["BuckHashSys"] = StorageMetric{ID: "BuckHashSys", MType: common.Gauge, Value: getAddress(float64(memStats.BuckHashSys))}
	listMetrics.Metrics["Frees"] = StorageMetric{ID: "Frees", MType: common.Gauge, Value: getAddress(float64(memStats.Frees))}
	listMetrics.Metrics["GCCPUFraction"] = StorageMetric{ID: "GCCPUFraction", MType: common.Gauge, Value: &memStats.GCCPUFraction}
	listMetrics.Metrics["GCSys"] = StorageMetric{ID: "GCSys", MType: common.Gauge, Value: getAddress(float64(memStats.GCSys))}
	listMetrics.Metrics["HeapAlloc"] = StorageMetric{ID: "HeapAlloc", MType: common.Gauge, Value: getAddress(float64(memStats.HeapAlloc))}
	listMetrics.Metrics["HeapIdle"] = StorageMetric{ID: "HeapIdle", MType: common.Gauge, Value: getAddress(float64(memStats.HeapIdle))}
	listMetrics.Metrics["HeapInuse"] = StorageMetric{ID: "HeapInuse", MType: common.Gauge, Value: getAddress(float64(memStats.HeapInuse))}
	listMetrics.Metrics["HeapObjects"] = StorageMetric{ID: "HeapObjects", MType: common.Gauge, Value: getAddress(float64(memStats.HeapObjects))}
	listMetrics.Metrics["HeapReleased"] = StorageMetric{ID: "HeapReleased", MType: common.Gauge, Value: getAddress(float64(memStats.HeapReleased))}
	listMetrics.Metrics["HeapSys"] = StorageMetric{ID: "HeapSys", MType: common.Gauge, Value: getAddress(float64(memStats.HeapSys))}
	listMetrics.Metrics["LastGC"] = StorageMetric{ID: "LastGC", MType: common.Gauge, Value: getAddress(float64(memStats.LastGC))}
	listMetrics.Metrics["Lookups"] = StorageMetric{ID: "Lookups", MType: common.Gauge, Value: getAddress(float64(memStats.Lookups))}
	listMetrics.Metrics["MCacheInuse"] = StorageMetric{ID: "MCacheInuse", MType: common.Gauge, Value: getAddress(float64(memStats.MCacheInuse))}
	listMetrics.Metrics["MSpanInuse"] = StorageMetric{ID: "MSpanInuse", MType: common.Gauge, Value: getAddress(float64(memStats.MSpanInuse))}
	listMetrics.Metrics["MSpanSys"] = StorageMetric{ID: "MSpanSys", MType: common.Gauge, Value: getAddress(float64(memStats.MSpanSys))}
	listMetrics.Metrics["Mallocs"] = StorageMetric{ID: "Mallocs", MType: common.Gauge, Value: getAddress(float64(memStats.Mallocs))}
	listMetrics.Metrics["NextGC"] = StorageMetric{ID: "NextGC", MType: common.Gauge, Value: getAddress(float64(memStats.NextGC))}
	listMetrics.Metrics["NumForcedGC"] = StorageMetric{ID: "NumForcedGC", MType: common.Gauge, Value: getAddress(float64(memStats.NumForcedGC))}
	listMetrics.Metrics["NumGC"] = StorageMetric{ID: "NumGC", MType: common.Gauge, Value: getAddress(float64(memStats.NumGC))}
	listMetrics.Metrics["OtherSys"] = StorageMetric{ID: "OtherSys", MType: common.Gauge, Value: getAddress(float64(memStats.OtherSys))}
	listMetrics.Metrics["PauseTotalNs"] = StorageMetric{ID: "PauseTotalNs", MType: common.Gauge, Value: getAddress(float64(memStats.PauseTotalNs))}
	listMetrics.Metrics["StackInuse"] = StorageMetric{ID: "StackInuse", MType: common.Gauge, Value: getAddress(float64(memStats.StackInuse))}
	listMetrics.Metrics["StackSys"] = StorageMetric{ID: "StackSys", MType: common.Gauge, Value: getAddress(float64(memStats.StackSys))}
	listMetrics.Metrics["Sys"] = StorageMetric{ID: "Sys", MType: common.Gauge, Value: getAddress(float64(memStats.Sys))}
	listMetrics.Metrics["TotalAlloc"] = StorageMetric{ID: "TotalAlloc", MType: common.Gauge, Value: getAddress(float64(memStats.TotalAlloc))}
	listMetrics.Metrics["MCacheSys"] = StorageMetric{ID: "MCacheSys", MType: common.Gauge, Value: getAddress(float64(memStats.MCacheSys))}
	listMetrics.Metrics["RandomValue"] = StorageMetric{ID: "RandomValue", MType: common.Gauge, Value: getAddress(float64(rand.Int()))}
	listMetrics.Metrics["PollCount"] = StorageMetric{ID: "PollCount", MType: common.Counter, Delta: &pollCount}
}

func (listMetrics *StorageMetrics) AddMetric(metric StorageMetric) {
	listMetrics.Metrics[metric.ID] = metric
}
