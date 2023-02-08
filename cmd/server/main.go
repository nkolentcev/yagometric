package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nkolentcev/yagometric/internal/compress"
	"github.com/nkolentcev/yagometric/internal/config"
	"github.com/nkolentcev/yagometric/internal/handlers"
	"github.com/nkolentcev/yagometric/internal/keeper"
	"github.com/nkolentcev/yagometric/internal/storage"
)

func main() {

	scfg := config.NewServerCfg()
	keeper := keeper.New(scfg)
	storage := storage.NewMemStorage(keeper)

	if scfg.Restore {
		rd, err := keeper.RestoreCache()
		if err != nil {
			log.Printf("restore error: %v", err)
		}
		storage.Counters = rd.Counters
		storage.Metrics = rd.Metrics
	}

	go storage.Keeper.Work(storage)

	zipper := compress.NewZipper()

	handler := handlers.NewMetricHandler(storage, zipper)
	routers := handler.Router()
	r := chi.NewRouter()
	r.Mount("/", routers)

	err := http.ListenAndServe(scfg.Address, r)
	if err != nil {
		log.Println(err)
	}
}
