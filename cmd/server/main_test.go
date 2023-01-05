package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/nkolentcev/yagometric/cmd/server/handlers"
	"github.com/nkolentcev/yagometric/cmd/server/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type gauge float64
type counter int64

func TestMain(t *testing.T) {
	ms := storage.NewMemStorage()
	handler := handlers.NewMetricHandler(ms)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Route("/", func(r chi.Router) {
		r.Get("/", handler.GetMetricsValuesList)
		r.Get("/value/{type}/{name}", handler.GetMetricValue)
		r.Post("/update/{type}/{name}/{value}", handler.UpdateMetric)
	})

	st := httptest.NewServer(r)
	defer st.Close()

	ms.AddMetric("Sys", 911.911, "gauge")
	ms.AddMetric("PolCount", 50, "counter")

	status, body := tsstRequest(t, st, "GET", "/")
	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, "-> Sys : 911.911\n", body)
}

func tsstRequest(t *testing.T, st *httptest.Server, method, uri string) (int, string) {

	request, err := http.NewRequest(method, st.URL+uri, nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(request)
	require.NoError(t, err)

	rb, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	resp.Body.Close()
	return resp.StatusCode, string(rb)
}
