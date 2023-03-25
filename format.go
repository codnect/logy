package logy

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

const (
	loggerTargetLen      = 40
	defaultErrorTypeName = "Error"
)

const (
	TimestampKey     = "timestamp"
	MappedContextKey = "mapped_context"
	LevelKey         = "level"
	LoggerKey        = "logger"
	MessageKey       = "message"
	ErrorKey         = "error"
	StackTraceKey    = "stack_trace"
)

const (
	timestampKeyEnabled = 1 << iota
	mappedContextKeyEnabled
	levelKeyEnabled
	loggerKeyEnabled
	messageKeyEnabled
	errorKeyEnabled
	stackTraceKeyEnabled
	allKeysEnabled = timestampKeyEnabled | mappedContextKeyEnabled | levelKeyEnabled | loggerKeyEnabled | messageKeyEnabled | errorKeyEnabled | stackTraceKeyEnabled
)

var (
	processId   = os.Getpid()
	processName = filepath.Base(os.Args[0])
)

func formatContextValues(encoder *textEncoder, ctx context.Context, withKeys bool) {
	contextFields := ContextFieldsFrom(ctx)

	if contextFields != nil {
		fieldLen := len(contextFields.fields)

		for i, field := range contextFields.fields {
			if withKeys {
				encoder.buf.WriteString(field.Key())
				encoder.buf.WriteByte('=')
			}

			encoder.buf.WriteString(field.ValueAsText())

			if i != fieldLen-1 {
				encoder.buf.WriteString(", ")
			}
		}
	}
}

func (h *commonHandler) formatContextValuesAsJson(encoder *jsonEncoder, ctx context.Context) {
	contextFields := ContextFieldsFrom(ctx)
	excludedFields := h.excludedFields.Load().(ExcludedKeys)
	encoder.buf.WriteString("{")

	if contextFields != nil {
		fieldLen := len(contextFields.fields)

		for i, field := range contextFields.fields {

			isExcluded := false
			for _, excluded := range excludedFields {
				if excluded == field.key {
					isExcluded = true
					break
				}
			}

			if isExcluded {
				continue
			}

			jsonVal := field.AsJson()

			if len(jsonVal) != 0 {
				encoder.buf.WriteString(jsonVal[1 : len(jsonVal)-1])
			}

			if i != fieldLen-1 {
				encoder.buf.WriteByte(',')
			}
		}
	}

	encoder.buf.WriteString("}")
}

func formatText(encoder *textEncoder, format string, record Record, color bool, noPadding bool) {
	contextFields := ContextFieldsFrom(record.Context)

	i := 0
	for j := 0; j < len(format); j++ {
		if format[j] == '%' && j+1 < len(format) {
			typ := format[j+1]
			w := 1

			switch typ {
			case 'd': // date
				layout, l := getPlaceholderName(format[j+2:])

				if layout != "" {
					encoder.AppendTimeLayout(record.Time, layout)
				} else {
					encoder.AppendTime(record.Time)
				}

				w = l + 1
			case 'c': // logger
				appendLoggerAsText(encoder.buf, record.LoggerName, color, noPadding)
			case 'p': // level
				appendLevelAsText(encoder.buf, record.Level, color)
			case 'x': // context value without key
				name, l := getPlaceholderName(format[j+2:])

				if contextFields != nil {
					if name != "" {
						field, ok := contextFields.Field(name)
						if ok {
							encoder.buf.WriteString(field.ValueAsText())
						}
					} else {
						formatContextValues(encoder, record.Context, false)
					}
				}

				w = l + 1
			case 'X': // context value with key
				name, l := getPlaceholderName(format[j+2:])

				if contextFields != nil {
					if name != "" {
						field, ok := contextFields.Field(name)
						if ok {
							encoder.buf.WriteString(name)
							encoder.buf.WriteByte('=')
							encoder.buf.WriteString(field.ValueAsText())
						}
					} else {
						formatContextValues(encoder, record.Context, true)
					}
				}

				w = l + 1
			case 'm': // full message
				encoder.AppendString(record.Message)

				if record.Error != nil {
					appendError(encoder.buf, record.Error)
				}

				if record.StackTrace != "" {
					if record.Error != nil {
						encoder.buf.WriteByte('\n')
					}
					encoder.buf.WriteString(strings.ReplaceAll(record.StackTrace, "\\n", "\n"))
				}
			case 's': // simple message
				encoder.AppendString(record.Message)
			case 'M': // method
				encoder.AppendString(record.Caller.Name())
			case 'L': // line
				encoder.AppendInt(record.Caller.Line())
			case 'F': // file
				encoder.AppendString(record.Caller.File())
			case 'C': // package
				encoder.AppendString(record.Caller.Package())
			case 'l': // location
				encoder.AppendString(record.Caller.Path())
			case 'e': // error and stack trace if exist
				if record.Error != nil {
					appendError(encoder.buf, record.Error)
				}

				if record.StackTrace != "" {
					if record.Error != nil {
						encoder.buf.WriteByte('\n')
					}
					encoder.buf.WriteString(strings.ReplaceAll(record.StackTrace, "\\n", "\n"))
				}
			case 'i': // process id
				encoder.buf.WriteIntWidth(processId, 4)
			case 'N': // process name
				encoder.AppendString(processName)
			case 'n': // newline
				encoder.buf.WriteByte('\n')
			case '%':
				encoder.buf.WriteByte('%')
			default:
				encoder.buf.WriteString(format[i:j])
			}

			j += w
			i = j + 1
		} else {
			encoder.buf.WriteByte(format[j])
			i = j + 1
		}
	}
}

