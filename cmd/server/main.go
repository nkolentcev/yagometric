package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/nkolentcev/yagometric/internal/handlers"
	"github.com/nkolentcev/yagometric/internal/storage"
)

const endpoint = ":8080"

func main() {

	memStorage := storage.NewMemStorage()
	handler := handlers.NewMetricHandler(memStorage)

	fmt.Println("Start server")

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Route("/", func(r chi.Router) {
		r.Get("/", handler.GetMetricsValuesList)
		r.Get("/value/{type}/{name}", handler.GetMetricValue)
		r.Post("/update/{type}/{name}/{value}", handler.UpdateMetric)
	})

	log.Println(http.ListenAndServe(":8080", r))
}
