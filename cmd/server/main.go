package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nkolentcev/yagometric/internal/handlers"
	"github.com/nkolentcev/yagometric/internal/storage"
)

const endpoint = ":8080"

func main() {

	storage := storage.NewMemStorage()
	handler := handlers.NewMetricHandler(storage)
	routers := handler.Router()
	r := chi.NewRouter()
	r.Mount("/", routers)

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Println(err)
	}
}
