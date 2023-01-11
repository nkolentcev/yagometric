package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

const (
	pollInterval   = time.Duration(2) * time.Second
	reportInterval = time.Duration(10) * time.Second
	host           = "127.0.0.1"
	port           = "8080"
)

type gauge float64
type counter int64

func main() {

	fmt.Println("start")
	ctx, cancel := context.WithCancel(context.Background())
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT)
	go readMetrics(ctx)

	<-sig
	cancel()
}

func sendMetrics(ctx context.Context, uri string) {

	client := http.Client{}
	request, err := http.NewRequest(http.MethodPost, uri, nil)
	if err != nil {
		fmt.Println(err)
	}
	request.Header.Add("Content-Type", "text/plain")

	response, err := client.Do(request)
	if err != nil {
		log.Println(err)
		response.Body.Close()
	}
}

func updateGaugeMetric(ctx context.Context, metricName string, metricValue gauge) {
	uri := fmt.Sprintf("http://%s:%s/update/%s/%s/%f", host, port, "gauge", metricName, metricValue)
	sendMetrics(ctx, uri)
}

func updateCounterMetric(ctx context.Context, metricName string, metricValue counter) {
	uri := fmt.Sprintf("http://%s:%s/update/%s/%s/%v", host, port, "counter", metricName, metricValue)
	sendMetrics(ctx, uri)
}

func readMetrics(ctx context.Context) {
	var rtm runtime.MemStats
	count := 0
	for {

		<-time.After(pollInterval)

		runtime.ReadMemStats(&rtm)

		updateGaugeMetric(ctx, "Alloc", gauge(rtm.Alloc))
		updateGaugeMetric(ctx, "BuckHashSys", gauge(rtm.BuckHashSys))
		updateGaugeMetric(ctx, "Frees", gauge(rtm.Frees))
		updateGaugeMetric(ctx, "GCCPUFraction", gauge(rtm.GCCPUFraction))
		updateGaugeMetric(ctx, "GCSys", gauge(rtm.GCSys))
		updateGaugeMetric(ctx, "HeapAlloc", gauge(rtm.HeapAlloc))
		updateGaugeMetric(ctx, "HeapIdle", gauge(rtm.HeapInuse))
		updateGaugeMetric(ctx, "HeapInuse", gauge(rtm.HeapInuse))
		updateGaugeMetric(ctx, "HeapObjects", gauge(rtm.HeapObjects))
		updateGaugeMetric(ctx, "HeapReleased", gauge(rtm.HeapReleased))
		updateGaugeMetric(ctx, "HeapSys", gauge(rtm.HeapSys))
		updateGaugeMetric(ctx, "LastGC", gauge(rtm.LastGC))
		updateGaugeMetric(ctx, "Lookups", gauge(rtm.Lookups))
		updateGaugeMetric(ctx, "MCacheInuse", gauge(rtm.MCacheInuse))
		updateGaugeMetric(ctx, "MCacheSys", gauge(rtm.MCacheSys))
		updateGaugeMetric(ctx, "MSpanInuse", gauge(rtm.MSpanInuse))
		updateGaugeMetric(ctx, "MSpanSys", gauge(rtm.MSpanSys))
		updateGaugeMetric(ctx, "Mallocs", gauge(rtm.Mallocs))
		updateGaugeMetric(ctx, "NextGC", gauge(rtm.NextGC))
		updateGaugeMetric(ctx, "NumForcedGC", gauge(rtm.NumForcedGC))
		updateGaugeMetric(ctx, "NumGC", gauge(rtm.NumGC))
		updateGaugeMetric(ctx, "OtherSys", gauge(rtm.OtherSys))
		updateGaugeMetric(ctx, "PauseTotalNs", gauge(rtm.PauseTotalNs))
		updateGaugeMetric(ctx, "StackInuse", gauge(rtm.StackInuse))
		updateGaugeMetric(ctx, "StackSys", gauge(rtm.StackSys))
		updateGaugeMetric(ctx, "Sys", gauge(rtm.Sys))
		updateGaugeMetric(ctx, "TotalAlloc", gauge(rtm.TotalAlloc))
		updateGaugeMetric(ctx, "RandomValue", gauge(rand.Float64()))

		count++
		updateCounterMetric(ctx, "PollCount", counter(count))

		log.Printf("metrics updated %v", count)
	}
}
