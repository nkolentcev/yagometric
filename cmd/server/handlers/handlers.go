package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/nkolentcev/yagometric/cmd/server/storage"
)

type MyMetricHandler struct {
	storage *storage.MemStorage
}

func NewMetricHandler(storage *storage.MemStorage) *MyMetricHandler {
	return &MyMetricHandler{storage: storage}
}

func (mh MyMetricHandler) UpdateMetric(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)

	name := chi.URLParam(r, "name")

	metricType := chi.URLParam(r, "type")

	if !(metricType == "gauge") && !(metricType == "counter") {
		w.WriteHeader(http.StatusNotFound) //404
		return
	}

	value, err := strconv.ParseFloat(chi.URLParam(r, "value"), 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) //400
		log.Println("unable convert string metric")
		return
	}
	mh.storage.AddMetric(name, value, metricType)

}

func (mh MyMetricHandler) GetMetricValue(w http.ResponseWriter, r *http.Request) {
	///w.WriteHeader(http.StatusOK)
	name := chi.URLParam(r, "name")

	metricType := chi.URLParam(r, "type")

	if !(metricType == "gauge") && !(metricType == "counter") {
		w.WriteHeader(http.StatusNotFound) //404
		return
	}

	resp := mh.storage.GetMetricValue(name)
	if resp == nil {
		w.WriteHeader(http.StatusNotFound) //501 метрика не найдена по имени
		log.Printf("wrong metric name %s\n", name)
		return
	}
	switch metricType {
	case "gauge":
		_, err := w.Write([]byte(fmt.Sprintf("%d\n", resp)))
		if err != nil {
			log.Printf("cant write response on body")
		}
	case "counter":
		_, err := w.Write([]byte(fmt.Sprintf("%f\n", resp)))
		if err != nil {
			log.Printf("cant write response on body")
		}
	default:
		_, err := w.Write([]byte(fmt.Sprintf("%v\n", resp)))
		if err != nil {
			log.Printf("cant write response on body")
		}
	}
}

func (mh MyMetricHandler) GetMetricsValuesList(w http.ResponseWriter, r *http.Request) {
	for n, v := range mh.storage.Metrics {
		samp := fmt.Sprintf("-> %s : %v ;\n", n, v)
		w.Write([]byte(samp))
	}
}
