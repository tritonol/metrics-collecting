package memstorage

import (
	"context"
	"fmt"
	"sync"

	m "github.com/tritonol/metrics-collecting.git/internal/structs/JSON"
)

const (
	Gauge   m.MetricType = "gauge"
	Counter m.MetricType = "counter"
)

type MemStorage struct {
	metrics map[string]m.Metric
	mu      sync.RWMutex
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		metrics: make(map[string]m.Metric),
	}
}

func (ms *MemStorage) StoreMetric(ctx context.Context, name string, mType string, value float64, delta int64) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if existingMetric, ok := ms.metrics[name]; ok && existingMetric.Type == Counter {
		existingMetric.Delta += delta
		ms.metrics[name] = existingMetric
	} else {
		ms.metrics[name] = m.Metric{
			Type:  m.MetricType(mType),
			Value: value,
			Delta: delta,
		}
	}

	return nil
}

func (ms *MemStorage) GetMetrics(ctx context.Context) (map[string]m.Metric, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	return ms.metrics, nil
}

func (ms *MemStorage) GetMetric(ctx context.Context, name string, mType string) (m.Metric, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	metric, ok := ms.metrics[name]
	if !ok || metric.Type != m.MetricType(mType) {
		return m.Metric{}, fmt.Errorf("metric not found for name '%s' and type '%s'", name, mType)
	}

	return metric, nil
}

func (ms *MemStorage) GetAllDataStructed(ctx context.Context) (map[string]m.Metrics, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	data := make(map[string]m.Metrics)

	for name, metric := range ms.metrics {
		newMetric := metric
		data[name] = m.Metrics{
			ID:    name,
			MType: string(metric.Type),
			Value: &newMetric.Value,
			Delta: &newMetric.Delta,
		}
	}

	return data, nil
}

func (ms *MemStorage) SaveAllDataStructured(ctx context.Context, metrics map[string]m.Metrics) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	for name, metric := range metrics {
		switch metric.MType {
		case "gauge":
			ms.metrics[name] = m.Metric{
				Type:  m.MetricType(metric.MType),
				Value: *metric.Value,
			}
		case "counter":
			ms.metrics[name] = m.Metric{
				Type:  m.MetricType(metric.MType),
				Delta: *metric.Delta,
			}
		}
	}
	return nil
}

func (ms *MemStorage) Ping(ctx context.Context) error {
	return nil
}

func (pg *MemStorage) BatchUpdate(ctx context.Context, metrics []m.Metrics) error {
	for _, v := range metrics {
		err := pg.StoreMetric(ctx, v.ID, v.MType, *v.Value, *v.Delta)
		if err != nil {
			return err
		}
	}

	return nil
}