package storage

import "sync"

type MemStorage struct {
	Metrics map[string]float64
	mutex   sync.Mutex
}

// // AddMetrics implements Storage
// func (*MemStorage) AddMetrics(name string, value float64) {
// 	panic("unimplemented")
// }

type Storage interface {
	AddMetric(name string, value float64)
	GetMetricValue(name string)
}

func NewMemStorage() *MemStorage {
	var ms MemStorage
	ms.Metrics = make(map[string]float64)
	ms.mutex = sync.Mutex{}
	return &ms
}

func (ms *MemStorage) AddMetric(name string, value float64) {
	ms.mutex.Lock()
	ms.Metrics[name] = value
	ms.mutex.Unlock()
}

func (ms *MemStorage) GetMetricValue(name string) (value any) {
	if _, ok := ms.Metrics[name]; ok {
		value = ms.Metrics[name]
		return value
	}
	return nil
}
