package logy

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	loggerTargetLen = 40
)

var (
	timestampKey     = []byte("timestamp")
	mappedContextKey = []byte("mappedContext")
	levelKey         = []byte("level")
	loggerKey        = []byte("logger")
	messageKey       = []byte("message")
	stackTraceKey    = []byte("stack_trace")
)

func appendValue(buf *[]byte, val any, json bool) {
	switch typed := val.(type) {
	case string:
		if json {
			*buf = append(*buf, '"')
		}

		*buf = append(*buf, typed...)

		if json {
			*buf = append(*buf, '"')
		}
	case error:
		if json {
			*buf = append(*buf, '"')
		}

		*buf = append(*buf, typed.Error()...)

		if json {
			*buf = append(*buf, '"')
		}
	case int:
		appendInt(buf, int64(typed))
	case int8:
		appendInt(buf, int64(typed))
	case int16:
		appendInt(buf, int64(typed))
	case int32:
		appendInt(buf, int64(typed))
	case int64:
		appendInt(buf, typed)
	case uint:
		appendUInt(buf, uint64(typed))
	case uint8:
		appendUInt(buf, uint64(typed))
	case uint16:
		appendUInt(buf, uint64(typed))
	case uint32:
		appendUInt(buf, uint64(typed))
	case uint64:
		appendUInt(buf, typed)
	case bool:
		appendBool(buf, typed)
	case float32:
		appendFloat32(buf, typed)
	case float64:
		appendFloat64(buf, typed)
	case time.Duration:
		appendInt(buf, int64(typed))
	case time.Time:
		formatDateAsText(buf, typed)
	default:
		rValue := reflect.ValueOf(typed)

		if rValue.Kind() == reflect.Map {
			appendMap(buf, rValue, json)
		} else if rValue.Kind() == reflect.Array || rValue.Kind() == reflect.Slice {
			appendSlice(buf, rValue, json)
		} else {
			if stringer, implements := typed.(fmt.Stringer); implements {
				appendValue(buf, stringer.String(), json)
			} else {
				appendValue(buf, rValue.String(), json)
			}
		}
	}
}

func appendInt(buf *[]byte, val int64) {
	*buf = append(*buf, strconv.FormatInt(val, 10)...)
}

func appendUInt(buf *[]byte, val uint64) {
	*buf = append(*buf, strconv.FormatUint(val, 10)...)
}

func appendBool(buf *[]byte, val bool) {
	*buf = append(*buf, strconv.FormatBool(val)...)
}

func appendFloat32(buf *[]byte, val float32) {
	*buf = append(*buf, strconv.FormatFloat(float64(val), 'f', 6, 32)...)
}

func appendFloat64(buf *[]byte, val float64) {
	*buf = append(*buf, strconv.FormatFloat(val, 'f', 6, 32)...)
}

func appendSlice(buf *[]byte, rtype reflect.Value, json bool) {
	*buf = append(*buf, '[')

	for i := 0; i < rtype.Len(); i++ {
		item := rtype.Index(i)

		appendValue(buf, item.Interface(), json)

		if rtype.Len()-1 != i {
			*buf = append(*buf, ',')
		}
	}

	*buf = append(*buf, ']')
}

func appendMap(buf *[]byte, rtype reflect.Value, json bool) {
	*buf = append(*buf, '{')
	iter := rtype.MapRange()
	i := 0

	for iter.Next() {
		appendValue(buf, iter.Key().Interface(), json)

		if json {
			*buf = append(*buf, ':')
		} else {
			*buf = append(*buf, '=')
		}

		appendValue(buf, iter.Value().Interface(), json)

		i++

		if rtype.Len() != i {
			*buf = append(*buf, ',')
		}
	}

	*buf = append(*buf, '}')
}

func formatText(buf *[]byte, format string, record Record, console bool) {
	mc := MappedContextFrom(record.Context)

	i := 0
	for j := 0; j < len(format); j++ {
		if format[j] == '%' && j+1 < len(format) {
			typ := format[j+1]
			w := 1

			switch typ {
			case 'd': // date
				formatDateAsText(buf, record.Time)
			case 'c': // logger
				formatLoggerAsText(buf, record.LoggerName, console)
			case 'p': // level
				formatLevelAsText(buf, record.Level, console)
			case 'x': // context value without key
				name, l := getPlaceholderName(format[j+2:])
				formatContextValueAsText(buf, name, mc, false)
				w = l + 1
			case 'X': // context value with key
				name, l := getPlaceholderName(format[j+2:])
				formatContextValueAsText(buf, name, mc, true)
				w = l + 1
			case 'm': // message
				*buf = append(*buf, record.Message...)
			case 'M': // method
				*buf = append(*buf, record.Caller.Name()...)
			case 'L': // line
				appendInt(buf, int64(record.Caller.Line()))
			case 'F': // file
				*buf = append(*buf, record.Caller.File()...)
			case 'C': // package
				*buf = append(*buf, record.Caller.Package()...)
			case 'l': // location
				*buf = append(*buf, record.Caller.Path()...)
			case 's': // stack trace if exist
				if record.StackTrace != "" {
					*buf = append(*buf, '\n')
					*buf = append(*buf, strings.ReplaceAll(record.StackTrace, "\\n", "\n")...)
				}
			case 'n': // newline
				*buf = append(*buf, '\n')
			default:
				*buf = append(*buf, format[i:j]...)
			}

			j += w
			i = j + 1
		} else {
			*buf = append(*buf, format[j])
			i = j + 1
		}
	}
}

