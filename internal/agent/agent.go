package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"time"

	"github.com/nkolentcev/yagometric/internal/config"
)

type gauge float64
type counter int64

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

type met struct {
	Alloc         gauge
	BuckHashSys   gauge
	Frees         gauge
	GCCPUFraction gauge
	GCSys         gauge
	HeapAlloc     gauge
	HeapIdle      gauge
	HeapInuse     gauge
	HeapObjects   gauge
	HeapReleased  gauge
	HeapSys       gauge
	LastGC        gauge
	Lookups       gauge
	MCacheInuse   gauge
	MCacheSys     gauge
	MSpanInuse    gauge
	MSpanSys      gauge
	Mallocs       gauge
	NextGC        gauge
	NumForcedGC   gauge
	NumGC         gauge
	OtherSys      gauge
	PauseTotalNs  gauge
	StackInuse    gauge
	StackSys      gauge
	Sys           gauge
	TotalAlloc    gauge
	RandomValue   gauge
	PollCount     counter
}
type contextKey int

type Agent struct {
	pollInterval   time.Duration
	reportInterval time.Duration
	reportAddres   string
}

const (
	Address contextKey = iota
)

func NewAgent(cfg *config.AgentCfg) *Agent {
	var agent Agent
	agent.pollInterval = cfg.PollInterval
	agent.reportInterval = cfg.ReportInterval
	agent.reportAddres = cfg.Address
	return &agent
}

func (a *Agent) Start(ctx context.Context) {

	ctx = context.WithValue(ctx, Address, a.reportAddres)

	mem := new(met)
	go readMetrics(ctx, a.pollInterval, mem)
	go updateMetrics(ctx, a.reportInterval, mem)
}

func updateMetric(ctx context.Context, metricType string, metricName string, metricValue gauge) {

	host := ctx.Value(Address)

	metrics := new(Metrics)
	metrics.ID = metricName
	metrics.MType = metricType

	switch metricType {
	case "gauge":
		ftype := float64(metricValue)
		metrics.Value = &ftype
	case "counter":
		ftype := int64(metricValue)
		metrics.Delta = &ftype
	}

	dataJSON, err := json.Marshal(metrics)
	if err != nil {
		log.Panicf("unable convert metric in json %s", err)
	}
	uri := fmt.Sprintf("http://%s/update/", host)

	client := http.Client{}
	request, err := http.NewRequest(http.MethodPost, uri, bytes.NewReader(dataJSON))
	if err != nil {
		fmt.Println(err)
	}
	request.Header.Add("Content-Type", "tapplication/json")

	response, err := client.Do(request)
	if err != nil {
		log.Println(err)
		response.Body.Close()
	}
}

