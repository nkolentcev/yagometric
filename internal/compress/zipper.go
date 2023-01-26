package compress

import (
	"bytes"
	"compress/gzip"
	"fmt"
)

type Zipper struct {
}

func NewZipper() *Zipper {
	return &Zipper{}
}

func (c *Zipper) GZip(data []byte) ([]byte, error) {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)

	if _, err := w.Write(data); err != nil {
		return nil, fmt.Errorf("unable compress data, err: %v", err)
	}
	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("unable write compresed data, err: %v", err)
	}
	return b.Bytes(), nil
}

func (c *Zipper) UnGZip(data []byte) ([]byte, error) {
	var b bytes.Buffer
	r, err := gzip.NewReader(bytes.NewReader(data))

	if err != nil {
		return nil, fmt.Errorf("unable read compresed data, err: %v", err)
	}
	if _, err := b.ReadFrom(r); err != nil {
		return nil, fmt.Errorf("unable decomress data, err: %v", err)
	}
	if err := r.Close(); err != nil {
		return nil, fmt.Errorf("unable write decompres data, err: %v", err)
	}

	return b.Bytes(), nil
}
