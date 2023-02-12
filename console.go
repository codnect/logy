package logy

import (
	"io"
	"os"
	"sync"
	"sync/atomic"
)

type ConsoleHandler struct {
	writer           atomic.Value
	enabled          atomic.Value
	level            atomic.Value
	json             atomic.Value
	format           atomic.Value
	color            atomic.Value
	excludedKeys     atomic.Value
	additionalFields atomic.Value
	writerMu         sync.RWMutex
}

func NewConsoleHandler() *ConsoleHandler {
	handler := &ConsoleHandler{}

	handler.writer.Store(os.Stderr)
	handler.enabled.Store(true)
	handler.level.Store(LevelDebug)
	handler.json.Store(false)
	handler.excludedKeys.Store(map[string]struct{}{})
	handler.additionalFields.Store(map[string]JsonAdditionalField{})
	return handler
}

func (h *ConsoleHandler) Handle(record Record) error {
	var (
		buf  []byte
		json bool
	)

	json = h.json.Load().(bool)

	if json {
		excludedKeys := h.excludedKeys.Load().(map[string]struct{})
		additionalFields := h.additionalFields.Load().(map[string]JsonAdditionalField)
		formatJson(&buf, record, excludedKeys, additionalFields)
	} else {
		format := h.format.Load().(string)
		formatText(&buf, format, record, true)
	}

	defer h.writerMu.Unlock()
	h.writerMu.Lock()

	consoleWriter := h.writer.Load().(io.Writer)
	_, err := consoleWriter.Write(buf)

	return err
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

	return record.Level <= h.Level()
}

func (h *ConsoleHandler) onConfigure(config *ConsoleConfig) {
	h.enabled.Store(config.Enable)
	h.level.Store(config.Level)

	switch config.Target {
	case TargetStdout:
		//h.writer.Store(os.Stdout)
	case TargetStderr:
		//h.writer.Store(os.Stdin)
	case TargetDiscard:
		h.writer.Store(io.Discard)
	}

	h.format.Store(config.Format)
	h.color.Store(config.Color)

	if config.Json != nil {
		h.json.Store(true)

		if len(config.Json.ExcludeKeys) != 0 {
			excludedKeys := map[string]struct{}{}

			for _, key := range config.Json.ExcludeKeys {
				excludedKeys[key] = struct{}{}
			}

			h.excludedKeys.Store(excludedKeys)
		}

		if len(config.Json.AdditionalFields) != 0 {
			h.additionalFields.Store(config.Json.AdditionalFields)
		}
	}
}
