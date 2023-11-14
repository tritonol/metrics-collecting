package routes

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/tritonol/metrics-collecting.git/internal/middleware/compressor"
	middleware "github.com/tritonol/metrics-collecting.git/internal/middleware/logger/zap"
	"github.com/tritonol/metrics-collecting.git/internal/server/handlers/metrics/get"
	"github.com/tritonol/metrics-collecting.git/internal/server/handlers/metrics/save"
	statuscheck "github.com/tritonol/metrics-collecting.git/internal/server/handlers/statusCheck"
	"github.com/tritonol/metrics-collecting.git/internal/storage/memstorage"
	"github.com/tritonol/metrics-collecting.git/internal/storage/pgstorage"
	"go.uber.org/zap"
)

func MetricRouter(ctx context.Context, db *pgstorage.Postgres, storage *memstorage.MemStorage, logger *zap.Logger) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.RequestLogger(logger))
	r.Use(compressor.GzipMiddleware)

	r.Post("/update/{type}/{name}/{value}", save.New(db))
	r.Post("/update/", save.NewJSON(db))

	r.Get("/value/{type}/{name}", get.Get(db))
	r.Post("/value/", get.GetJSON(db))

	r.Get("/", get.MainPage(db))

	r.Get("/ping", statuscheck.Ping(ctx, db))

	return r
}