package metrics

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

type Metrics struct {
	mu		sync.RWMutex
	gauge   map[string]float64
	counter map[string]int64
}

func NewMetrics() *Metrics {
	return &Metrics{
		gauge:   make(map[string]float64),
		counter: make(map[string]int64),
	}
}

func (m *Metrics) CollectGauge() {
	m.mu.Lock()
	var mem = new(runtime.MemStats)
	runtime.ReadMemStats(mem)

	m.gauge["Alloc"] = float64(mem.Alloc)
	m.gauge["BuckHashSys"] = float64(mem.BuckHashSys)
	m.gauge["Frees"] = float64(mem.Frees)
	m.gauge["GCCPUFraction"] = mem.GCCPUFraction
	m.gauge["GCSys"] = float64(mem.GCSys)
	m.gauge["HeapAlloc"] = float64(mem.HeapAlloc)
	m.gauge["HeapIdle"] = float64(mem.HeapIdle)
	m.gauge["HeapInuse"] = float64(mem.HeapInuse)
	m.gauge["HeapObjects"] = float64(mem.HeapObjects)
	m.gauge["HeapReleased"] = float64(mem.HeapReleased)
	m.gauge["HeapSys"] = float64(mem.HeapSys)
	m.gauge["LastGC"] = float64(mem.LastGC)
	m.gauge["Lookups"] = float64(mem.Lookups)
	m.gauge["MCacheInuse"] = float64(mem.MCacheInuse)
	m.gauge["MCacheSys"] = float64(mem.MCacheSys)
	m.gauge["MSpanInuse"] = float64(mem.MSpanInuse)
	m.gauge["MSpanSys"] = float64(mem.MSpanSys)
	m.gauge["Mallocs"] = float64(mem.Mallocs)
	m.gauge["NextGC"] = float64(mem.NextGC)
	m.gauge["NumForcedGC"] = float64(mem.NumForcedGC)
	m.gauge["NumGC"] = float64(mem.NumGC)
	m.gauge["OtherSys"] = float64(mem.OtherSys)
	m.gauge["PauseTotalNs"] = float64(mem.PauseTotalNs)
	m.gauge["StackInuse"] = float64(mem.StackInuse)
	m.gauge["StackSys"] = float64(mem.StackSys)
	m.gauge["Sys"] = float64(mem.Sys)
	m.gauge["TotalAlloc"] = float64(mem.TotalAlloc)

	m.gauge["RandomValue"] = float64(rand.Intn(100))
	m.mu.Unlock()
}

func (m *Metrics) CollectAdditionalGauge() error {
	m.mu.Lock()
	v, err := mem.VirtualMemory()
	if err != nil {
		return fmt.Errorf("read cpu gopsutil err: %w", err)
	}

	utils, err := cpu.Percent(1*time.Second, true)
	if err != nil {
		return fmt.Errorf("read cpu gopsutil err: %w", err)
	}

	m.gauge["TotalMemory"] = float64(v.Total)
	m.gauge["FreeMemory"] = float64(v.Free)

	for k, value := range utils {
		index := fmt.Sprintf("CPUutilization%d", k)
		m.gauge[index] = value
	}

	m.mu.Unlock()
	return nil
}

func (m *Metrics) GetGauge() map[string]float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.gauge
}

func (m *Metrics) CollectCounter() {
	m.mu.Lock()
	m.counter["PollCount"] += 1
	m.mu.Unlock()
}

func (m *Metrics) GetCounter() map[string]int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.counter
}
