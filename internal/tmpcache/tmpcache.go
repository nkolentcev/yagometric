package tmpcache

import (
	"bufio"
	"encoding/json"
	"os"
	"time"

	"github.com/nkolentcev/yagometric/internal/config"
	"github.com/nkolentcev/yagometric/internal/storage"
)

type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

type Cache struct {
	cfg     *config.ServerCfg
	storage *storage.MemStorage
	file    *os.File
	reader  *bufio.Reader
	writer  *bufio.Writer
}

func NewSaveCache(cfg *config.ServerCfg, storage *storage.MemStorage) *Cache {
	file, _ := os.OpenFile(cfg.FilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	return &Cache{
		cfg:     cfg,
		storage: storage,
		file:    file,
		reader:  bufio.NewReader(file),
		writer:  bufio.NewWriter(file),
	}
}

func NewReaderCache(cfg *config.ServerCfg, storage *storage.MemStorage) *Cache {
	file, _ := os.OpenFile(cfg.FilePath, os.O_RDONLY|os.O_CREATE, 0777)
	return &Cache{
		cfg:     cfg,
		storage: storage,
		file:    file,
		reader:  bufio.NewReader(file),
		writer:  bufio.NewWriter(file),
	}
}

func (c *Cache) Work() (err error) {
	for {
		<-time.After(c.cfg.StoreInterval)
		err := c.WriteCash()
		if err != nil {
			return err
		}
	}
}

func (c *Cache) Stop() (err error) {
	return c.file.Close()
}

func (c *Cache) WriteCash() (err error) {
	var data []byte
	met := new(Metrics)
	met.MType = "gauge"
	for m, v := range c.storage.Metrics {
		met.ID = m
		met.Value = &v
		met.Delta = nil
		data, err = json.Marshal(met)
		if err != nil {
			return
		}
		_, err = c.writer.Write(data)

		if err != nil {
			return
		}
		if err = c.writer.WriteByte('\n'); err != nil {
			return
		}
	}

	met.MType = "counter"
	for m, v := range c.storage.Metrics {
		met.ID = m
		met.Value = nil
		tmp := int64(v)
		met.Delta = &tmp
		data, err = json.Marshal(met)
		if err != nil {
			return
		}
		_, err = c.writer.Write(data)

		if err != nil {
			return
		}
		if err = c.writer.WriteByte('\n'); err != nil {
			return
		}
	}

	return c.writer.Flush()
}

func (c *Cache) ReadeCache() (data []byte, err error) {
	met := new(Metrics)
	for {
		if data, err = c.reader.ReadBytes('\n'); err != nil {
			return nil, err
		}

		if err = json.Unmarshal(data, met); err != nil {
			return nil, err
		}

		if met.MType == "gauge" {
			c.storage.Metrics[met.ID] = *met.Value
			*met.Value = 0
			*met.Delta = 0
		}
		if met.MType == "counter" {
			c.storage.Counters[met.ID] = int(*met.Delta)
			*met.Value = 0
			*met.Delta = 0
		}
	}
}
