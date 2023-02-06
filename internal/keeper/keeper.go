package keeper

import (
	"bufio"
	"encoding/json"
	"os"
	"time"

	"github.com/nkolentcev/yagometric/internal/config"
	"github.com/nkolentcev/yagometric/internal/storage"
)

type Keeper struct {
	Restore       bool
	StoreInterval time.Duration
	file          *os.File
	reader        *bufio.Reader
	writer        *bufio.Writer
}

// type MemStorage struct {
// 	Metrics  map[string]float64
// 	Counters map[string]int
// 	mutex    sync.Mutex
// 	keeper   interface{}
// }

type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

func New(cfg *config.ServerCfg) *Keeper {
	var k Keeper
	k.Restore = cfg.Restore
	k.StoreInterval = cfg.StoreInterval
	return &k
}

func (k *Keeper) Work(ms *storage.MemStorage) error {
	for {
		<-time.After(k.StoreInterval)
		err := k.WriteCaсhe(ms)
		if err != nil {

		}

	}
}

func (k *Keeper) WriteCaсhe(ms *storage.MemStorage) (err error) {
	var data []byte
	met := new(Metrics)
	met.MType = "gauge"
	for m, v := range ms.Metrics {
		met.ID = m
		met.Value = &v
		met.Delta = nil
		data, err = json.Marshal(met)
		if err != nil {
			return
		}
		_, err = k.writer.Write(data)

		if err != nil {
			return
		}
		if err = k.writer.WriteByte('\n'); err != nil {
			return
		}
	}

	met.MType = "counter"
	for m, v := range ms.Counters {
		met.ID = m
		met.Value = nil
		tmp := int64(v)
		met.Delta = &tmp
		data, err = json.Marshal(met)
		if err != nil {
			return
		}
		_, err = k.writer.Write(data)

		if err != nil {
			return
		}
		if err = k.writer.WriteByte('\n'); err != nil {
			return
		}
	}

	return k.writer.Flush()
}