func formatDateAsText(buf *[]byte, t time.Time) {
	year, month, day := t.Date()
	itoa(buf, year, 4)
	*buf = append(*buf, '-')
	itoa(buf, int(month), 2)
	*buf = append(*buf, '-')
	itoa(buf, day, 2)
	*buf = append(*buf, ' ')

	hour, min, sec := t.Clock()
	itoa(buf, hour, 2)
	*buf = append(*buf, ':')
	itoa(buf, min, 2)
	*buf = append(*buf, ':')
	itoa(buf, sec, 2)

	*buf = append(*buf, '.')
	itoa(buf, t.Nanosecond()/1e3, 6)
}

func formatLoggerAsText(buf *[]byte, logger string, console bool) {
	loggerName := abbreviateLoggerName(logger, loggerTargetLen)

	if console {
		colorCyan.print(buf, loggerName)
	} else {
		*buf = append(*buf, loggerName...)
	}

	appendPadding(buf, loggerTargetLen-len(loggerName))
}

func formatLevelAsText(buf *[]byte, level Level, console bool) {
	str := level.String()
	appendPadding(buf, 5-len(str))

	if console {
		levelColors[level-1].print(buf, str)
	} else {
		*buf = append(*buf, str...)
	}
}

func formatContextValueAsText(buf *[]byte, key string, mc *MappedContext, includeKey bool) {
	if mc == nil {
		return
	}

	if key == "" {
		return
	}

	if includeKey {
		*buf = append(*buf, key...)
		*buf = append(*buf, '=')
	}

	val := mc.Value(key)

	if val != nil {
		appendValue(buf, val, false)
	}
}

func appendJsonDateField(buf *[]byte, key []byte, t time.Time) {
	*buf = append(*buf, '"')
	*buf = append(*buf, key...)
	*buf = append(*buf, '"')
	*buf = append(*buf, ':')
	*buf = append(*buf, '"')

	year, month, day := t.Date()
	itoa(buf, year, 4)
	*buf = append(*buf, '-')
	itoa(buf, int(month), 2)
	*buf = append(*buf, '-')
	itoa(buf, day, 2)
	*buf = append(*buf, 'T')

	hour, min, sec := t.Clock()
	itoa(buf, hour, 2)
	*buf = append(*buf, ':')
	itoa(buf, min, 2)
	*buf = append(*buf, ':')
	itoa(buf, sec, 2)

	*buf = append(*buf, '.')
	itoa(buf, t.Nanosecond()/1e3, 6)
	*buf = append(*buf, '"')
}

func appendJsonStringField(buf *[]byte, key []byte, value string) {
	*buf = append(*buf, '"')
	*buf = append(*buf, key...)
	*buf = append(*buf, '"')
	*buf = append(*buf, ':')
	*buf = append(*buf, '"')
	*buf = append(*buf, value...)
	*buf = append(*buf, '"')
}

func appendJsonAnyField(buf *[]byte, key string, val any) {
	*buf = append(*buf, '"')
	*buf = append(*buf, key...)
	*buf = append(*buf, '"')
	*buf = append(*buf, ':')

	appendValue(buf, val, true)
}

func appendContextValues(buf *[]byte, ctx context.Context, excludedKeys map[string]struct{}) {
	*buf = append(*buf, '"')
	*buf = append(*buf, mappedContextKey...)
	*buf = append(*buf, '"')
	*buf = append(*buf, ':')
	*buf = append(*buf, '{')

	if ctx != nil {
		mc := MappedContextFrom(ctx)

		if mc != nil {
			keys := mc.Keys()

			for i, key := range keys {
				if _, ok := excludedKeys[key]; ok {
					continue
				}

				val := mc.Value(key)

				if val != nil {
					appendJsonAnyField(buf, key, val)

					if len(keys)-1 != i {
						*buf = append(*buf, ',')
					}
				}
			}
		}
	}

	*buf = append(*buf, '}')
}

func appendAdditionalFields(buf *[]byte, additionalFields map[string]JsonAdditionalField) {
	i := 0

	for key, field := range additionalFields {
		appendJsonAnyField(buf, key, field.Value)

		if len(additionalFields)-1 != i {
			*buf = append(*buf, ',')
		}

		i++
	}
}

func formatJson(buf *[]byte, record Record, excludedKeys map[string]struct{}, additionalFields map[string]JsonAdditionalField) {
	*buf = append(*buf, '{')

	// timestamp
	appendJsonDateField(buf, timestampKey, record.Time)
	*buf = append(*buf, ',')

	// level
	appendJsonStringField(buf, levelKey, record.Level.String())
	*buf = append(*buf, ',')

	// logger name
	appendJsonStringField(buf, loggerKey, record.LoggerName)
	*buf = append(*buf, ',')

	// message
	appendJsonStringField(buf, messageKey, record.Message)

	*buf = append(*buf, ',')

	// mapped context
	appendContextValues(buf, record.Context, excludedKeys)

	if record.StackTrace != "" {
		// stack trace
		appendJsonStringField(buf, stackTraceKey, record.StackTrace)
	}

	// additional fields
	if len(additionalFields) != 0 {
		*buf = append(*buf, ',')
		appendAdditionalFields(buf, additionalFields)
	}

	*buf = append(*buf, '}')
	*buf = append(*buf, '\n')
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

func appendPadding(buf *[]byte, n int) {
	if n <= 0 {
		return
	}

	for i := 0; i < n; i++ {
		*buf = append(*buf, ' ')
	}
}

func abbreviateLoggerName(name string, targetLen int) string {
	inLen := len(name)
	if inLen < targetLen {
		return name
	}

	var buf []byte
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
			buf = append(buf, c)
			inDotState = true
		} else {
			if inDotState {
				buf = append(buf, c)
				inDotState = false
			} else {
				trimmed++
			}
		}
	}

	buf = append(buf, name[i:]...)
	return string(buf)
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

func itoa(buf *[]byte, i int, wid int) {
	var b [20]byte
	bp := len(b) - 1

	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}

	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}
