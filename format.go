package logy

import (
	"strings"
)

const (
	loggerTargetLen = 40
)

var (
	timestampKey     = "timestamp"
	mappedContextKey = "mappedContext"
	levelKey         = "level"
	loggerKey        = "logger"
	messageKey       = "message"
	stackTraceKey    = "stack_trace"
)

func formatText(encoder *textEncoder, format string, record Record, console bool) {
	mc := MappedContextFrom(record.Context)

	i := 0
	for j := 0; j < len(format); j++ {
		if format[j] == '%' && j+1 < len(format) {
			typ := format[j+1]
			w := 1

			switch typ {
			case 'd': // date
				encoder.AppendTime(record.Time)
			case 'c': // logger
				appendLoggerAsText(encoder.buf, record.LoggerName, console)
			case 'p': // level
				appendLevelAsText(encoder.buf, record.Level, console)
			case 'x': // context value without key
				name, l := getPlaceholderName(format[j+2:])

				if mc != nil && name != "" {
					val := mc.value(name)
					if val != nil {
						encoder.AppendAny(val)
					}
				}

				w = l + 1
			case 'X': // context value with key
				name, l := getPlaceholderName(format[j+2:])

				if mc != nil && name != "" {
					encoder.AppendString(name)
					encoder.buf.WriteByte('=')
					val := mc.value(name)
					if val != nil {
						encoder.AppendAny(val)
					}
				}

				w = l + 1
			case 'm': // message
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
			case 's': // stack trace if exist
				if record.StackTrace != "" {
					encoder.buf.WriteByte('\n')
					encoder.buf.WriteString(strings.ReplaceAll(record.StackTrace, "\\n", "\n"))
				}
			case 'n': // newline
				encoder.buf.WriteByte('\n')
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

func appendLoggerAsText(buf *buffer, logger string, console bool) {
	loggerName := abbreviateLoggerName(logger, loggerTargetLen)

	if console {
		colorCyan.print(buf, loggerName)
	} else {
		buf.WriteString(loggerName)
	}

	buf.WritePadding(loggerTargetLen - len(loggerName))
}

func appendLevelAsText(buf *buffer, level Level, console bool) {
	str := level.String()
	buf.WritePadding(5 - len(str))

	if console {
		levelColors[level-1].print(buf, str)
	} else {
		buf.WriteString(str)
	}
}

func formatJson(encoder *jsonEncoder, record Record, additionalFieldJson string) {
	// timestamp
	encoder.AddTime(timestampKey, record.Time)
	// level
	encoder.AddString(loggerKey, record.Level.String())

	// logger name
	encoder.AddString(loggerKey, record.LoggerName)

	// message
	encoder.AddString(messageKey, record.Message)

	if record.StackTrace != "" {
		// stack trace
		encoder.AddString(stackTraceKey, record.StackTrace)
	}

	// mapped context
	if record.Context != nil {
		mc := MappedContextFrom(record.Context)
		encoder.addKey(mappedContextKey)
		encoder.buf.WriteString(mc.ValuesAsJSON(nil))
	}

	// additional fields
	if len(additionalFieldJson) != 0 {
		encoder.buf.WriteByte(',')
		encoder.buf.WriteString(additionalFieldJson)
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

func abbreviateLoggerName(name string, targetLen int) string {
	inLen := len(name)
	if inLen < targetLen {
		return name
	}

	buf := newBuffer()
	defer buf.Free()

	rightMostDotIndex := strings.LastIndex(name, ".")

	if rightMostDotIndex == -1 {
		return name
	}

	lastSegmentLen := inLen - rightMostDotIndex

	leftSegmentsTargetLen := targetLen - lastSegmentLen
	if leftSegmentsTargetLen < 0 {
		leftSegmentsTargetLen = 0
	}

	leftSegmentsLen := inLen - lastSegmentLen
	maxPossibleTrim := leftSegmentsLen - leftSegmentsTargetLen

	trimmed := 0
	inDotState := true

	i := 0
	for ; i < rightMostDotIndex; i++ {
		c := name[i]
		if c == '.' {
			if trimmed >= maxPossibleTrim {
				break
			}
			buf.WriteByte(c)
			inDotState = true
		} else {
			if inDotState {
				buf.WriteByte(c)
				inDotState = false
			} else {
				trimmed++
			}
		}
	}

	buf.WriteString(name[i:])
	return buf.String()
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
