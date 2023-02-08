package keeper

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/nkolentcev/yagometric/internal/config"
	"github.com/nkolentcev/yagometric/internal/storage"
)

type Keeper struct {
	Restore       bool
	StoreInterval time.Duration
	filepath      string
	file          *os.File
	reader        *bufio.Reader
	writer        *bufio.Writer
}

type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

type RestoredCache struct {
	Metrics  map[string]float64
	Counters map[string]int
}

func New(cfg *config.ServerCfg) *Keeper {
	var k Keeper
	k.Restore = cfg.Restore
	k.StoreInterval = cfg.StoreInterval
	k.filepath = cfg.FilePath
	return &k
}

func (k *Keeper) Work(ms *storage.MemStorage) error {
	for {
		<-time.After(k.StoreInterval)
		err := k.WriteCaсhe(ms)
		if err != nil {
			return err
		}

	}
}

func (k *Keeper) WriteCaсhe(ms *storage.MemStorage) (err error) {

	file, _ := os.OpenFile(k.filepath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	k.writer = bufio.NewWriter(file)

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

func (k *Keeper) RestoreCache() (m RestoredCache, err error) {
	var data []byte
	file, _ := os.OpenFile(k.filepath, os.O_RDONLY|os.O_CREATE, 0777)
	k.reader = bufio.NewReader(file)

	met := new(Metrics)
	for {
		if data, err = k.reader.ReadBytes('\n'); err != nil {
			return m, err
		}
		log.Println(data)
		if err = json.Unmarshal(data, met); err != nil {
			return m, err
		}

		if met.MType == "gauge" {
			m.Metrics[met.ID] = *met.Value
		}
		if met.MType == "counter" {
			m.Counters[met.ID] = int(*met.Delta)
		}
	}
}
