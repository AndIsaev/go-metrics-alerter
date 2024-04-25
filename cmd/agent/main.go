package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"time"
)

const (
	address                      = "http://localhost:8080/update/%v/%v/%v"
	gauge          string        = "gauge"
	counter        string        = "counter"
	reportInterval time.Duration = 10
	pollInterval   time.Duration = 2
)

var memStats runtime.MemStats
var PollCount int

type Metrics []struct {
	name       string
	value      interface{}
	metricType string
}

func getMetrics() Metrics {
	time.Sleep(pollInterval * time.Second)
	PollCount++

	runtime.ReadMemStats(&memStats)

	metrics := Metrics{
		{name: "Alloc", value: memStats.Alloc, metricType: gauge},
		{name: "BuckHashSys", value: memStats.BuckHashSys, metricType: gauge},
		{name: "Frees", value: memStats.Frees, metricType: gauge},
		{name: "GCCPUFraction", value: memStats.GCCPUFraction, metricType: gauge},
		{name: "GCSys", value: memStats.GCSys, metricType: gauge},
		{name: "HeapAlloc", value: memStats.HeapAlloc, metricType: gauge},
		{name: "HeapIdle", value: memStats.HeapIdle, metricType: gauge},
		{name: "HeapInuse", value: memStats.HeapInuse, metricType: gauge},
		{name: "HeapObjects", value: memStats.HeapObjects, metricType: gauge},
		{name: "HeapReleased", value: memStats.HeapReleased, metricType: gauge},
		{name: "HeapSys", value: memStats.HeapSys, metricType: gauge},
		{name: "LastGC", value: memStats.LastGC, metricType: gauge},
		{name: "Lookups", value: memStats.Lookups, metricType: gauge},
		{name: "MCacheInuse", value: memStats.MCacheInuse, metricType: gauge},
		{name: "MCacheSys", value: memStats.MCacheSys, metricType: gauge},
		{name: "MSpanInuse", value: memStats.MSpanInuse, metricType: gauge},
		{name: "MSpanSys", value: memStats.MSpanSys, metricType: gauge},
		{name: "Mallocs", value: memStats.Mallocs, metricType: gauge},
		{name: "NextGC", value: memStats.NextGC, metricType: gauge},
		{name: "NumForcedGC", value: memStats.NumForcedGC, metricType: gauge},
		{name: "NumGC", value: memStats.NumGC, metricType: gauge},
		{name: "OtherSys", value: memStats.OtherSys, metricType: gauge},
		{name: "PauseTotalNs", value: memStats.PauseTotalNs, metricType: gauge},
		{name: "StackInuse", value: memStats.StackInuse, metricType: gauge},
		{name: "StackSys", value: memStats.StackSys, metricType: gauge},
		{name: "Sys", value: memStats.Sys, metricType: gauge},
		{name: "TotalAlloc", value: memStats.TotalAlloc, metricType: gauge},
		{name: "PollCount", value: PollCount, metricType: counter},
		{name: "RandomValue", value: rand.Int(), metricType: counter},
	}

	return metrics
}

//func writeToChannel(c chan<- Metrics, metrics Metrics) {
//	c <- metrics
//}
//
//func readFromChannel(metrics <-chan Metrics) Metrics {
//	m := <-metrics
//	return m
//}

func postRequest(url, contentType string, body []byte) {

	resp, err := http.Post(url, contentType, bytes.NewBuffer(body))
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
		return
	}

	defer resp.Body.Close()
}

func sendReport(m Metrics) {
	time.Sleep(reportInterval * time.Second)

	for _, v := range m {
		url := fmt.Sprintf(address, v.metricType, v.name, v.value)
		postRequest(url, "text/plain", nil)
	}
}

func main() {

	newMetrics := getMetrics()
	sendReport(newMetrics)
	//ch := make(chan Metrics)

	//metricsTick := time.NewTicker(pollInterval * time.Second)
	//reportTick := time.NewTicker(reportInterval * time.Second)
	//
	//for {
	//	select {
	//	case <-metricsTick.C:
	//
	//		newMetrics := getMetrics()
	//		go writeToChannel(ch, newMetrics)
	//	case <-reportTick.C:
	//		v := readFromChannel(ch)
	//		go sendReport(v)
	//	}
	//
	//}

}
