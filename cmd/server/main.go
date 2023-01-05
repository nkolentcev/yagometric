package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/nkolentcev/yagometric/cmd/server/handlers"
	"github.com/nkolentcev/yagometric/cmd/server/storage"
)

const endpoint = ":8080"

func main() {

	memStorage := storage.NewMemStorage()
	handler := handlers.NewMetricHandler(memStorage)
	// mux := http.NewServeMux()
	// mux.Handle("/update/", handler)

	fmt.Println("Start server")
	// server := &http.Server{
	// 	Addr:    endpoint,
	// 	Handler: mux,
	// }
	// err := server.ListenAndServe()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Route("/", func(r chi.Router) {
		r.Get("/", handler.GetMetricsValuesList)
		r.Get("/value/gauge/{name}", handler.GetMetricValue)
		r.Get("/value/counter/{name}", handler.GetMetricValue)
		r.Post("/update/gauge/{name}/{value}", handler.UpdateMetric)
		r.Post("/update/counter/{name}/{value}", handler.UpdateMetric)
	})

	log.Println(http.ListenAndServe(":8080", r))
}
