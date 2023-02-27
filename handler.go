package logy

import (
	"io"
	"sync"
	"sync/atomic"
)

var (
	handlers = map[string]Handler{
		"console": NewConsoleHandler(),
		"file":    NewFileHandler(),
	}
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

	if name == "console" || name == "file" {
		panic("logy: 'console' and 'file' registers are reserved")
	}

	handlers[name] = handler
}

type commonHandler struct {
	target           atomic.Value
	writer           atomic.Value
	enabled          atomic.Value
	level            atomic.Value
	json             atomic.Value
	format           atomic.Value
	excludedKeys     atomic.Value
	additionalFields atomic.Value
	writerMu         sync.RWMutex
	isConsole        atomic.Value
}

func (h *commonHandler) initializeHandler() {
	h.SetExcludedKeys([]string{})
	h.SetAdditionalFields(map[string]JsonAdditionalField{})
	h.SetJsonEnabled(false)
	h.isConsole.Store(true)
}

func (h *commonHandler) applyJsonConfig(jsonConfig *JsonConfig) {
	if jsonConfig != nil {
		h.SetJsonEnabled(true)
		h.SetExcludedKeys(jsonConfig.ExcludeKeys)
		h.SetAdditionalFields(jsonConfig.AdditionalFields)
	} else {
		h.SetJsonEnabled(false)
		h.SetExcludedKeys([]string{})
		h.SetAdditionalFields(map[string]JsonAdditionalField{})
	}
}

func (h *commonHandler) Handle(record Record) error {
	buf := newBuffer()
	defer buf.Free()

	json := h.json.Load().(bool)
	console := h.isConsole.Load().(bool)

	if json {
		encoder := getJSONEncoder()
		encoder.buf = buf

		buf.WriteByte('{')
		excludedKeys := h.excludedKeys.Load().(map[string]struct{})
		additionalFields := h.additionalFields.Load().(map[string]JsonAdditionalField)
		formatJson(encoder, record, excludedKeys, additionalFields)

		buf.WriteByte('}')
		buf.WriteByte('\n')
		putJSONEncoder(encoder)
	} else {
		encoder := getTextEncoder()
		encoder.buf = buf

		format := h.format.Load().(string)
		formatText(encoder, format, record, console)

		putTextEncoder(encoder)
	}

	target := h.target.Load()

	if console {
		targetVal, ok := target.(Target)
		if ok && targetVal == TargetDiscard {
			io.Discard.Write(*buf)
			return nil
		}
	}

	consoleWriter := h.writer.Load().(io.Writer)
	_, err := consoleWriter.Write(*buf)

	return err
}

func (h *commonHandler) setWriter(writer io.Writer) {
	h.writer.Store(writer)
}

func (h *commonHandler) SetLevel(level Level) {
	h.level.Store(level)
}

func (h *commonHandler) Level() Level {
	return h.level.Load().(Level)
}

func (h *commonHandler) SetEnabled(enabled bool) {
	h.enabled.Store(enabled)
}

func (h *commonHandler) IsEnabled() bool {
	return h.enabled.Load().(bool)
}

func (h *commonHandler) IsLoggable(record Record) bool {
	if !h.IsEnabled() {
		return false
	}

	return record.Level <= h.Level()
}

func (h *commonHandler) SetFormat(format string) {
	h.format.Store(format)
}

func (h *commonHandler) Format() string {
	return h.format.Load().(string)
}

func (h *commonHandler) SetJsonEnabled(json bool) {
	h.json.Store(json)
}

func (h *commonHandler) JsonEnabled() bool {
	return h.json.Load().(bool)
}

func (h *commonHandler) SetExcludedKeys(excludedKeys []string) {
	if len(excludedKeys) != 0 {
		excludedKeyMap := map[string]struct{}{}

		for _, key := range excludedKeys {
			excludedKeyMap[key] = struct{}{}
		}

		h.excludedKeys.Store(excludedKeyMap)
	} else {
		h.excludedKeys.Store(map[string]struct{}{})
	}
}

func (h *commonHandler) ExcludedKeys() []string {
	excludedKeys := make([]string, 0)
	excludedKeyMap := h.excludedKeys.Load().(map[string]struct{})

	for key := range excludedKeyMap {
		excludedKeys = append(excludedKeys, key)
	}

	return excludedKeys
}

func (h *commonHandler) SetAdditionalFields(additionalFields map[string]JsonAdditionalField) {
	if len(additionalFields) != 0 {
		h.additionalFields.Store(additionalFields)
	} else {
		h.additionalFields.Store(map[string]JsonAdditionalField{})
	}
}

func (h *commonHandler) AdditionalFields() map[string]JsonAdditionalField {
	additionalFields := h.additionalFields.Load().(map[string]JsonAdditionalField)
	copyOfFields := make(map[string]JsonAdditionalField, 0)

	for key, value := range additionalFields {
		copyOfFields[key] = value
	}

	return copyOfFields
}
