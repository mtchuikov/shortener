package decompress

import (
	"io"
	"sync"
)

type pooledZlibReader struct {
	zr   io.ReadCloser
	pool *sync.Pool
}

func (r *pooledZlibReader) Read(p []byte) (int, error) {
	return r.zr.Read(p)
}

func (r *pooledZlibReader) Close() error {
	err := r.zr.Close()
	r.pool.Put(r.zr)
	return err
}
