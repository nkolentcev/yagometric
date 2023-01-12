package main

import (
	"log"
	"net/http"

	"github.com/nkolentcev/yagometric/internal/handlers"
)

const endpoint = ":8080"

func main() {

	r := handlers.Router()
	log.Println(http.ListenAndServe(":8080", r))

}
