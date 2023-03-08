package logy

import (
	"io"
	"log"
)

type discarder struct {
}

func (d *discarder) Write(p []byte) (n int, err error) {
	return io.Discard.Write(p)
}

type writer struct {
	logger *Logger
	flags  int
}

func newWriter(logger *Logger) *writer {
	return &writer{
		logger: logger,
		flags:  log.Flags(),
	}
}

func (w *writer) Write(buf []byte) (n int, err error) {
	if !w.logger.IsLoggable(LevelDebug) {
		return 0, nil
	}

	origLen := len(buf)
	if len(buf) > 0 && buf[len(buf)-1] == '\n' {
		buf = buf[:len(buf)-1]
	}

	return origLen, w.logger.logDepth(3, nil, LevelDebug, string(buf))
}
