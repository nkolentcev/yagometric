package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	compress "github.com/nkolentcev/yagometric/internal/Compress"
	"github.com/nkolentcev/yagometric/internal/config"
	"github.com/nkolentcev/yagometric/internal/handlers"
	"github.com/nkolentcev/yagometric/internal/storage"
	"github.com/nkolentcev/yagometric/internal/tmpcache"
)

func main() {

	scfg := config.NewServerCfg()
	storage := storage.NewMemStorage()
	zipper := compress.NewZipper()
	cache := tmpcache.NewReaderCache(scfg, storage)

	if scfg.Restore {
		cache.ReadCache()
	}

	cache = tmpcache.NewSaveCache(scfg, storage)
	go cache.Work()

	handler := handlers.NewMetricHandler(storage, zipper)
	routers := handler.Router()
	r := chi.NewRouter()
	r.Mount("/", routers)

	err := http.ListenAndServe(scfg.Address, r)
	if err != nil {
		log.Println(err)
	}
}