func updateMetrics(ctx context.Context, reportInterval time.Duration, mem *met) {

	for {
		<-time.After(reportInterval)

		updateMetric(ctx, "gauge", "Alloc", gauge(mem.Alloc))
		updateMetric(ctx, "gauge", "BuckHashSys", gauge(mem.BuckHashSys))
		updateMetric(ctx, "gauge", "Frees", gauge(mem.Frees))
		updateMetric(ctx, "gauge", "GCCPUFraction", gauge(mem.GCCPUFraction))
		updateMetric(ctx, "gauge", "GCSys", gauge(mem.GCSys))
		updateMetric(ctx, "gauge", "HeapAlloc", gauge(mem.HeapAlloc))
		updateMetric(ctx, "gauge", "HeapIdle", gauge(mem.HeapInuse))
		updateMetric(ctx, "gauge", "HeapInuse", gauge(mem.HeapInuse))
		updateMetric(ctx, "gauge", "HeapObjects", gauge(mem.HeapObjects))
		updateMetric(ctx, "gauge", "HeapReleased", gauge(mem.HeapReleased))
		updateMetric(ctx, "gauge", "HeapSys", gauge(mem.HeapSys))
		updateMetric(ctx, "gauge", "LastGC", gauge(mem.LastGC))
		updateMetric(ctx, "gauge", "Lookups", gauge(mem.Lookups))
		updateMetric(ctx, "gauge", "MCacheInuse", gauge(mem.MCacheInuse))
		updateMetric(ctx, "gauge", "MCacheSys", gauge(mem.MCacheSys))
		updateMetric(ctx, "gauge", "MSpanInuse", gauge(mem.MSpanInuse))
		updateMetric(ctx, "gauge", "MSpanSys", gauge(mem.MSpanSys))
		updateMetric(ctx, "gauge", "Mallocs", gauge(mem.Mallocs))
		updateMetric(ctx, "gauge", "NextGC", gauge(mem.NextGC))
		updateMetric(ctx, "gauge", "NumForcedGC", gauge(mem.NumForcedGC))
		updateMetric(ctx, "gauge", "NumGC", gauge(mem.NumGC))
		updateMetric(ctx, "gauge", "OtherSys", gauge(mem.OtherSys))
		updateMetric(ctx, "gauge", "PauseTotalNs", gauge(mem.PauseTotalNs))
		updateMetric(ctx, "gauge", "StackInuse", gauge(mem.StackInuse))
		updateMetric(ctx, "gauge", "StackSys", gauge(mem.StackSys))
		updateMetric(ctx, "gauge", "Sys", gauge(mem.Sys))
		updateMetric(ctx, "gauge", "TotalAlloc", gauge(mem.TotalAlloc))
		updateMetric(ctx, "gauge", "RandomValue", gauge(mem.RandomValue))
		updateMetric(ctx, "counter", "PollCount", gauge(mem.PollCount))

		log.Printf("metrics sent %v", mem.PollCount)
	}
}

func readMetrics(ctx context.Context, pollInterval time.Duration, mem *met) {
	var rtm runtime.MemStats
	count := 0
	for {

		<-time.After(pollInterval)

		runtime.ReadMemStats(&rtm)

		mem.Alloc = gauge(rtm.Alloc)
		mem.BuckHashSys = gauge(rtm.BuckHashSys)
		mem.Frees = gauge(rtm.Frees)
		mem.GCCPUFraction = gauge(rtm.GCCPUFraction)
		mem.GCSys = gauge(rtm.GCSys)
		mem.HeapAlloc = gauge(rtm.HeapAlloc)
		mem.HeapIdle = gauge(rtm.HeapInuse)
		mem.HeapInuse = gauge(rtm.HeapInuse)
		mem.HeapObjects = gauge(rtm.HeapObjects)
		mem.HeapReleased = gauge(rtm.HeapReleased)
		mem.HeapSys = gauge(rtm.HeapSys)
		mem.LastGC = gauge(rtm.LastGC)
		mem.Lookups = gauge(rtm.Lookups)
		mem.MCacheInuse = gauge(rtm.MCacheInuse)
		mem.MCacheSys = gauge(rtm.MCacheSys)
		mem.MSpanInuse = gauge(rtm.MSpanInuse)
		mem.MSpanSys = gauge(rtm.MSpanSys)
		mem.Mallocs = gauge(rtm.Mallocs)
		mem.NextGC = gauge(rtm.NextGC)
		mem.NumForcedGC = gauge(rtm.NumForcedGC)
		mem.NumGC = gauge(rtm.NumGC)
		mem.OtherSys = gauge(rtm.OtherSys)
		mem.PauseTotalNs = gauge(rtm.PauseTotalNs)
		mem.StackInuse = gauge(rtm.StackInuse)
		mem.StackSys = gauge(rtm.StackSys)
		mem.Sys = gauge(rtm.Sys)
		mem.TotalAlloc = gauge(rtm.TotalAlloc)
		mem.RandomValue = gauge(rand.Float64())

		count++
		mem.PollCount = counter(count)

		log.Printf("metrics updated %v", count)
	}
}
