package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nkolentcev/yagometric/cmd/server/handlers"
	"github.com/nkolentcev/yagometric/cmd/server/storage"
)

func TestStorage_Handlers(t *testing.T) {

	ms := storage.NewMemStorage()
	handler := handlers.NewMetricHandler(ms)

	tests := []struct {
		name     string
		url      string
		wantCode int
	}{
		{
			name:     "test code 200 - ok",
			url:      "http://127.0.0.1:8080/update/gauge/Alloc/100",
			wantCode: 200,
		},
		{
			name:     "test code 404 - bad request schema",
			url:      "http://127.0.0.1:8080/update/counter/100",
			wantCode: 400,
		},
		{
			name:     "test code 200 - ok",
			url:      "http://127.0.0.1:8080/update/gauge/Alloc/aaabbb",
			wantCode: 200,
		},
	}

	// test request endpoint avaiable

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, tt.url, nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(handler.ServeHTTP)
			h.ServeHTTP(w, req)
			resp := w.Result()

			if statusCode := resp.StatusCode; statusCode != tt.wantCode {
				t.Errorf("want %d, got %d", tt.wantCode, statusCode)
			}
		})
	}
}
