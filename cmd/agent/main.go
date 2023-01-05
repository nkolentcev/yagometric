package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"reflect"
	"runtime"
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

type Metrics struct {
	Alloc,
	BuckHashSys,
	Frees,
	GCCPUFraction,
	GCSys,
	HeapAlloc,
	HeapIdle,
	HeapInuse,
	HeapObjects,
	HeapReleased,
	HeapSys,
	LastGC,
	Lookups,
	MCacheInuse,
	MCacheSys,
	MSpanInuse,
	MSpanSys,
	Mallocs,
	NextGC,
	NumForcedGC,
	NumGC,
	OtherSys,
	PauseTotalNs,
	StackInuse,
	StackSys,
	Sys,
	TotalAlloc,
	RandomValue gauge
	PollCount counter
	//mu        sync.Mutex
}

func main() {
	//ctx, cancel := context.WithCancel(context.Background())
	fmt.Println("start")
	ctx := context.Background()
	var m = new(Metrics)
	go readMetrics(ctx, m)
	go updateMetrics(ctx, m)

	<-ctx.Done()
}

func updateMetrics(ctx context.Context, m *Metrics) {
	var metricType string
	var metricValue string

	client := http.Client{}

	for {
		<-time.After(reportInterval)
		dv := reflect.ValueOf(m).Elem()
		for i := 0; i < dv.NumField(); i++ {

			vf := dv.Field(i)
			data := vf.Interface()
			tf := dv.Type().Field(i)
			switch data.(type) {
			case gauge:
				data = vf.Interface().(gauge)
				metricType = "gauge"
				metricValue = fmt.Sprintf("%f", data)
			case counter:
				data = vf.Interface().(counter)
				metricType = "counter"
				metricValue = fmt.Sprintf("%v", data)
			default:
				log.Panicf("Undefined metric type")
				metricType = "undefined"
				metricValue = fmt.Sprintf("%v", data)
			}

			endpoint := fmt.Sprintf("http://%s:%s/update/%s/%s/%s", host, port, metricType, tf.Name, metricValue) //-> TODO -  сборку адреса в отдельную функцию?
			log.Println(endpoint)
			request, err := http.NewRequest(http.MethodPost, endpoint, nil)
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
	}
}

func readMetrics(ctx context.Context, m *Metrics) {
	var rtm runtime.MemStats
	for {
		<-time.After(pollInterval)

		runtime.ReadMemStats(&rtm)

		//m.mu.Lock()

		m.Alloc = gauge(rtm.Alloc)
		m.BuckHashSys = gauge(rtm.BuckHashSys)
		m.Frees = gauge(rtm.Frees)
		m.GCCPUFraction = gauge(rtm.GCCPUFraction)
		m.GCSys = gauge(rtm.GCSys)
		m.HeapAlloc = gauge(rtm.HeapAlloc)
		m.HeapIdle = gauge(rtm.HeapInuse)
		m.HeapInuse = gauge(rtm.HeapInuse)
		m.HeapObjects = gauge(rtm.HeapObjects)
		m.HeapReleased = gauge(rtm.HeapReleased)
		m.HeapSys = gauge(rtm.HeapSys)
		m.LastGC = gauge(rtm.LastGC)
		m.Lookups = gauge(rtm.Lookups)
		m.MCacheInuse = gauge(rtm.MCacheInuse)
		m.MCacheSys = gauge(rtm.MCacheSys)
		m.MSpanInuse = gauge(rtm.MSpanInuse)
		m.MSpanSys = gauge(rtm.MSpanSys)
		m.Mallocs = gauge(rtm.Mallocs)
		m.NextGC = gauge(rtm.NextGC)
		m.NumForcedGC = gauge(rtm.NumForcedGC)
		m.NumGC = gauge(rtm.NumGC)
		m.OtherSys = gauge(rtm.OtherSys)
		m.PauseTotalNs = gauge(rtm.PauseTotalNs)
		m.StackInuse = gauge(rtm.StackInuse)
		m.StackSys = gauge(rtm.StackSys)
		m.Sys = gauge(rtm.Sys)
		m.TotalAlloc = gauge(rtm.TotalAlloc)
		m.RandomValue = gauge(rand.Float64())

		m.PollCount++

		//m.mu.Unlock()

		log.Printf("metrics updated %v", m.PollCount)
	}
}
