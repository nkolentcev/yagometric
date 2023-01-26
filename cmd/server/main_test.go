package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	compress "github.com/nkolentcev/yagometric/internal/Compress"
	"github.com/nkolentcev/yagometric/internal/handlers"
	"github.com/nkolentcev/yagometric/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type gauge float64
type counter int64

func TestMain(t *testing.T) {
	ms := storage.NewMemStorage()
	zp := compress.NewZipper()
	handler := handlers.NewMetricHandler(ms, zp)
	r := handler.Router()
	st := httptest.NewServer(r)
	defer st.Close()

	ms.AddMetric("sys", 911.911)
	ms.UpdateCounter("PollCount", 50)

	status, body := tsstRequest(t, st, "GET", "/value/gauge/sys")
	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, "911.911\n", body)

	status, body = tsstRequest(t, st, "GET", "/value/counter/PollCount")
	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, "50\n", body)

	status, body = tsstRequest(t, st, "GET", "/value/undercover/PollCount")
	assert.Equal(t, http.StatusNotImplemented, status)
	assert.Equal(t, "", body)

	status, body = tsstRequest(t, st, "POST", "/update/counter/PollCount/3")
	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, "", body)

	status, body = tsstRequest(t, st, "POST", "/update/counter/PollCount/3.141592")
	assert.Equal(t, http.StatusBadRequest, status)
	assert.Equal(t, "", body)
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
