package slog

import (
	"log"
	"time"
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

	record := w.logger.recordPool.Get().(*Record)
	record.Time = time.Now()
	record.Level = LevelInfo
	record.Message = string(buf)
	record.Context = nil
	record.depth = depth

	return origLen, w.logger.logRecord(record)
}
