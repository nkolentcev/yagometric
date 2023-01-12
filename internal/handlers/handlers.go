package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/nkolentcev/yagometric/internal/storage"
)

func Router() *chi.Mux {
	memStorage := storage.NewMemStorage()
	handler := newMetricHandler(memStorage)
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Route("/", func(r chi.Router) {
		r.Get("/", handler.getMetricsValuesList)
		r.Get("/value/{type}/{name}", handler.getMetricValue)
		r.Post("/update/{type}/{name}/{value}", handler.updateMetric)
	})
	return r
}

type MyMetricHandler struct {
	storage *storage.MemStorage
}

func newMetricHandler(storage *storage.MemStorage) *MyMetricHandler {
	return &MyMetricHandler{storage: storage}
}

func (mh MyMetricHandler) updateMetric(w http.ResponseWriter, r *http.Request) {

	name := chi.URLParam(r, "name")

	metricType := chi.URLParam(r, "type")

	if !(metricType == "gauge") && !(metricType == "counter") {
		w.WriteHeader(http.StatusNotImplemented) //404
		return
	}

	value, err := strconv.ParseFloat(chi.URLParam(r, "value"), 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) //400
		log.Println("unable convert string metric")
		return
	}

	mh.storage.AddMetric(name, value, metricType)
	w.WriteHeader(http.StatusOK)
}

func (mh MyMetricHandler) getMetricValue(w http.ResponseWriter, r *http.Request) {
	///w.WriteHeader(http.StatusOK)
	name := chi.URLParam(r, "name")

	metricType := chi.URLParam(r, "type")

	if !(metricType == "gauge") && !(metricType == "counter") {
		w.WriteHeader(http.StatusNotImplemented) // 501 в тесте
		return
	}

	resp := mh.storage.GetMetricValue(name)
	if resp == 0 {
		w.WriteHeader(http.StatusNotFound) //501 метрика не найдена по имени
		log.Printf("wrong metric name %s\n", name)
		return
	}
	switch metricType {
	case "gauge":
		_, err := w.Write([]byte(fmt.Sprintf("%.3f\n", resp)))
		if err != nil {
			log.Printf("cant write response on body")
		}
	case "counter":
		_, err := w.Write([]byte(fmt.Sprintf("%v\n", int(resp))))
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

func (mh MyMetricHandler) getMetricsValuesList(w http.ResponseWriter, r *http.Request) {
	for n, v := range mh.storage.Metrics {
		samp := fmt.Sprintf("-> %s : %v ;\n", n, v)
		w.Write([]byte(samp))
	}
}
