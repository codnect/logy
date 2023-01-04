package slog

import (
	"io"
	"os"
	"sync"
)

type Handler interface {
	Handle(record *Record) error
	SetFormatter(formatter Formatter)
	Formatter() Formatter
	SetLevel(level Level)
	Level() Level
	IsLoggable(record *Record) bool
}

type ConsoleHandler struct {
	formatter Formatter
	writer    io.Writer
	level     Level

	mu sync.RWMutex
}

func NewConsoleHandler() *ConsoleHandler {
	return &ConsoleHandler{
		formatter: NewSimpleFormatter(),
		writer:    os.Stderr,
		level:     LevelInfo,
		mu:        sync.RWMutex{},
	}
}

func (h *ConsoleHandler) Handle(record *Record) error {
	if h.IsLoggable(record) {
		msg := h.Formatter().Format(record)
		msg = msg + "\n"
		_, err := h.writer.Write([]byte(msg))

		if err != nil {
			return err
		}
	}

	return nil
}

func (h *ConsoleHandler) SetFormatter(formatter Formatter) {
	if formatter == nil {
		return
	}

	defer h.mu.Unlock()
	h.mu.Lock()
	h.formatter = formatter
}

func (h *ConsoleHandler) Formatter() Formatter {
	defer h.mu.Unlock()
	h.mu.Lock()
	return h.formatter
}

func (h *ConsoleHandler) SetLevel(level Level) {
	defer h.mu.Unlock()
	h.mu.Lock()
	h.level = level
}

func (h *ConsoleHandler) Level() Level {
	defer h.mu.Unlock()
	h.mu.Lock()
	return h.level
}

func (h *ConsoleHandler) IsLoggable(record *Record) bool {
	return record.Level >= h.Level()
}
