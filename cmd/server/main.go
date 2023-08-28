package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	compress "github.com/tritonol/metrics-collecting.git/internal/middleware/compressor"
	middleware "github.com/tritonol/metrics-collecting.git/internal/middleware/logger/zap"
	"github.com/tritonol/metrics-collecting.git/internal/server/config"
	"github.com/tritonol/metrics-collecting.git/internal/server/handlers/metrics/get"
	"github.com/tritonol/metrics-collecting.git/internal/server/handlers/metrics/save"
	"github.com/tritonol/metrics-collecting.git/internal/storage/memstorage"
	"go.uber.org/zap"
)

func main() {
	cfg := config.MustLoad()
	err := http.ListenAndServe(cfg.Server.Address, MetricRouter())
	if err != nil {
		panic(err)
	}
}

func MetricRouter() chi.Router {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	r := chi.NewRouter()
	r.Use(middleware.RequestLogger(logger))
	r.Use(compress.GzipMiddleware)

	storage := memstorage.NewMemStorage()

	r.Post("/update/{type}/{name}/{value}", save.New(storage))
	r.Post("/update/", save.NewJSON(storage))

	r.Get("/value/{type}/{name}", get.Get(storage))
	r.Post("/value/", get.GetJSON(storage))

	r.Get("/", get.MainPage(storage))

	return r
}
