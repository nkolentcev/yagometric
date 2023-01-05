package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/nkolentcev/yagometric/cmd/server/handlers"
	"github.com/nkolentcev/yagometric/cmd/server/storage"
)

const endpoint = "127.0.0.1:8080"

func main() {

	memStorage := storage.NewMemStorage()
	handler := handlers.NewMetricHandler(memStorage)
	mux := http.NewServeMux()
	mux.Handle("/update/", handler)

	fmt.Println("Start server")
	server := &http.Server{
		Addr:    endpoint,
		Handler: mux,
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
