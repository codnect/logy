package logy

import (
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSyslogHandler_Handle(t *testing.T) {
	handler := newSysLogHandler(true)

	testWriter := &TestWriter{}
	handler.setWriter(testWriter)

	ctx := WithValue(context.Background(), "traceId", "anyTraceId")
	ctx = WithValue(ctx, "spanId", "anySpanId")

	timestamp := time.Now()
	testCases := []struct {
		record            Record
		format            string
		syslogType        SysLogType
		expectedInRFC5424 string
		expectedInRFC3164 string
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
					file:     filepath.FromSlash("/test/any/main.go"),
					line:     41,
					function: "TestFunction",
				},
			},
			format: "%d %% %p %M %L %F %C %l %i %N %c : %s%e%n",
			expectedInRFC5424: fmt.Sprintf("<13>1 %s - %s "+
				"%d anyLoggerName - \ufeff%s %% TRACE TestFunction 41 main.go TestFunction /test/any %d %s %s : anyMessage\nError: anyError\nanyStackTrace\n",
				timestamp.Format(time.RFC3339),
				os.Args[0],
				processId,
				timestamp.Format(time.RFC3339),
				processId,
				processName,
				"anyLoggerName"),
			expectedInRFC3164: fmt.Sprintf("<13>%s UNKNOWN_HOSTNAME [%d]: %s %% TRACE TestFunction 41 main.go TestFunction /test/any %d %s %s : anyMessage\nError: anyError\nanyStackTrace\n",
				timestamp.Format(time.Stamp),
				processId,
				timestamp.Format(time.RFC3339),
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
				Caller: Caller{
					defined:  true,
					file:     filepath.FromSlash("/test/any/main.go"),
					line:     41,
					function: "TestFunction",
				},
			},
			format: "%d %p [%x{traceId},%x{spanId}] %c : %s%e%n",
			expectedInRFC5424: fmt.Sprintf("<13>1 %s - %s "+
				"%d anyLoggerName - \ufeff%s TRACE [anyTraceId,anySpanId] %s : anyMessage\nError: anyError\nanyStackTrace\n",
				timestamp.Format(time.RFC3339),
				os.Args[0],
				processId,
				timestamp.Format(time.RFC3339),
				"anyLoggerName"),
			expectedInRFC3164: fmt.Sprintf("<13>%s UNKNOWN_HOSTNAME [%d]: %s TRACE [anyTraceId,anySpanId] %s : anyMessage\nError: anyError\nanyStackTrace\n",
				timestamp.Format(time.Stamp),
				processId,
				timestamp.Format(time.RFC3339),
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
				Caller: Caller{
					defined:  true,
					file:     filepath.FromSlash("/test/any/main.go"),
					line:     41,
					function: "TestFunction",
				},
			},
			format: "%d{2006-01-02 15:04:05.000000} %p [%x{unknownKey1},%x{unknownKey12}] %c : %s%e%n",
			expectedInRFC5424: fmt.Sprintf("<13>1 %s - %s "+
				"%d anyLoggerName - \ufeff%s TRACE [,] %s : anyMessage\nError: anyError\nanyStackTrace\n",
				timestamp.Format(time.RFC3339),
				os.Args[0],
				processId,
				timestamp.Format("2006-01-02 15:04:05.000000"),
				"anyLoggerName"),
			expectedInRFC3164: fmt.Sprintf("<13>%s UNKNOWN_HOSTNAME [%d]: %s TRACE [,] %s : anyMessage\nError: anyError\nanyStackTrace\n",
				timestamp.Format(time.Stamp),
				processId,
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
				Caller: Caller{
					defined:  true,
					file:     filepath.FromSlash("/test/any/main.go"),
					line:     41,
					function: "TestFunction",
				},
			},
			format: "%d{2006-01-02 15:04:05.000000} %p [%X{traceId},%X{spanId}] %c : %s%e%n",
			expectedInRFC5424: fmt.Sprintf("<13>1 %s - %s "+
				"%d anyLoggerName - \ufeff%s TRACE [traceId=anyTraceId,spanId=anySpanId] %s : anyMessage\nError: anyError\nanyStackTrace\n",
				timestamp.Format(time.RFC3339),
				os.Args[0],
				processId,
				timestamp.Format("2006-01-02 15:04:05.000000"),
				"anyLoggerName"),
			expectedInRFC3164: fmt.Sprintf("<13>%s UNKNOWN_HOSTNAME [%d]: %s TRACE [traceId=anyTraceId,spanId=anySpanId] %s : anyMessage\nError: anyError\nanyStackTrace\n",
				timestamp.Format(time.Stamp),
				processId,
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
				Caller: Caller{
					defined:  true,
					file:     filepath.FromSlash("/test/any/main.go"),
					line:     41,
					function: "TestFunction",
				},
			},
			format: "%d{2006-01-02 15:04:05.000000} %p [%X{unknownKey1},%X{unknownKey12}] %c : %s%e%n",
			expectedInRFC5424: fmt.Sprintf("<13>1 %s - %s "+
				"%d anyLoggerName - \ufeff%s TRACE [,] %s : anyMessage\nError: anyError\nanyStackTrace\n",
				timestamp.Format(time.RFC3339),
				os.Args[0],
				processId,
				timestamp.Format("2006-01-02 15:04:05.000000"),
				"anyLoggerName"),
			expectedInRFC3164: fmt.Sprintf("<13>%s UNKNOWN_HOSTNAME [%d]: %s TRACE [,] %s : anyMessage\nError: anyError\nanyStackTrace\n",
				timestamp.Format(time.Stamp),
				processId,
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
				Caller: Caller{
					defined:  true,
					file:     filepath.FromSlash("/test/any/main.go"),
					line:     41,
					function: "TestFunction",
				},
			},
			format: "%d{2006-01-02 15:04:05.000000} %p [%x] %c : %s%e%n",
			expectedInRFC5424: fmt.Sprintf("<13>1 %s - %s "+
				"%d anyLoggerName - \ufeff%s TRACE [anyTraceId, anySpanId] %s : anyMessage\nError: anyError\nanyStackTrace\n",
				timestamp.Format(time.RFC3339),
				os.Args[0],
				processId,
				timestamp.Format("2006-01-02 15:04:05.000000"),
				"anyLoggerName"),
			expectedInRFC3164: fmt.Sprintf("<13>%s UNKNOWN_HOSTNAME [%d]: %s TRACE [anyTraceId, anySpanId] %s : anyMessage\nError: anyError\nanyStackTrace\n",
				timestamp.Format(time.Stamp),
				processId,
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
				Caller: Caller{
					defined:  true,
					file:     filepath.FromSlash("/test/any/main.go"),
					line:     41,
					function: "TestFunction",
				},
			},
			format: "%d{2006-01-02 15:04:05.000000} %p [%x] %c : %s%e%n",
			expectedInRFC5424: fmt.Sprintf("<13>1 %s - %s "+
				"%d anyLoggerName - \ufeff%s TRACE [] %s : anyMessage\nError: anyError\nanyStackTrace\n",
				timestamp.Format(time.RFC3339),
				os.Args[0],
				processId,
				timestamp.Format("2006-01-02 15:04:05.000000"),
				"anyLoggerName"),
			expectedInRFC3164: fmt.Sprintf("<13>%s UNKNOWN_HOSTNAME [%d]: %s TRACE [] %s : anyMessage\nError: anyError\nanyStackTrace\n",
				timestamp.Format(time.Stamp),
				processId,
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
				Caller: Caller{
					defined:  true,
					file:     filepath.FromSlash("/test/any/main.go"),
					line:     41,
					function: "TestFunction",
				},
			},
			format: "%d{2006-01-02 15:04:05.000000} %p [%X] %c : %s%e%n",
			expectedInRFC5424: fmt.Sprintf("<13>1 %s - %s "+
				"%d anyLoggerName - \ufeff%s TRACE [traceId=anyTraceId, spanId=anySpanId] %s : anyMessage\nError: anyError\nanyStackTrace\n",
				timestamp.Format(time.RFC3339),
				os.Args[0],
				processId,
				timestamp.Format("2006-01-02 15:04:05.000000"),
				"anyLoggerName"),
			expectedInRFC3164: fmt.Sprintf("<13>%s UNKNOWN_HOSTNAME [%d]: %s TRACE [traceId=anyTraceId, spanId=anySpanId] %s : anyMessage\nError: anyError\nanyStackTrace\n",
				timestamp.Format(time.Stamp),
				processId,
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
				Caller: Caller{
					defined:  true,
					file:     filepath.FromSlash("/test/any/main.go"),
					line:     41,
					function: "TestFunction",
				},
			},
			format: "%d{2006-01-02 15:04:05.000000} %p [%X] %c : %s%e%n",
			expectedInRFC5424: fmt.Sprintf("<13>1 %s - %s "+
				"%d anyLoggerName - \ufeff%s TRACE [] %s : anyMessage\nError: anyError\nanyStackTrace\n",
				timestamp.Format(time.RFC3339),
				os.Args[0],
				processId,
				timestamp.Format("2006-01-02 15:04:05.000000"),
				"anyLoggerName"),
			expectedInRFC3164: fmt.Sprintf("<13>%s UNKNOWN_HOSTNAME [%d]: %s TRACE [] %s : anyMessage\nError: anyError\nanyStackTrace\n",
				timestamp.Format(time.Stamp),
				processId,
				timestamp.Format("2006-01-02 15:04:05.000000"),
				"anyLoggerName"),
		},
	}

	for _, testCase := range testCases {
		handler.SetFormat(testCase.format)
		handler.setLogType(RFC5424)

		handler.Handle(testCase.record)
		assert.Equal(t, testCase.expectedInRFC5424, testWriter.message)

		handler.setLogType(RFC3164)
		handler.Handle(testCase.record)
		assert.Equal(t, testCase.expectedInRFC3164, testWriter.message)
	}

}

