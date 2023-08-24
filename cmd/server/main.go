package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/tritonol/metrics-collecting.git/internal/server/config"
	"github.com/tritonol/metrics-collecting.git/internal/server/handlers/metrics/get"
	"github.com/tritonol/metrics-collecting.git/internal/server/handlers/metrics/save"
	"github.com/tritonol/metrics-collecting.git/internal/storage/memstorage"
)

func main() {
	cfg := config.MustLoad()

	err := http.ListenAndServe(cfg.Server.Address, MetricRouter())
	if err != nil {
		panic(err)
	}
}

func MetricRouter() chi.Router {
	r := chi.NewRouter()
	storage := memstorage.NewMemStorage()

	r.Post("/update/{type}/{name}/{value}", save.New(storage))
	r.Get("/value/{type}/{name}", get.Get(storage))
	r.Get("/", get.MainPage(storage))
	return r
}
