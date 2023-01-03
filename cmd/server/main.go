package main

import (
	"fmt"
	"handlers"
	"log"
	"net/http"
	"storage"
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
