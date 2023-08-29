package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/tritonol/metrics-collecting.git/internal/middleware/compressor"
	middleware "github.com/tritonol/metrics-collecting.git/internal/middleware/logger/zap"
	"github.com/tritonol/metrics-collecting.git/internal/server/handlers/metrics/get"
	"github.com/tritonol/metrics-collecting.git/internal/server/handlers/metrics/save"
	"github.com/tritonol/metrics-collecting.git/internal/storage/memstorage"
	"go.uber.org/zap"
)

func MetricRouter(storage *memstorage.MemStorage) chi.Router {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	r := chi.NewRouter()
	r.Use(middleware.RequestLogger(logger))
	r.Use(compressor.GzipMiddleware)

	r.Post("/update/{type}/{name}/{value}", save.New(storage))
	r.Post("/update/", save.NewJSON(storage))

	r.Get("/value/{type}/{name}", get.Get(storage))
	r.Post("/value/", get.GetJSON(storage))

	r.Get("/", get.MainPage(storage))

	return r
}
