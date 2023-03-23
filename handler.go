package logy

import (
	"io"
	"sync"
	"sync/atomic"
)

const (
	ConsoleHandlerName = "console"
	FileHandlerName    = "file"
	SyslogHandlerName  = "syslog"
)

var (
	handlers = map[string]Handler{
		ConsoleHandlerName: newConsoleHandler(),
		FileHandlerName:    newFileHandler(false),
		SyslogHandlerName:  newSysLogHandler(false),
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
	Writer() io.Writer
}

type ConfigurableHandler interface {
	OnConfigure(config Config) error
}

func RegisterHandler(name string, handler Handler) {
	defer handlerMu.Unlock()
	handlerMu.Lock()

	if name == ConsoleHandlerName || name == FileHandlerName || name == SyslogHandlerName {
		panic("logy: 'console', 'file' and 'syslog' handlers are reserved")
	}

	if handler == nil {
		panic("logy: handler cannot be nil")
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
	keyOverrides         atomic.Value
	additionalFields     atomic.Value
	additionalFieldsJson atomic.Value
	color                atomic.Value

	excludedFields atomic.Value

	timestampKey     atomic.Value
	mappedContextKey atomic.Value
	levelKey         atomic.Value
	loggerKey        atomic.Value
	messageKey       atomic.Value
	errorKey         atomic.Value
	stackTraceKey    atomic.Value

	enabledKeys atomic.Value
}

func (h *commonHandler) initializeHandler() {
	h.SetAdditionalFields(AdditionalFields{})
	h.SetJsonEnabled(false)
	h.SetFormat(DefaultTextFormat)
	h.additionalFieldsJson.Store("")
	h.isConsole.Store(true)
	h.color.Store(false)
	h.excludedFields.Store(ExcludedKeys{})
	h.enabledKeys.Store(allKeysEnabled)
	h.keyOverrides.Store(KeyOverrides{})

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

func (h *commonHandler) applyJsonConfig(jsonConfig *JsonConfig) {
	if jsonConfig != nil {
		h.SetJsonEnabled(jsonConfig.Enabled)
		h.SetKeyOverrides(jsonConfig.KeyOverrides)
		h.SetExcludedKeys(jsonConfig.ExcludedKeys)
		h.SetAdditionalFields(jsonConfig.AdditionalFields)
	} else {
		h.SetJsonEnabled(false)
		h.SetKeyOverrides(KeyOverrides{})
		h.SetExcludedKeys(ExcludedKeys{})
		h.SetAdditionalFields(AdditionalFields{})
	}
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
		formatText(encoder, format, record, console && color, false)

		putTextEncoder(encoder)
	}

	_, err := h.writer.Write(*buf)

	return err
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

func (h *commonHandler) setWriter(writer io.Writer) {
	h.writer = writer
}

func (h *commonHandler) Writer() io.Writer {
	return h.writer
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

func (h *commonHandler) SetKeyOverrides(overrides KeyOverrides) {
	h.resetKeys()

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

	if len(overrides) == 0 {
		h.keyOverrides.Store(KeyOverrides{})
	} else {
		h.keyOverrides.Store(overrides)
	}
}

func (h *commonHandler) KeyOverrides() KeyOverrides {
	copyOfOverrides := make(KeyOverrides)
	overrides := h.keyOverrides.Load().(KeyOverrides)

	for key, value := range overrides {
		copyOfOverrides[key] = value
	}

	return copyOfOverrides
}

func (h *commonHandler) SetExcludedKeys(excludedKeys ExcludedKeys) {
	filteredKeys := make(ExcludedKeys, 0)
	enabledKeysFlag := allKeysEnabled

	timestampKey := h.timestampKey.Load().(string)
	mappedContextKey := h.mappedContextKey.Load().(string)
	levelKey := h.levelKey.Load().(string)
	loggerKey := h.loggerKey.Load().(string)
	messageKey := h.messageKey.Load().(string)
	errorKey := h.errorKey.Load().(string)
	stackTraceKey := h.stackTraceKey.Load().(string)

	for _, excludedKey := range excludedKeys {
		switch excludedKey {
		case timestampKey:
			enabledKeysFlag ^= timestampKeyEnabled
		case mappedContextKey:
			enabledKeysFlag ^= mappedContextKeyEnabled
		case levelKey:
			enabledKeysFlag ^= levelKeyEnabled
		case loggerKey:
			enabledKeysFlag ^= loggerKeyEnabled
		case messageKey:
			enabledKeysFlag ^= messageKeyEnabled
		case errorKey:
			enabledKeysFlag ^= errorKeyEnabled
		case stackTraceKey:
			enabledKeysFlag ^= stackTraceKeyEnabled
		default:
			filteredKeys = append(filteredKeys, excludedKey)
		}
	}

	h.enabledKeys.Store(enabledKeysFlag)
	h.excludedFields.Store(filteredKeys)
}

func (h *commonHandler) ExcludedKeys() ExcludedKeys {
	excludedKeys := make(ExcludedKeys, 0)
	enabledKeys := h.enabledKeys.Load().(int)

	if enabledKeys&timestampKeyEnabled == 0 {
		excludedKeys = append(excludedKeys, TimestampKey)
	}

	if enabledKeys&mappedContextKeyEnabled == 0 {
		excludedKeys = append(excludedKeys, MappedContextKey)
	}

	if enabledKeys&levelKeyEnabled == 0 {
		excludedKeys = append(excludedKeys, LevelKey)
	}

	if enabledKeys&loggerKeyEnabled == 0 {
		excludedKeys = append(excludedKeys, LoggerKey)
	}

	if enabledKeys&messageKeyEnabled == 0 {
		excludedKeys = append(excludedKeys, MessageKey)
	}

	if enabledKeys&errorKeyEnabled == 0 {
		excludedKeys = append(excludedKeys, ErrorKey)
	}

	if enabledKeys&stackTraceKeyEnabled == 0 {
		excludedKeys = append(excludedKeys, StackTraceKey)
	}

	excludedKeys = append(excludedKeys, h.excludedFields.Load().(ExcludedKeys)...)
	return excludedKeys
}

func (h *commonHandler) SetAdditionalFields(additionalFields AdditionalFields) {
	if len(additionalFields) == 0 {
		additionalFields = AdditionalFields{}
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

func (h *commonHandler) AdditionalFields() AdditionalFields {
	additionalFields := h.additionalFields.Load().(AdditionalFields)
	copyOfFields := make(AdditionalFields, 0)

	for key, value := range additionalFields {
		copyOfFields[key] = value
	}

	return copyOfFields
}
