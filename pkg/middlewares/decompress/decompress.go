package decompress

import (
	"io"
	"net/http"

	"github.com/klauspost/compress/gzip"
	"github.com/klauspost/compress/zlib"
)

func deleteContentEncoding(req *http.Request) {
	req.Header.Del("Content-Encoding")
	req.ContentLength = -1
}

const errMsgFailedToDecompressBody = "failed to decompress body"

// Decompress is an HTTP middleware that decompresses request body.
// It supports gzip and deflate encoded requests only as most
// popular one. Before calling this middleware, ensure that a
// request encoding format are verified for support. If not, the
// request body, encoded with, for example, Brotli, will be passed
// directly to the handler, which will not be able to properly read
// it.
func Decompress(next http.Handler) http.Handler {
	hfn := func(rw http.ResponseWriter, req *http.Request) {
		encoding := req.Header.Get("Content-Encoding")
		switch encoding {
		case "deflate":
			zr := zlibReaderPool.Get().(io.ReadCloser)
			err := zr.(zlib.Resetter).Reset(req.Body, nil)
			if err != nil {
				http.Error(rw, errMsgFailedToDecompressBody, http.StatusBadRequest)
				return
			}
			defer zr.Close()

			deleteContentEncoding(req)
			req.Body = &pooledZlibReader{
				zr:   zr,
				pool: &zlibReaderPool,
			}

		case "gzip":
			gr := gzipReaderPool.Get().(*gzip.Reader)
			err := gr.Reset(req.Body)
			if err != nil {
				http.Error(rw, errMsgFailedToDecompressBody, http.StatusBadRequest)
				return
			}
			defer gr.Close()

			deleteContentEncoding(req)
			req.Body = &pooledGzipReader{
				gr:   gr,
				pool: &gzipReaderPool,
			}
		}

		next.ServeHTTP(rw, req)
	}

	return http.HandlerFunc(hfn)
}
