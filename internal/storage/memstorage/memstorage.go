package memstorage

type MemStorage struct{
	gaugeMetrics map[string]float64
	counterMetrics map[string]int64
}

func NewMemStorage() *MemStorage{
	return &MemStorage{
		gaugeMetrics: make(map[string]float64),
		counterMetrics: make(map[string]int64),
	}
}

func (ms *MemStorage) StoreGauge(name string, value float64) {
	ms.gaugeMetrics[name] = value
}

func (ms *MemStorage) IncrCounter(name string, value int64) {
	ms.counterMetrics[name] += value
}