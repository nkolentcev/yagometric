package storage

import "sync"

type MemStorage struct {
	Metrics  map[string]float64
	Counters map[string]int
	mutex    sync.Mutex
}

type Storage interface {
	AddMetric(name string, value float64)
	UpdateCounter(name string, value int)
	GetMetricValue(name string)
	GetCounter(name string)
}

func NewMemStorage() *MemStorage {
	var ms MemStorage
	ms.Metrics = make(map[string]float64)
	ms.Counters = make(map[string]int)
	ms.mutex = sync.Mutex{}
	return &ms
}

func (ms *MemStorage) AddMetric(name string, value float64) {
	ms.mutex.Lock()
	ms.Metrics[name] = value
	ms.mutex.Unlock()
}
func (ms *MemStorage) UpdateCounter(name string, value int) {
	ms.mutex.Lock()
	ms.Counters[name] += value
	ms.mutex.Unlock()
}

func (ms *MemStorage) GetMetricValue(name string) (value float64) {
	if _, ok := ms.Metrics[name]; ok {
		value = ms.Metrics[name]
		return value
	}
	return 0
}

func (ms *MemStorage) GetCounter(name string) (value int) {
	if _, ok := ms.Counters[name]; ok {
		value = ms.Counters[name]
		return value
	}
	return 0
}