func appendError(buf *buffer, err error) {
	var (
		errorTypeName = defaultErrorTypeName
	)

	typ := reflect.TypeOf(err)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	if typ.Name() != "errorString" {
		errorTypeName = typ.String()
	}

	buf.WriteByte('\n')
	buf.WriteString(errorTypeName)
	buf.WriteString(": ")
	buf.WriteString(err.Error())
}

func appendLoggerAsText(buf *buffer, logger string, color bool, noPadding bool) {
	if color {
		colorCyan.start(buf)
		abbreviateLoggerName(buf, logger, loggerTargetLen, noPadding)
		colorCyan.end(buf)
	} else {
		abbreviateLoggerName(buf, logger, loggerTargetLen, noPadding)
	}
}

func appendLevelAsText(buf *buffer, level Level, color bool) {
	str := level.String()
	buf.WritePadding(5 - len(str))

	if color {
		levelColors[level-1].print(buf, str)
	} else {
		buf.WriteString(str)
	}
}

func (h *commonHandler) formatJson(encoder *jsonEncoder, record Record) {
	enabledKeys := h.enabledKeys.Load().(int)

	// timestamp
	if enabledKeys&timestampKeyEnabled == timestampKeyEnabled {
		encoder.AddTime(h.timestampKey.Load().(string), record.Time)
	}

	// level
	if enabledKeys&levelKeyEnabled == levelKeyEnabled {
		encoder.AddString(h.levelKey.Load().(string), record.Level.String())
	}

	// logger name
	if enabledKeys&loggerKeyEnabled == loggerKeyEnabled {
		encoder.AddString(h.loggerKey.Load().(string), record.LoggerName)
	}

	// message
	if enabledKeys&messageKeyEnabled == messageKeyEnabled {
		encoder.AddString(h.messageKey.Load().(string), record.Message)
	}

	if enabledKeys&stackTraceKeyEnabled == stackTraceKeyEnabled && record.StackTrace != "" {
		// stack trace
		encoder.AddString(h.stackTraceKey.Load().(string), record.StackTrace)
	}

	if enabledKeys&errorKeyEnabled == errorKeyEnabled && record.Error != nil {
		// error
		encoder.AddString(h.errorKey.Load().(string), record.Error.Error())
	}

	// mapped context
	if enabledKeys&mappedContextKeyEnabled == mappedContextKeyEnabled && record.Context != nil {
		encoder.addKey(h.mappedContextKey.Load().(string))
		h.formatContextValuesAsJson(encoder, record.Context)
	}

	// additional fields
	additionalFieldsJson := h.additionalFieldsJson.Load().(string)
	if len(additionalFieldsJson) != 0 {
		encoder.buf.WriteByte(',')
		encoder.buf.WriteString(additionalFieldsJson)
	}
}

func getPlaceholderName(s string) (string, int) {
	switch {
	case s[0] == '{':
		if len(s) > 2 && isSpecialVar(s[1]) && s[2] == '}' {
			return s[1:2], 3
		}

		for i := 1; i < len(s); i++ {
			if s[i] == '}' {
				if i == 1 {
					return "", 2
				}
				return s[1:i], i + 1
			}
		}

		return "", 1
	case isSpecialVar(s[0]):
		return s[0:1], 1
	}

	var i int
	for i = 0; i < len(s) && isAlphaNum(s[i]); i++ {
	}

	return s[:i], i
}

func abbreviateLoggerName(buf *buffer, name string, targetLen int, noPadding bool) {
	inLen := len(name)
	if inLen < targetLen {
		buf.WriteString(name)
		if !noPadding {
			buf.WritePadding(loggerTargetLen - inLen)
		}
		return
	}

	trimmed := 0
	inDotState := true
	inSlashState := false
	start := buf.Len()

	rightMostDotIndex := strings.LastIndex(name, ".")
	rightMostIndex := rightMostDotIndex

	rightMostSlashIndex := strings.LastIndex(name, "/")
	if rightMostIndex < rightMostSlashIndex {
		rightMostIndex = rightMostSlashIndex
		inSlashState = true
		inDotState = false
	}

	if rightMostIndex == -1 {
		buf.WriteString(name)
		if !noPadding {
			buf.WritePadding(loggerTargetLen - inLen)
		}
		return
	}

	lastSegmentLen := inLen - rightMostIndex

	leftSegmentsTargetLen := targetLen - lastSegmentLen
	if leftSegmentsTargetLen < 0 {
		leftSegmentsTargetLen = 0
	}

	leftSegmentsLen := inLen - lastSegmentLen
	maxPossibleTrim := leftSegmentsLen - leftSegmentsTargetLen

	i := 0
	for ; i < rightMostIndex; i++ {
		c := name[i]
		if c == '.' {
			if trimmed >= maxPossibleTrim {
				break
			}
			buf.WriteByte(c)
			inDotState = true
		} else if c == '/' {
			if trimmed >= maxPossibleTrim {
				break
			}
			buf.WriteByte(c)
			inSlashState = true
		} else {
			if inDotState {
				buf.WriteByte(c)
				inDotState = false
			} else if inSlashState {
				buf.WriteByte(c)
				inSlashState = false
			} else {
				trimmed++
			}
		}
	}

	buf.WriteString(name[i:])
	end := buf.Len()
	if !noPadding {
		buf.WritePadding(loggerTargetLen - (end - start))
	}
}

func isSpecialVar(c uint8) bool {
	switch c {
	case '*', '#', '$', '@', '!', '?', '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return true
	}
	return false
}

func isAlphaNum(c uint8) bool {
	return c == '_' || '0' <= c && c <= '9' || 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z'
}
