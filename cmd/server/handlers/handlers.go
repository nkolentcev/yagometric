package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/nkolentcev/yagometric/cmd/server/storage"
)

type MyMetricHandler struct {
	storage storage.Storage
}

func NewMetricHandler(storage storage.Storage) MyMetricHandler {
	return MyMetricHandler{storage: storage}
}

func (mh MyMetricHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(rw, "its not a post method on request", http.StatusNotFound) //404
		return
	}
	url := r.URL.Path
	url = strings.TrimPrefix(url, "/")
	url = strings.TrimSuffix(url, "/")
	sl := strings.Split(url, "/")

	if len(sl) != 4 {
		http.Error(rw, "wrong endpoint schema", http.StatusBadGateway) //400
		return
	}

	metricName := sl[2]
	metricValue, err := strconv.ParseFloat(sl[3], 64)
	if err != nil {
		http.Error(rw, "wrong metric value", http.StatusBadRequest) //400
	}
	mh.storage.AddMetric(metricName, metricValue)

	rw.Header().Set("Content-Type", "text/plain")
	rw.WriteHeader(http.StatusOK)
	rw.Write(nil)
}
