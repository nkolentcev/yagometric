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
	name := chi.URLParam(r, "name")
	value, err := strconv.ParseFloat(chi.URLParam(r, "value"), 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Panicln("unable convert string metric")
	}
	mh.storage.AddMetric(name, value)
}

func (mh MyMetricHandler) GetMetricValue(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	resp := mh.storage.GetMetricValue(name)
	if resp == nil {
		w.WriteHeader(http.StatusNotFound)
		log.Println("wrong metric name %s", name)
	}
	_, err := w.Write([]byte(fmt.Sprintf("%v\n", resp)))
	if err != nil {
		log.Printf("cant write response on body")
	}
	w.WriteHeader(http.StatusOK)
}

func (mh MyMetricHandler) GetMetricsValuesList(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Actual metric value\n <ul>"))
	fmt.Printf("%v", mh.storage.Metrics)
	for n, v := range mh.storage.Metrics {
		samp := fmt.Sprintf("<li> %s : %v </li>", n, v)
		w.Write([]byte(samp))
	}
	w.Write([]byte("</ul>"))
}
