package logy

import (
	"io"
	"log"
	"net"
	"sync"
	"sync/atomic"
)

type syncWriter struct {
	discard atomic.Value
	mu      sync.Mutex
	writer  io.Writer
}

func newSyncWriter(writer io.Writer, discard bool) *syncWriter {
	syncW := &syncWriter{
		writer: writer,
	}

	syncW.discard.Store(discard)
	return syncW
}

func (sw *syncWriter) setDiscarded(discarded bool) {
	sw.discard.Store(discarded)
}

func (sw *syncWriter) isDiscarded() bool {
	return sw.discard.Load().(bool)
}

func (sw *syncWriter) Write(p []byte) (n int, err error) {
	if sw.isDiscarded() {
		return 0, nil
	}

	defer sw.mu.Unlock()
	sw.mu.Lock()
	return sw.writer.Write(p)
}

type syslogWriter struct {
	discard atomic.Value
	mu      sync.Mutex
	writer  net.Conn
	network string
	address string
	retry   bool
}

func newSyslogWriter(network, address string, retry bool, discarded bool) *syslogWriter {
	syslogWriter := &syslogWriter{
		network: network,
		address: address,
		retry:   retry,
	}
	syslogWriter.discard.Store(discarded)
	return syslogWriter
}

func (sw *syslogWriter) setDiscarded(discarded bool) {
	sw.discard.Store(discarded)
}

func (sw *syslogWriter) isDiscarded() bool {
	return sw.discard.Load().(bool)
}

func (sw *syslogWriter) connect() error {
	if sw.writer != nil {
		sw.writer.Close()
		sw.writer = nil
	}

	con, err := net.Dial(sw.network, sw.address)
	if err != nil {
		return err
	}

	sw.writer = con
	return nil
}

func (sw *syslogWriter) Write(p []byte) (n int, err error) {
	if sw.isDiscarded() {
		return 0, nil
	}

	defer sw.mu.Unlock()
	sw.mu.Lock()

	if sw.writer != nil {
		return sw.writer.Write(p)
	}

	if !sw.retry {
		return 0, nil
	}

	err = sw.connect()
	if err != nil {
		return 0, err
	}

	return sw.writer.Write(p)
}

type globalWriter struct {
	logger *Logger
	flags  int
}

func newGlobalWriter(logger *Logger) *globalWriter {
	return &globalWriter{
		logger: logger,
		flags:  log.Flags(),
	}
}

func (w *globalWriter) Write(buf []byte) (n int, err error) {
	if !w.logger.IsLoggable(LevelDebug) {
		return 0, nil
	}

	origLen := len(buf)
	if len(buf) > 0 && buf[len(buf)-1] == '\n' {
		buf = buf[:len(buf)-1]
	}

	return origLen, w.logger.logDepth(3, nil, LevelDebug, string(buf))
}
