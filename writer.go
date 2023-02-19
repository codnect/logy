package logy

import (
	"log"
)

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
	if !w.logger.IsLoggable(LevelInfo) {
		return 0, nil
	}

	var depth int
	if w.flags&(log.Lshortfile|log.Llongfile) != 0 {
		depth = 2
	}

	origLen := len(buf)
	if len(buf) > 0 && buf[len(buf)-1] == '\n' {
		buf = buf[:len(buf)-1]
	}

	return origLen, w.logger.logDepth(depth, nil, LevelDebug, string(buf))
}
