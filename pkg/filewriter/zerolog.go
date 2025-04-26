package filewriter

import (
	"github.com/rs/zerolog"
)

type ZerologdWriter struct {
	fw  *FileWriter
	lvl zerolog.Level
}

func NewZerologdWriter(fw *FileWriter, lvl zerolog.Level) *ZerologdWriter {
	return &ZerologdWriter{
		fw:  fw,
		lvl: lvl,
	}
}

func (zw *ZerologdWriter) Write(p []byte) (n int, err error) {
	return zw.fw.Write(p)
}

func (zw *ZerologdWriter) WriteLevel(lvl zerolog.Level, p []byte) (int, error) {
	if lvl >= zw.lvl {
		return zw.fw.Write(p)
	}

	return len(p), nil
}

func (zw *ZerologdWriter) Close() error {
	return zw.fw.Close()
}