func TestSyslogHandler_Endpoint(t *testing.T) {
	handler := newSysLogHandler(true)
	handler.setEndpoint("anyEndpoint")
	assert.Equal(t, "anyEndpoint", handler.Endpoint())
}

func TestSyslogHandler_SetApplicationName(t *testing.T) {
	handler := newSysLogHandler(true)
	handler.SetApplicationName("anyApplicationName")
	assert.Equal(t, "anyApplicationName", handler.appName.Load().(string))
}

func TestSyslogHandler_ApplicationName(t *testing.T) {
	handler := newSysLogHandler(true)
	handler.SetApplicationName("anyApplicationName")
	assert.Equal(t, "anyApplicationName", handler.ApplicationName())
}

func TestSyslogHandler_SetHostname(t *testing.T) {
	handler := newSysLogHandler(true)
	handler.SetHostname("anyHostname")
	assert.Equal(t, "anyHostname", handler.hostname.Load().(string))
}

func TestSyslogHandler_Hostname(t *testing.T) {
	handler := newSysLogHandler(true)
	handler.SetHostname("anyHostname")
	assert.Equal(t, "anyHostname", handler.Hostname())
}

func TestSyslogHandler_SetFacility(t *testing.T) {
	handler := newSysLogHandler(false)
	handler.SetFacility(FacilityLogAlert)
	assert.Equal(t, FacilityLogAlert, handler.facility.Load().(Facility))
}

func TestSyslogHandler_Facility(t *testing.T) {
	handler := newSysLogHandler(true)
	handler.SetFacility(FacilityLogAlert)
	assert.Equal(t, FacilityLogAlert, handler.Facility())
}

func TestSyslogHandler_LogType(t *testing.T) {
	handler := newSysLogHandler(true)
	handler.setLogType(RFC5424)
	assert.Equal(t, RFC5424, handler.LogType())
}

func TestSyslogHandler_Protocol(t *testing.T) {
	handler := newSysLogHandler(true)
	handler.setProtocol(ProtocolTCP)
	assert.Equal(t, ProtocolTCP, handler.protocol.Load().(Protocol))
}

func TestSyslogHandler_IsBlockOnReconnect(t *testing.T) {
	handler := newSysLogHandler(true)
	handler.setBlockOnReconnect(true)
	assert.True(t, handler.IsBlockOnReconnect())

	handler.setBlockOnReconnect(false)
	assert.False(t, handler.IsBlockOnReconnect())
}

func TestSyslogHandler_OnConfigure(t *testing.T) {

}
