package storage

import (
	"context"

	m "github.com/tritonol/metrics-collecting.git/internal/structs/JSON"
)

type Storage interface {
	StoreMetric(ctx context.Context, name string, mType string, value float64, delta int64) error
	GetAllDataStructed(ctx context.Context) (map[string]m.Metrics, error)
	SaveAllDataStructured(ctx context.Context, metrics map[string]m.Metrics) error
	GetMetrics(ctx context.Context) (map[string]m.Metric, error)
	GetMetric(ctx context.Context ,name string, mType string) (m.Metric, error)
	BatchUpdate(ctx context.Context, metrics []m.Metrics) error
	Ping(ctx context.Context) error
}