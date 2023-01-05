package storage

import "sync"

type MemStorage struct {
	metrics map[string]float64
	mutex   sync.Mutex
}

// AddMetrics implements Storage
func (*MemStorage) AddMetrics(name string, value float64) {
	panic("unimplemented")
}

type Storage interface {
	AddMetrics(name string, value float64)
}

func NewMemStorage() *MemStorage {
	var ms MemStorage
	ms.metrics = make(map[string]float64)
	ms.mutex = sync.Mutex{}
	return &ms
}

func (ms *MemStorage) AddMetric(name string, value float64) {
	ms.mutex.Lock()
	ms.metrics[name] = value
	ms.mutex.Unlock()
}
