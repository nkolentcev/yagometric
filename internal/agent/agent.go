package agent

import (
	"context"
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

type metrics struct {
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

type Agent struct {
	pollInterval   time.Duration
	reportInterval time.Duration
	reportHost     string
	reportPort     string
}

func NewAgent(cfg *config.AgentCfg) *Agent {
	var agent Agent
	agent.pollInterval = cfg.PollInterval
	agent.reportInterval = cfg.ReportInterval
	agent.reportHost = cfg.Host
	agent.reportPort = cfg.Port
	return &agent
}

func (a *Agent) Start(ctx context.Context) {

	ctx = context.WithValue(ctx, "host", a.reportHost)
	ctx = context.WithValue(ctx, "port", a.reportPort)

	mem := new(metrics)
	go readMetrics(ctx, a.pollInterval, mem)
	go updateMetrics(ctx, a.reportInterval, mem)
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
	host := ctx.Value("host")
	port := ctx.Value("port")
	uri := fmt.Sprintf("http://%s:%s/update/%s/%s/%f", host, port, "gauge", metricName, metricValue)

	fmt.Println(ctx.Value("port"))
	fmt.Println(ctx.Value("host"))
	fmt.Println(uri)

	sendMetrics(ctx, uri)
}

func updateCounterMetric(ctx context.Context, metricName string, metricValue counter) {
	host := ctx.Value("host")
	port := ctx.Value("port")
	uri := fmt.Sprintf("http://%s:%s/update/%s/%s/%v", host, port, "counter", metricName, metricValue)
	sendMetrics(ctx, uri)
}

func updateMetrics(ctx context.Context, reportInterval time.Duration, mem *metrics) {

	for {
		<-time.After(reportInterval)
		updateGaugeMetric(ctx, "Alloc", mem.Alloc)
		updateGaugeMetric(ctx, "BuckHashSys", mem.BuckHashSys)
		updateGaugeMetric(ctx, "Frees", mem.Frees)
		updateGaugeMetric(ctx, "GCCPUFraction", mem.GCCPUFraction)
		updateGaugeMetric(ctx, "GCSys", mem.GCSys)
		updateGaugeMetric(ctx, "HeapAlloc", mem.HeapAlloc)
		updateGaugeMetric(ctx, "HeapIdle", mem.HeapInuse)
		updateGaugeMetric(ctx, "HeapInuse", mem.HeapInuse)
		updateGaugeMetric(ctx, "HeapObjects", mem.HeapObjects)
		updateGaugeMetric(ctx, "HeapReleased", mem.HeapReleased)
		updateGaugeMetric(ctx, "HeapSys", mem.HeapSys)
		updateGaugeMetric(ctx, "LastGC", mem.LastGC)
		updateGaugeMetric(ctx, "Lookups", mem.Lookups)
		updateGaugeMetric(ctx, "MCacheInuse", mem.MCacheInuse)
		updateGaugeMetric(ctx, "MCacheSys", mem.MCacheSys)
		updateGaugeMetric(ctx, "MSpanInuse", mem.MSpanInuse)
		updateGaugeMetric(ctx, "MSpanSys", mem.MSpanSys)
		updateGaugeMetric(ctx, "Mallocs", mem.Mallocs)
		updateGaugeMetric(ctx, "NextGC", mem.NextGC)
		updateGaugeMetric(ctx, "NumForcedGC", mem.NumForcedGC)
		updateGaugeMetric(ctx, "NumGC", mem.NumGC)
		updateGaugeMetric(ctx, "OtherSys", mem.OtherSys)
		updateGaugeMetric(ctx, "PauseTotalNs", mem.PauseTotalNs)
		updateGaugeMetric(ctx, "StackInuse", mem.StackInuse)
		updateGaugeMetric(ctx, "StackSys", mem.StackSys)
		updateGaugeMetric(ctx, "Sys", mem.Sys)
		updateGaugeMetric(ctx, "TotalAlloc", mem.TotalAlloc)
		updateGaugeMetric(ctx, "RandomValue", mem.RandomValue)
		updateCounterMetric(ctx, "PollCount", mem.PollCount)
		log.Printf("metrics sent %v", mem.PollCount)
	}
}

func readMetrics(ctx context.Context, pollInterval time.Duration, mem *metrics) {
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
