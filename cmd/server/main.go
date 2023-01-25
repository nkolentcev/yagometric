package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nkolentcev/yagometric/internal/config"
	"github.com/nkolentcev/yagometric/internal/handlers"
	"github.com/nkolentcev/yagometric/internal/storage"
	"github.com/nkolentcev/yagometric/internal/tmpcache"
)

func main() {
	scfg := config.NewServerCfg()
	storage := storage.NewMemStorage()
	// cache := tmpcache.NewReaderCache(scfg, storage)

	// if scfg.Restore {
	// 	cache.ReadeCache()
	// }

	cache := tmpcache.NewSaveCache(scfg, storage)
	go cache.WriteCash()

	handler := handlers.NewMetricHandler(storage)
	routers := handler.Router()
	r := chi.NewRouter()
	r.Mount("/", routers)

	err := http.ListenAndServe(scfg.Address, r)
	if err != nil {
		log.Println(err)
	}
}
