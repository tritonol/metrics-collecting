package routes

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/tritonol/metrics-collecting.git/internal/middleware/compressor"
	middleware "github.com/tritonol/metrics-collecting.git/internal/middleware/logger/zap"
	"github.com/tritonol/metrics-collecting.git/internal/server/handlers/metrics/get"
	"github.com/tritonol/metrics-collecting.git/internal/server/handlers/metrics/save"
	statuscheck "github.com/tritonol/metrics-collecting.git/internal/server/handlers/statusCheck"
	"github.com/tritonol/metrics-collecting.git/internal/storage"
	"go.uber.org/zap"
)

func MetricRouter(ctx context.Context, storage storage.Storage, logger *zap.Logger, key string) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.RequestLogger(logger))
	r.Use(compressor.GzipMiddleware)

	r.Post("/update/{type}/{name}/{value}", save.New(storage))
	r.Post("/update/", save.NewJSON(storage))
	r.Post("/updates/", save.Update(storage, key))

	r.Get("/value/{type}/{name}", get.Get(storage))
	r.Post("/value/", get.GetJSON(storage))

	r.Get("/", get.MainPage(storage))

	r.Get("/ping", statuscheck.Ping(ctx, storage))

	return r
}
