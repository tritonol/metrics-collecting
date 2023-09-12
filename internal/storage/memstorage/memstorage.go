package memstorage

import (
	"fmt"
	"sync"

	jsonstructs "github.com/tritonol/metrics-collecting.git/internal/structs/JSON"
)

type MemStorage struct {
	gaugeMetrics   map[string]float64
	counterMetrics map[string]int64
	mu             sync.Mutex
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gaugeMetrics:   make(map[string]float64),
		counterMetrics: make(map[string]int64),
	}
}

func (ms *MemStorage) StoreGauge(name string, value float64) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.gaugeMetrics[name] = value
}

func (ms *MemStorage) IncrCounter(name string, value int64) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.counterMetrics[name] += value
}

func (ms *MemStorage) StoreCounter(name string, value int64) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.counterMetrics[name] = value
}

func (ms *MemStorage) GetCounter(name string) (int64, bool) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	resp, ok := ms.counterMetrics[name]
	return resp, ok
}

func (ms *MemStorage) GetGauge(name string) (float64, bool) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	resp, ok := ms.gaugeMetrics[name]
	return resp, ok
}

func (ms *MemStorage) GetAllGauge() map[string]float64 {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	return ms.gaugeMetrics
}

func (ms *MemStorage) GetAllCounter() map[string]int64 {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	return ms.counterMetrics
}

func (ms *MemStorage) GetAllDataStructed() map[string]jsonstructs.Metrics {
	data := make(map[string]jsonstructs.Metrics, 30)
	for k, v := range ms.GetAllGauge() {
		data[k] = jsonstructs.Metrics{
			ID:    k,
			MType: "gauge",
			Value: &v,
		}
	}
	for k, v := range ms.GetAllCounter() {
		data[k] = jsonstructs.Metrics{
			ID:    k,
			MType: "counter",
			Delta: &v,
		}
	}

	return data
}

func (ms *MemStorage) SaveAllDataStructured(metrics map[string]jsonstructs.Metrics) error {
	for k, v := range metrics {
		switch v.MType {
		case "gauge":
			ms.StoreGauge(k, *v.Value)
		case "counter":
			ms.StoreCounter(k, *v.Delta)
		default:
			return fmt.Errorf("invalid type: %s", v.MType)
		}
	}

	return nil
}
