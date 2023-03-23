package logy

import (
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type TestWriter struct {
	message string
}

func (tw *TestWriter) Write(p []byte) (n int, err error) {
	tw.message = string(p)
	return 0, nil
}

func TestRegisterHandler_PanicsIfGivenHandlerIsNil(t *testing.T) {
	assert.Panics(t, func() {
		RegisterHandler("anyHandler", nil)
	})
}

func TestRegisterHandler_PanicsIfGivenHandlerNameIsConsole(t *testing.T) {
	assert.Panics(t, func() {
		RegisterHandler("console", nil)
	})
}

func TestRegisterHandler_PanicsIfGivenHandlerNameIsFile(t *testing.T) {
	assert.Panics(t, func() {
		RegisterHandler("file", nil)
	})
}

func TestRegisterHandler_PanicsIfGivenHandlerNameIsSyslog(t *testing.T) {
	assert.Panics(t, func() {
		RegisterHandler("syslog", nil)
	})
}

func TestRegisterHandler_ShouldRegisterSuccessfully(t *testing.T) {
	RegisterHandler("anyHandler", newConsoleHandler())
	handler, ok := handlers["anyHandler"]
	assert.True(t, ok)
	assert.IsType(t, &ConsoleHandler{}, handler)
}

func TestGetHandler_ShouldReturnIfHandlerExists(t *testing.T) {
	RegisterHandler("anyHandler", newConsoleHandler())
	handler, ok := GetHandler("anyHandler")
	assert.True(t, ok)
	assert.IsType(t, &ConsoleHandler{}, handler)
}

func TestGetHandler_ShouldReturnIfHandlerDoesNotExist(t *testing.T) {
	handler, ok := GetHandler("anotherHandler")
	assert.False(t, ok)
	assert.Nil(t, handler)
}

func TestCommonHandler_Handle(t *testing.T) {
	handler := newConsoleHandler()

	testWriter := &TestWriter{}
	handler.SetColorEnabled(false)
	handler.setWriter(testWriter)

	ctx := WithValue(context.Background(), "traceId", "anyTraceId")
	ctx = WithValue(ctx, "spanId", "anySpanId")

	timestamp := time.Now()
	testCases := []struct {
		record   Record
		format   string
		json     bool
		expected string
	}{
		{
			record: Record{
				Time:       timestamp,
				Level:      LevelTrace,
				Message:    "anyMessage",
				Context:    ctx,
				LoggerName: "anyLoggerName",
				StackTrace: "anyStackTrace",
				Error:      errors.New("anyError"),
				Caller: Caller{
					defined:  true,
					file:     "main.go",
					line:     41,
					function: "TestFunction",
				},
			},
			format: "%d{2006-01-02 15:04:05.000000} %% %p %M %L %F %C %l %i %N %c : %s%e%n",
			json:   false,
			expected: fmt.Sprintf("%s %% TRACE TestFunction 41 main.go TestFunction main.go %d %s %-40s : anyMessage\nError: anyError\nanyStackTrace\n",
				timestamp.Format("2006-01-02 15:04:05.000000"),
				processId,
				processName,
				"anyLoggerName"),
		},
		{
			record: Record{
				Time:       timestamp,
				Level:      LevelTrace,
				Message:    "anyMessage",
				Context:    ctx,
				LoggerName: "anyLoggerName",
				StackTrace: "anyStackTrace",
				Error:      errors.New("anyError"),
				Caller:     Caller{},
			},
			format: "%d{2006-01-02 15:04:05.000000} %p [%x{traceId},%x{spanId}] %c : %s%e%n",
			json:   false,
			expected: fmt.Sprintf("%s TRACE [anyTraceId,anySpanId] %-40s : anyMessage\nError: anyError\nanyStackTrace\n",
				timestamp.Format("2006-01-02 15:04:05.000000"),
				"anyLoggerName"),
		},
		{
			record: Record{
				Time:       timestamp,
				Level:      LevelTrace,
				Message:    "anyMessage",
				Context:    ctx,
				LoggerName: "anyLoggerName",
				StackTrace: "anyStackTrace",
				Error:      errors.New("anyError"),
				Caller:     Caller{},
			},
			format: "%d{2006-01-02 15:04:05.000000} %p [%x{unknownKey1},%x{unknownKey12}] %c : %s%e%n",
			json:   false,
			expected: fmt.Sprintf("%s TRACE [,] %-40s : anyMessage\nError: anyError\nanyStackTrace\n",
				timestamp.Format("2006-01-02 15:04:05.000000"),
				"anyLoggerName"),
		},
		{
			record: Record{
				Time:       timestamp,
				Level:      LevelTrace,
				Message:    "anyMessage",
				Context:    ctx,
				LoggerName: "anyLoggerName",
				StackTrace: "anyStackTrace",
				Error:      errors.New("anyError"),
				Caller:     Caller{},
			},
			format: "%d{2006-01-02 15:04:05.000000} %p [%X{traceId},%X{spanId}] %c : %s%e%n",
			json:   false,
			expected: fmt.Sprintf("%s TRACE [traceId=anyTraceId,spanId=anySpanId] %-40s : anyMessage\nError: anyError\nanyStackTrace\n",
				timestamp.Format("2006-01-02 15:04:05.000000"),
				"anyLoggerName"),
		},
		{
			record: Record{
				Time:       timestamp,
				Level:      LevelTrace,
				Message:    "anyMessage",
				Context:    ctx,
				LoggerName: "anyLoggerName",
				StackTrace: "anyStackTrace",
				Error:      errors.New("anyError"),
				Caller:     Caller{},
			},
			format: "%d{2006-01-02 15:04:05.000000} %p [%X{unknownKey1},%X{unknownKey12}] %c : %s%e%n",
			json:   false,
			expected: fmt.Sprintf("%s TRACE [,] %-40s : anyMessage\nError: anyError\nanyStackTrace\n",
				timestamp.Format("2006-01-02 15:04:05.000000"),
				"anyLoggerName"),
		},
		{
			record: Record{
				Time:       timestamp,
				Level:      LevelTrace,
				Message:    "anyMessage",
				Context:    ctx,
				LoggerName: "anyLoggerName",
				StackTrace: "anyStackTrace",
				Error:      errors.New("anyError"),
				Caller:     Caller{},
			},
			format: "%d{2006-01-02 15:04:05.000000} %p [%x] %c : %s%e%n",
			json:   false,
			expected: fmt.Sprintf("%s TRACE [anyTraceId, anySpanId] %-40s : anyMessage\nError: anyError\nanyStackTrace\n",
				timestamp.Format("2006-01-02 15:04:05.000000"),
				"anyLoggerName"),
		},
		{
			record: Record{
				Time:       timestamp,
				Level:      LevelTrace,
				Message:    "anyMessage",
				Context:    context.Background(),
				LoggerName: "anyLoggerName",
				StackTrace: "anyStackTrace",
				Error:      errors.New("anyError"),
				Caller:     Caller{},
			},
			format: "%d{2006-01-02 15:04:05.000000} %p [%x] %c : %s%e%n",
			json:   false,
			expected: fmt.Sprintf("%s TRACE [] %-40s : anyMessage\nError: anyError\nanyStackTrace\n",
				timestamp.Format("2006-01-02 15:04:05.000000"),
				"anyLoggerName"),
		},
		{
			record: Record{
				Time:       timestamp,
				Level:      LevelTrace,
				Message:    "anyMessage",
				Context:    ctx,
				LoggerName: "anyLoggerName",
				StackTrace: "anyStackTrace",
				Error:      errors.New("anyError"),
				Caller:     Caller{},
			},
			format: "%d{2006-01-02 15:04:05.000000} %p [%X] %c : %s%e%n",
			json:   false,
			expected: fmt.Sprintf("%s TRACE [traceId=anyTraceId, spanId=anySpanId] %-40s : anyMessage\nError: anyError\nanyStackTrace\n",
				timestamp.Format("2006-01-02 15:04:05.000000"),
				"anyLoggerName"),
		},
		{
			record: Record{
				Time:       timestamp,
				Level:      LevelTrace,
				Message:    "anyMessage",
				Context:    context.Background(),
				LoggerName: "anyLoggerName",
				StackTrace: "anyStackTrace",
				Error:      errors.New("anyError"),
				Caller:     Caller{},
			},
			format: "%d{2006-01-02 15:04:05.000000} %p [%X] %c : %s%e%n",
			json:   false,
			expected: fmt.Sprintf("%s TRACE [] %-40s : anyMessage\nError: anyError\nanyStackTrace\n",
				timestamp.Format("2006-01-02 15:04:05.000000"),
				"anyLoggerName"),
		},
		{
			record: Record{
				Time:       timestamp,
				Level:      LevelTrace,
				Message:    "anyMessage",
				Context:    ctx,
				LoggerName: "anyLoggerName",
				StackTrace: "anyStackTrace",
				Error:      errors.New("anyError"),
				Caller:     Caller{},
			},
			json: true,
			expected: fmt.Sprintf("{"+
				"\"timestamp\":\"%s\","+
				"\"level\":\"TRACE\","+
				"\"logger\":\"anyLoggerName\","+
				"\"message\":\"anyMessage\","+
				"\"stack_trace\":\"anyStackTrace\","+
				"\"error\":\"anyError\","+
				"\"mapped_context\":{"+
				"\"traceId\":\"anyTraceId\","+
				"\"spanId\":\"anySpanId\""+
				"}"+
				"}\n", timestamp.Format(time.RFC3339)),
		},
	}

	for _, testCase := range testCases {
		handler.SetJsonEnabled(testCase.json)

		if !testCase.json {
			handler.SetFormat(testCase.format)
		}

		handler.Handle(testCase.record)
		assert.Equal(t, testCase.expected, testWriter.message)
	}

}

func TestCommonHandler_SetLevel(t *testing.T) {
	handler := newConsoleHandler()
	handler.SetLevel(LevelTrace)
	assert.Equal(t, LevelTrace, handler.level.Load().(Level))
}

func TestCommonHandler_Level(t *testing.T) {
	handler := newConsoleHandler()
	handler.SetLevel(LevelTrace)
	assert.Equal(t, LevelTrace, handler.Level())
}

func TestCommonHandler_SetEnabled(t *testing.T) {
	handler := newConsoleHandler()
	handler.SetEnabled(true)
	assert.True(t, handler.enabled.Load().(bool))

	handler.SetEnabled(false)
	assert.False(t, handler.enabled.Load().(bool))
}

func TestCommonHandler_IsEnabled(t *testing.T) {
	handler := newConsoleHandler()
	handler.SetEnabled(true)
	assert.True(t, handler.IsEnabled())

	handler.SetEnabled(false)
	assert.False(t, handler.IsEnabled())
}

func TestCommonHandler_IsLoggable(t *testing.T) {
	handler := newConsoleHandler()
	handler.SetEnabled(false)
	assert.False(t, handler.IsLoggable(Record{}))

	handler.SetEnabled(true)
	handler.SetLevel(LevelOff)
	assert.False(t, handler.IsLoggable(Record{
		Level: LevelTrace,
	}))
}

func TestCommonHandler_Writer(t *testing.T) {
	testWriter := newSyncWriter(&discarder{})
	handler := newConsoleHandler()
	handler.setWriter(testWriter)
	assert.Equal(t, testWriter, handler.Writer())
}

func TestCommonHandler_SetFormat(t *testing.T) {
	handler := newConsoleHandler()
	handler.SetFormat("anyFormat")
	assert.Equal(t, "anyFormat", handler.format.Load().(string))
}

func TestCommonHandler_Format(t *testing.T) {
	handler := newConsoleHandler()
	handler.SetFormat("anyFormat")
	assert.Equal(t, "anyFormat", handler.Format())
}

func TestCommonHandler_SetJsonEnabled(t *testing.T) {
	handler := newConsoleHandler()
	handler.SetJsonEnabled(true)
	assert.True(t, handler.json.Load().(bool))

	handler.SetJsonEnabled(false)
	assert.False(t, handler.json.Load().(bool))
}

func TestCommonHandler_IsJsonEnabled(t *testing.T) {
	handler := newConsoleHandler()
	handler.SetJsonEnabled(true)
	assert.True(t, handler.IsJsonEnabled())

	handler.SetJsonEnabled(false)
	assert.False(t, handler.IsJsonEnabled())
}

func TestCommonHandler_SetExcludedKeys(t *testing.T) {
	handler := newConsoleHandler()
	excludedKeys := ExcludedKeys{
		"timestamp", "anyKey1",
	}
	handler.SetExcludedKeys(excludedKeys)
	assert.Equal(t, allKeysEnabled^timestampKeyEnabled, handler.enabledKeys.Load().(int))
	assert.Equal(t, ExcludedKeys{"anyKey1"}, handler.excludedFields.Load().(ExcludedKeys))
}

func TestCommonHandler_ExcludedKeys(t *testing.T) {
	handler := newConsoleHandler()
	excludedKeys := ExcludedKeys{
		"timestamp", "mapped_context", "level", "logger", "message", "error", "stack_trace", "anyKey1",
	}
	handler.SetExcludedKeys(excludedKeys)
	assert.Equal(t, 0, handler.enabledKeys.Load().(int))
	assert.Equal(t, excludedKeys, handler.ExcludedKeys())
}

func TestCommonHandler_SetKeyOverrides(t *testing.T) {
	handler := newConsoleHandler()
	keyOverrides := KeyOverrides{
		"timestamp":      "@timestamp",
		"logger":         "@logger",
		"level":          "@level",
		"mapped_context": "@mapped_context",
		"message":        "@message",
		"error":          "@error",
		"stack_trace":    "@stack_trace",
		"anyKey1":        "anyValue1",
	}
	handler.SetKeyOverrides(keyOverrides)
	assert.Equal(t, keyOverrides, handler.keyOverrides.Load().(KeyOverrides))
}

func TestCommonHandler_KeyOverrides(t *testing.T) {
	handler := newConsoleHandler()
	keyOverrides := KeyOverrides{
		"anyKey1": "anyValue1",
	}
	handler.SetKeyOverrides(keyOverrides)
	assert.Equal(t, keyOverrides, handler.KeyOverrides())
}

func TestCommonHandler_SetAdditionalFields(t *testing.T) {
	handler := newConsoleHandler()
	additionalFields := AdditionalFields{
		"anyKey1": "anyValue1",
	}
	handler.SetAdditionalFields(additionalFields)
	assert.Equal(t, "\"anyKey1\":\"anyValue1\"", handler.additionalFieldsJson.Load().(string))
	assert.Equal(t, additionalFields, handler.additionalFields.Load().(AdditionalFields))
}

func TestCommonHandler_AdditionalFields(t *testing.T) {
	handler := newConsoleHandler()
	additionalFields := AdditionalFields{
		"anyKey1": "anyValue1",
	}
	handler.SetAdditionalFields(additionalFields)
	assert.Equal(t, additionalFields, handler.AdditionalFields())
}
