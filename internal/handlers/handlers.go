package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/nkolentcev/yagometric/internal/compress"
	"github.com/nkolentcev/yagometric/internal/storage"
)

type MyMetricHandler struct {
	storage *storage.MemStorage
	zipper  *compress.Zipper
}
type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

func (mh MyMetricHandler) Router() *chi.Mux {

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", mh.getMetricsValuesList)
	r.Get("/value/{type}/{name}", mh.getMetricValue)
	r.Post("/update/{type}/{name}/{value}", mh.updateMetric)

	r.Post("/value/", mh.getJSONMetricValue)
	r.Post("/update/", mh.updateJSONMetricValue)

	return r
}

func NewMetricHandler(storage *storage.MemStorage, zipper *compress.Zipper) *MyMetricHandler {
	return &MyMetricHandler{
		storage: storage,
		zipper:  zipper,
	}
}

func (mh MyMetricHandler) updateMetric(w http.ResponseWriter, r *http.Request) {

	name := chi.URLParam(r, "name")

	metricType := chi.URLParam(r, "type")

	switch metricType {
	case "gauge":
		value, err := strconv.ParseFloat(chi.URLParam(r, "value"), 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println("unable convert string metric")
			return
		}
		mh.storage.AddMetric(name, value)
		w.WriteHeader(http.StatusOK)
	case "counter":
		value, err := strconv.Atoi(chi.URLParam(r, "value"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println("unable convert string metric")
			return
		}
		mh.storage.UpdateCounter(name, value)
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusNotImplemented)
		return
	}
}

func (mh MyMetricHandler) getMetricValue(w http.ResponseWriter, r *http.Request) {

	name := chi.URLParam(r, "name")

	metricType := chi.URLParam(r, "type")

	if !(metricType == "gauge") && !(metricType == "counter") {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	switch metricType {
	case "gauge":
		resp := mh.storage.GetMetricValue(name)
		if resp == 0 {
			w.WriteHeader(http.StatusNotFound)
			log.Printf("wrong metric name %s\n", name)
			return
		}
		_, err := w.Write([]byte(fmt.Sprintf("%.3f\n", resp)))
		if err != nil {
			log.Printf("cant write response on body")
		}
	case "counter":
		resp := mh.storage.GetCounter(name)
		if resp == 0 {
			w.WriteHeader(http.StatusNotFound)
			log.Printf("wrong metric name %s\n", name)
			return
		}
		_, err := w.Write([]byte(fmt.Sprintf("%v\n", int(resp))))
		if err != nil {
			log.Printf("cant write response on body")
		}
	default:
		_, err := w.Write([]byte(fmt.Sprintf("%v\n", nil)))
		if err != nil {
			log.Printf("cant write response on body")
		}
	}
}

func (mh MyMetricHandler) getMetricsValuesList(w http.ResponseWriter, r *http.Request) {
	for n, v := range mh.storage.Metrics {
		samp := fmt.Sprintf("-> %s : %v ;\n", n, v)
		w.Write([]byte(samp))
	}
	for n, v := range mh.storage.Counters {
		samp := fmt.Sprintf("-> %s : %v ;\n", n, v)
		w.Write([]byte(samp))
	}
}

func (mh MyMetricHandler) getJSONMetricValue(w http.ResponseWriter, r *http.Request) {

	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("wrong content type")
		return
	}

	if r != nil {
		defer r.Body.Close()
	}

	w.Header().Add("Content-Type", "application/json")

	metric := new(Metrics)
	err := json.NewDecoder(r.Body).Decode(&metric)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Print("unable decode metric data")
		return
	}

	if metric.ID == "" {
		w.WriteHeader(http.StatusNotFound)
		log.Printf("undefined metric name")
		return
	}

	name := string(metric.ID)

	switch metric.MType {
	case "gauge":
		resp := mh.storage.GetMetricValue(name)
		metric.Value = &resp
		metric.Delta = nil

	case "counter":
		resp := mh.storage.GetCounter(name)
		tmp := int64(resp)
		metric.Delta = &tmp
		metric.Value = nil

	default:
		log.Printf("unknown metric type")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	dataJSON, err := json.Marshal(metric)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("unable serialise metric data: %v", metric)
		return
	}

	if r.Header.Get("Accept-Encoding") == "gzip" {
		dataJSON, err = mh.zipper.GZip(dataJSON)
		if err != nil {
			log.Printf("err: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	if r.Header.Get("Accept-Encoding") == "gzip" {
		w.Header().Add("Content-Encoding", "gzip")
	}
	_, _ = w.Write(dataJSON)
}

func (mh MyMetricHandler) updateJSONMetricValue(w http.ResponseWriter, r *http.Request) {

	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("wrong content type")
		return
	}

	b, _ := io.ReadAll(r.Body)
	if r.Header.Get("Content-Encoding") == "gzip" {
		b, _ = mh.zipper.UnGZip(b)
		if b == nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	defer r.Body.Close()

	var metric Metrics
	err := json.Unmarshal(b, &metric)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Print("unable decode metric data")
		return
	}

	if metric.Value == nil && metric.Delta == nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("nil data in receving metric")
		return
	}

	switch metric.MType {
	case "gauge":
		mh.storage.AddMetric(metric.ID, *metric.Value)
		resp := mh.storage.GetMetricValue(metric.ID)
		metric.Value = &resp
		if resp == 0 {
			metric.Value = nil
		}

		metric.Delta = nil
	case "counter":
		mh.storage.UpdateCounter(metric.ID, int(*metric.Delta))
		resp := mh.storage.GetCounter(metric.ID)
		tmp := int64(resp)
		metric.Delta = &tmp
		if resp == 0 {
			metric.Delta = nil
		}
		metric.Value = nil
	default:
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("unknown metric type: %s", metric.ID)
		return
	}

	dataJSON, err := json.Marshal(metric)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("unable serialise metric data: %v", metric)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	if r.Header.Get("Content-Encoding") == "gzip" {
		w.Header().Add("Accept-Encoding", "gzip")
	}
	_, _ = w.Write(dataJSON)
}
