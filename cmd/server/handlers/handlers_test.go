package handlers_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nkolentcev/yagometric/cmd/server/handlers"
	"github.com/nkolentcev/yagometric/cmd/server/storage"
)

func TestMetricHandler(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name    string
		address string
		want    want
	}{
		{
			name:    "Positive test: Gauge metric",
			address: "/update/gauge/Name1/111",
			want: want{
				code:        200,
				contentType: "text/plain",
			},
		},
		{
			name:    "Positive test: Counter metric",
			address: "/update/counter/Name1/111",
			want: want{
				code:        200,
				contentType: "text/plain",
			},
		},
		{
			name:    "Negative test: Unkonwn metric",
			address: "/update/unknown/Name1/111",
			want: want{
				code:        501,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:    "Negative test: Empty value",
			address: "/update/counter/Name1",
			want: want{
				code:        404,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.address, bytes.NewBufferString(""))
			w := httptest.NewRecorder()
			memStorage := storage.NewMemStorage()
			handler := handlers.NewMetricHandler(memStorage)
			h := http.Handler(handler)
			h.ServeHTTP(w, request)
			res := w.Result()
			if res.StatusCode != tt.want.code {
				t.Errorf("Expected status code %d, got %d", tt.want.code, w.Code)
			}
			defer res.Body.Close()
			if res.Header.Get("Content-Type") != tt.want.contentType {
				t.Errorf("Expected Content-Type %s, got %s", tt.want.contentType, res.Header.Get("Content-Type"))
			}
		})
	}
}
