package logy

import (
	"io"
	"os"
	"sync"
	"sync/atomic"
)

var (
	handlers  = map[string]Handler{}
	handlerMu sync.RWMutex
)

type Handler interface {
	Handle(record Record) error
	SetLevel(level Level)
	Level() Level
	SetEnabled(enabled bool)
	IsEnabled() bool
	IsLoggable(record Record) bool
}

type ConfigurableHandler interface {
	OnConfigure(properties ConfigProperties)
}

func RegisterHandler(name string, handler Handler) {
	defer handlerMu.Unlock()
	handlerMu.Lock()

	handlers[name] = handler
}

type ConsoleHandler struct {
	writer  io.Writer
	enabled atomic.Value
	level   atomic.Value
	format  string
	color   bool
}

func NewConsoleHandler() *ConsoleHandler {
	handler := &ConsoleHandler{
		writer: os.Stderr,
	}

	handler.enabled.Store(true)
	handler.level.Store(LevelDebug)
	return handler
}

func (h *ConsoleHandler) Handle(record Record) error {
	return nil
}

func (h *ConsoleHandler) SetLevel(level Level) {
	h.level.Store(level)
}

func (h *ConsoleHandler) Level() Level {
	return h.level.Load().(Level)
}

func (h *ConsoleHandler) SetEnabled(enabled bool) {
	h.enabled.Store(enabled)
}

func (h *ConsoleHandler) IsEnabled() bool {
	return h.enabled.Load().(bool)
}

func (h *ConsoleHandler) IsLoggable(record Record) bool {
	if !h.IsEnabled() {
		return false
	}

	return record.Level >= h.Level()
}

func (h *ConsoleHandler) onConfigure(config *ConsoleConfig) {
	h.enabled.Store(config.Enabled)
	h.level.Store(config.Level)

	switch config.Target {
	case TargetStdout:
		h.writer = os.Stdout
	case TargetStderr:
		h.writer = os.Stdin
	case TargetDiscard:
		h.writer = io.Discard
	}

	h.format = config.Format
	h.color = config.Color
}
