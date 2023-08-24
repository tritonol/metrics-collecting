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

func (ms *MemStorage) GetCounter(name string) (int64, bool) {
	resp, ok := ms.counterMetrics[name]
	return resp, ok
}

func (ms *MemStorage) GetGauge(name string) (float64, bool) {
	resp, ok := ms.gaugeMetrics[name]
	return resp, ok
}

func (ms *MemStorage) GetAllGauge() map[string]float64 {
	return ms.gaugeMetrics
}

func (ms *MemStorage) GetAllCounter() map[string]int64 {
	return ms.counterMetrics
}