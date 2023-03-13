package logy

import (
	"io"
	"sync"
	"sync/atomic"
)

var (
	handlers = map[string]Handler{
		"console": newConsoleHandler(),
		"file":    newFileHandler(),
		"syslog":  newSysLogHandler(),
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
	OnConfigure(config Config) error
}

func RegisterHandler(name string, handler Handler) {
	defer handlerMu.Unlock()
	handlerMu.Lock()

	if name == "console" || name == "file" || name == "syslog" {
		panic("logy: 'console', 'file' and 'syslog' handlers are reserved")
	}

	handlers[name] = handler
}

func GetHandler(name string) (Handler, bool) {
	defer handlerMu.Unlock()
	handlerMu.Lock()

	if handler, ok := handlers[name]; ok {
		return handler, true
	}

	return nil, false
}

type commonHandler struct {
	target               atomic.Value
	writer               io.Writer
	enabled              atomic.Value
	level                atomic.Value
	json                 atomic.Value
	format               atomic.Value
	isConsole            atomic.Value
	additionalFields     atomic.Value
	additionalFieldsJson atomic.Value
	color                atomic.Value

	timestampKey     atomic.Value
	mappedContextKey atomic.Value
	levelKey         atomic.Value
	loggerKey        atomic.Value
	messageKey       atomic.Value
	errorKey         atomic.Value
	stackTraceKey    atomic.Value
}

func (h *commonHandler) initializeHandler() {
	h.SetAdditionalFields(JsonAdditionalFields{})
	h.SetJsonEnabled(false)
	h.additionalFieldsJson.Store("")
	h.isConsole.Store(true)
	h.color.Store(false)

	h.resetKeys()
}

func (h *commonHandler) resetKeys() {
	h.timestampKey.Store(TimestampKey)
	h.mappedContextKey.Store(MappedContextKey)
	h.levelKey.Store(LevelKey)
	h.loggerKey.Store(LoggerKey)
	h.messageKey.Store(MessageKey)
	h.errorKey.Store(ErrorKey)
	h.stackTraceKey.Store(StackTraceKey)
}

func (h *commonHandler) overrideKeys(overrides KeyOverrides) {
	if len(overrides) == 0 {
		return
	}

	for key, value := range overrides {
		switch key {
		case TimestampKey:
			h.timestampKey.Store(value)
		case MappedContextKey:
			h.mappedContextKey.Store(value)
		case LevelKey:
			h.levelKey.Store(value)
		case LoggerKey:
			h.loggerKey.Store(value)
		case MessageKey:
			h.messageKey.Store(value)
		case ErrorKey:
			h.errorKey.Store(value)
		case StackTraceKey:
			h.stackTraceKey.Store(value)
		}
	}
}

func (h *commonHandler) applyJsonConfig(jsonConfig *JsonConfig) {
	if jsonConfig != nil {
		h.SetJsonEnabled(true)
		h.SetAdditionalFields(jsonConfig.AdditionalFields)
		h.overrideKeys(jsonConfig.KeyOverrides)
	} else {
		h.SetJsonEnabled(false)
		h.SetAdditionalFields(JsonAdditionalFields{})
		h.resetKeys()
	}
}

func (h *commonHandler) Writer() io.Writer {
	return h.writer
}

func (h *commonHandler) Handle(record Record) error {
	buf := newBuffer()
	defer buf.Free()

	json := h.json.Load().(bool)
	console := h.isConsole.Load().(bool)
	color := h.color.Load().(bool)

	if json {
		encoder := getJSONEncoder(buf)

		buf.WriteByte('{')
		h.formatJson(encoder, record)

		buf.WriteByte('}')
		buf.WriteByte('\n')
		putJSONEncoder(encoder)
	} else {
		encoder := getTextEncoder(buf)

		format := h.format.Load().(string)
		h.formatText(encoder, format, record, console && color, false)

		putTextEncoder(encoder)
	}

	_, err := h.writer.Write(*buf)

	return err
}

func (h *commonHandler) setWriter(writer io.Writer) {
	h.writer = writer
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

func (h *commonHandler) IsJsonEnabled() bool {
	return h.json.Load().(bool)
}

func (h *commonHandler) SetAdditionalFields(additionalFields JsonAdditionalFields) {
	if len(additionalFields) == 0 {
		additionalFields = JsonAdditionalFields{}
	}

	h.additionalFields.Store(additionalFields)

	buf := newBuffer()
	encoder := getJSONEncoder(buf)

	for name, value := range additionalFields {
		encoder.AddAny(name, value)
	}

	h.additionalFieldsJson.Store(buf.String())

	buf.Free()
	putJSONEncoder(encoder)
}

func (h *commonHandler) AdditionalFields() JsonAdditionalFields {
	additionalFields := h.additionalFields.Load().(JsonAdditionalFields)
	copyOfFields := make(JsonAdditionalFields, 0)

	for key, value := range additionalFields {
		copyOfFields[key] = value
	}

	return copyOfFields
}
