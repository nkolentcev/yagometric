package main

import (
	"log"
	"net/http"

	"github.com/nkolentcev/yagometric/internal/handlers"
)

const endpoint = ":8080"

func main() {

	r := handlers.Router()
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Println(err)
	}
}
