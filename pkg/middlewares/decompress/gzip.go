package decompress

import (
	"bytes"
	"compress/zlib"
	"sync"

	"github.com/klauspost/compress/gzip"
)

var gzipReaderPool = sync.Pool{
	New: func() any {
		return &gzip.Reader{}
	},
}

type pooledGzipReader struct {
	gr   *gzip.Reader
	pool *sync.Pool
}

func (r *pooledGzipReader) Read(b []byte) (int, error) {
	return r.gr.Read(b)
}

func (r *pooledGzipReader) Close() error {
	err := r.gr.Close()
	r.pool.Put(r.gr)
	return err
}

var zlibHeaderReaderBytes = []byte{0x78, 0x9c, 0x01, 0x00, 0x00, 0xff, 0xff, 0x00, 0x00, 0x00, 0x01}

var zlibReaderPool = sync.Pool{
	New: func() any {
		br := bytes.NewReader(zlibHeaderReaderBytes)
		zr, _ := zlib.NewReader(br)
		return zr
	},
}
