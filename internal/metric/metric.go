package metric

import (
	"github.com/shirou/gopsutil/v4/process"
	"sync"
)

var instance Metric
var once sync.Once

type metrics struct {
	mu sync.Mutex
	v  map[int]int
}

type Metric interface {
	GetHttpStats() map[int]int
	IncreaseHttpStat(int, *sync.WaitGroup) error
	ResetHttpStat(wg *sync.WaitGroup) error
	GetCpuInfo() []float64
	GetMemoryInfo() *process.MemoryInfoStat
}

func MetricsInstance() Metric {
	once.Do(func() {
		instance = NewMetrics()
	})

	return instance
}

func NewMetrics() Metric {
	return &metrics{v: make(map[int]int)}
}

func (m *metrics) GetHttpStats() map[int]int {
	return m.v
}

func (m *metrics) IncreaseHttpStat(code int, wg *sync.WaitGroup) error {
	m.mu.Lock()
	m.v[code]++

	defer func() {
		wg.Done()
		m.mu.Unlock()
	}()
	return nil
}

func (m *metrics) ResetHttpStat(wg *sync.WaitGroup) error {
	m.mu.Lock()

	for k := range m.v {
		delete(m.v, k)
	}

	defer func() {
		wg.Done()
		m.mu.Unlock()
	}()

	return nil
}
