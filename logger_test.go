package logy

import (
	"context"
	"errors"
	"github.com/stretchr/testify/mock"
	"io"
	"runtime"
	"testing"
	"time"
)

type testHandler struct {
	mock.Mock
}

func (h *testHandler) Handle(record Record) error {
	arg := h.Called(record)

	if len(arg) == 1 {
		err, ok := arg[0].(error)
		if ok {
			return err
		}
	}

	return nil
}

func (h *testHandler) SetLevel(level Level) {
}

func (h *testHandler) Level() Level {
	return LevelAll
}

func (h *testHandler) SetEnabled(enabled bool) {

}

func (h *testHandler) IsEnabled() bool {
	return true
}

func (h *testHandler) IsLoggable(record Record) bool {
	return true
}

func (h *testHandler) Writer() io.Writer {
	return nil
}

var (
	mockTestHandler   = &testHandler{}
	_, filename, _, _ = runtime.Caller(0)
)

func init() {
	RegisterHandler("testHandler", mockTestHandler)
	LoadConfig(&Config{
		Level:         LevelAll,
		IncludeCaller: true,
		Handlers:      Handlers{"testHandler"},
	})
}

func TestLogger_I(t *testing.T) {
	current := time.Now()
	now = func() time.Time {
		return current
	}

	ctx := WithValue(context.Background(), "anyKey", "anyValue")

	mockTestHandler.On("Handle", mock.AnythingOfType("Record")).Return(nil)

	logger := Named("anyLogger")
	logger.I(ctx, "anyMessage")

	mockTestHandler.AssertCalled(t, "Handle", Record{
		Time:       current,
		Level:      LevelInfo,
		Message:    "anyMessage",
		Context:    ctx,
		LoggerName: "anyLogger",
		Caller: Caller{
			defined:  true,
			file:     filename,
			line:     78,
			function: "github.com/procyon-projects/logy.TestLogger_I",
		},
	})
}

func TestLogger_IWithStackTrace(t *testing.T) {
	current := time.Now()
	now = func() time.Time {
		return current
	}
	err := errors.New("anyError")

	ctx := WithValue(context.Background(), "anyKey", "anyValue")

	mockTestHandler.On("Handle", mock.AnythingOfType("Record")).Return(nil)

	logger := Named("anyLogger")
	logger.I(ctx, "anyMessage", err)

	mockTestHandler.AssertCalled(t, "Handle", Record{
		Time:       current,
		Level:      LevelInfo,
		Message:    "anyMessage",
		Context:    ctx,
		LoggerName: "anyLogger",
		Error:      err,
		StackTrace: "github.com/procyon-projects/logy.TestLogger_IWithStackTrace()\\n    /Users/burakkoken/GolandProjects/slog/logger_test.go:107",
		Caller: Caller{
			defined:  true,
			file:     filename,
			line:     107,
			function: "github.com/procyon-projects/logy.TestLogger_IWithStackTrace",
		},
	})
}

func TestLogger_IWithArguments(t *testing.T) {
	current := time.Now()
	now = func() time.Time {
		return current
	}

	ctx := WithValue(context.Background(), "anyKey", "anyValue")

	mockTestHandler.On("Handle", mock.AnythingOfType("Record")).Return(nil)

	logger := Named("anyLogger")
	logger.I(ctx, "anyMessage {}", "anyValue")

	mockTestHandler.AssertCalled(t, "Handle", Record{
		Time:       current,
		Level:      LevelInfo,
		Message:    "anyMessage anyValue",
		Context:    ctx,
		LoggerName: "anyLogger",
		Caller: Caller{
			defined:  true,
			file:     filename,
			line:     137,
			function: "github.com/procyon-projects/logy.TestLogger_IWithArguments",
		},
	})
}

func TestLogger_Info(t *testing.T) {
	current := time.Now()
	now = func() time.Time {
		return current
	}

	mockTestHandler.On("Handle", mock.AnythingOfType("Record")).Return(nil)

	logger := Named("anyLogger")
	logger.Info("anyMessage")

	mockTestHandler.AssertCalled(t, "Handle", Record{
		Time:       current,
		Level:      LevelInfo,
		Message:    "anyMessage",
		LoggerName: "anyLogger",
		Caller: Caller{
			defined:  true,
			file:     filename,
			line:     163,
			function: "github.com/procyon-projects/logy.TestLogger_Info",
		},
	})
}

func TestLogger_InfoWithStackTrace(t *testing.T) {
	current := time.Now()
	now = func() time.Time {
		return current
	}
	err := errors.New("anyError")

	mockTestHandler.On("Handle", mock.AnythingOfType("Record")).Return(nil)

	logger := Named("anyLogger")
	logger.Info("anyMessage", err)

	mockTestHandler.AssertCalled(t, "Handle", Record{
		Time:       current,
		Level:      LevelInfo,
		Message:    "anyMessage",
		LoggerName: "anyLogger",
		Error:      err,
		StackTrace: "github.com/procyon-projects/logy.TestLogger_InfoWithStackTrace()\\n    /Users/burakkoken/GolandProjects/slog/logger_test.go:189",
		Caller: Caller{
			defined:  true,
			file:     filename,
			line:     189,
			function: "github.com/procyon-projects/logy.TestLogger_InfoWithStackTrace",
		},
	})
}

func TestLogger_InfoWithArguments(t *testing.T) {
	current := time.Now()
	now = func() time.Time {
		return current
	}

	mockTestHandler.On("Handle", mock.AnythingOfType("Record")).Return(nil)

	logger := Named("anyLogger")
	logger.Info("anyMessage {}", "anyValue")

	mockTestHandler.AssertCalled(t, "Handle", Record{
		Time:       current,
		Level:      LevelInfo,
		Message:    "anyMessage anyValue",
		LoggerName: "anyLogger",
		Caller: Caller{
			defined:  true,
			file:     filename,
			line:     216,
			function: "github.com/procyon-projects/logy.TestLogger_InfoWithArguments",
		},
	})
}

func TestLogger_D(t *testing.T) {
	current := time.Now()
	now = func() time.Time {
		return current
	}

	ctx := WithValue(context.Background(), "anyKey", "anyValue")

	mockTestHandler.On("Handle", mock.AnythingOfType("Record")).Return(nil)

	logger := Named("anyLogger")
	logger.D(ctx, "anyMessage")

	mockTestHandler.AssertCalled(t, "Handle", Record{
		Time:       current,
		Level:      LevelDebug,
		Message:    "anyMessage",
		Context:    ctx,
		LoggerName: "anyLogger",
		Caller: Caller{
			defined:  true,
			file:     filename,
			line:     243,
			function: "github.com/procyon-projects/logy.TestLogger_D",
		},
	})
}

func TestLogger_DWithStackTrace(t *testing.T) {
	current := time.Now()
	now = func() time.Time {
		return current
	}
	err := errors.New("anyError")

	ctx := WithValue(context.Background(), "anyKey", "anyValue")

	mockTestHandler.On("Handle", mock.AnythingOfType("Record")).Return(nil)

	logger := Named("anyLogger")
	logger.D(ctx, "anyMessage", err)

	mockTestHandler.AssertCalled(t, "Handle", Record{
		Time:       current,
		Level:      LevelDebug,
		Message:    "anyMessage",
		Context:    ctx,
		LoggerName: "anyLogger",
		Error:      err,
		StackTrace: "github.com/procyon-projects/logy.TestLogger_DWithStackTrace()\\n    /Users/burakkoken/GolandProjects/slog/logger_test.go:272",
		Caller: Caller{
			defined:  true,
			file:     filename,
			line:     272,
			function: "github.com/procyon-projects/logy.TestLogger_DWithStackTrace",
		},
	})
}

func TestLogger_DWithArguments(t *testing.T) {
	current := time.Now()
	now = func() time.Time {
		return current
	}

	ctx := WithValue(context.Background(), "anyKey", "anyValue")

	mockTestHandler.On("Handle", mock.AnythingOfType("Record")).Return(nil)

	logger := Named("anyLogger")
	logger.D(ctx, "anyMessage {}", "anyValue")

	mockTestHandler.AssertCalled(t, "Handle", Record{
		Time:       current,
		Level:      LevelDebug,
		Message:    "anyMessage anyValue",
		Context:    ctx,
		LoggerName: "anyLogger",
		Caller: Caller{
			defined:  true,
			file:     filename,
			line:     302,
			function: "github.com/procyon-projects/logy.TestLogger_DWithArguments",
		},
	})
}

func TestLogger_Debug(t *testing.T) {
	current := time.Now()
	now = func() time.Time {
		return current
	}

	mockTestHandler.On("Handle", mock.AnythingOfType("Record")).Return(nil)

	logger := Named("anyLogger")
	logger.Debug("anyMessage")

	mockTestHandler.AssertCalled(t, "Handle", Record{
		Time:       current,
		Level:      LevelDebug,
		Message:    "anyMessage",
		LoggerName: "anyLogger",
		Caller: Caller{
			defined:  true,
			file:     filename,
			line:     328,
			function: "github.com/procyon-projects/logy.TestLogger_Debug",
		},
	})
}

func TestLogger_DebugWithStackTrace(t *testing.T) {
	current := time.Now()
	now = func() time.Time {
		return current
	}
	err := errors.New("anyError")

	mockTestHandler.On("Handle", mock.AnythingOfType("Record")).Return(nil)

	logger := Named("anyLogger")
	logger.Debug("anyMessage", err)

	mockTestHandler.AssertCalled(t, "Handle", Record{
		Time:       current,
		Level:      LevelDebug,
		Message:    "anyMessage",
		LoggerName: "anyLogger",
		Error:      err,
		StackTrace: "github.com/procyon-projects/logy.TestLogger_DebugWithStackTrace()\\n    /Users/burakkoken/GolandProjects/slog/logger_test.go:354",
		Caller: Caller{
			defined:  true,
			file:     filename,
			line:     354,
			function: "github.com/procyon-projects/logy.TestLogger_DebugWithStackTrace",
		},
	})
}

func TestLogger_DebugWithArguments(t *testing.T) {
	current := time.Now()
	now = func() time.Time {
		return current
	}

	mockTestHandler.On("Handle", mock.AnythingOfType("Record")).Return(nil)

	logger := Named("anyLogger")
	logger.Debug("anyMessage {}", "anyValue")

	mockTestHandler.AssertCalled(t, "Handle", Record{
		Time:       current,
		Level:      LevelDebug,
		Message:    "anyMessage anyValue",
		LoggerName: "anyLogger",
		Caller: Caller{
			defined:  true,
			file:     filename,
			line:     381,
			function: "github.com/procyon-projects/logy.TestLogger_DebugWithArguments",
		},
	})
}

func TestLogger_T(t *testing.T) {
	current := time.Now()
	now = func() time.Time {
		return current
	}

	ctx := WithValue(context.Background(), "anyKey", "anyValue")

	mockTestHandler.On("Handle", mock.AnythingOfType("Record")).Return(nil)

	logger := Named("anyLogger")
	logger.T(ctx, "anyMessage")

	mockTestHandler.AssertCalled(t, "Handle", Record{
		Time:       current,
		Level:      LevelTrace,
		Message:    "anyMessage",
		Context:    ctx,
		LoggerName: "anyLogger",
		Caller: Caller{
			defined:  true,
			file:     filename,
			line:     408,
			function: "github.com/procyon-projects/logy.TestLogger_T",
		},
	})
}

func TestLogger_TWithStackTrace(t *testing.T) {
	current := time.Now()
	now = func() time.Time {
		return current
	}
	err := errors.New("anyError")

	ctx := WithValue(context.Background(), "anyKey", "anyValue")

	mockTestHandler.On("Handle", mock.AnythingOfType("Record")).Return(nil)

	logger := Named("anyLogger")
	logger.T(ctx, "anyMessage", err)

	mockTestHandler.AssertCalled(t, "Handle", Record{
		Time:       current,
		Level:      LevelTrace,
		Message:    "anyMessage",
		Context:    ctx,
		LoggerName: "anyLogger",
		Error:      err,
		StackTrace: "github.com/procyon-projects/logy.TestLogger_TWithStackTrace()\\n    /Users/burakkoken/GolandProjects/slog/logger_test.go:437",
		Caller: Caller{
			defined:  true,
			file:     filename,
			line:     437,
			function: "github.com/procyon-projects/logy.TestLogger_TWithStackTrace",
		},
	})
}

func TestLogger_TWithArguments(t *testing.T) {
	current := time.Now()
	now = func() time.Time {
		return current
	}

	ctx := WithValue(context.Background(), "anyKey", "anyValue")

	mockTestHandler.On("Handle", mock.AnythingOfType("Record")).Return(nil)

	logger := Named("anyLogger")
	logger.T(ctx, "anyMessage {}", "anyValue")

	mockTestHandler.AssertCalled(t, "Handle", Record{
		Time:       current,
		Level:      LevelTrace,
		Message:    "anyMessage anyValue",
		Context:    ctx,
		LoggerName: "anyLogger",
		Caller: Caller{
			defined:  true,
			file:     filename,
			line:     467,
			function: "github.com/procyon-projects/logy.TestLogger_TWithArguments",
		},
	})
}

func TestLogger_Trace(t *testing.T) {
	current := time.Now()
	now = func() time.Time {
		return current
	}

	mockTestHandler.On("Handle", mock.AnythingOfType("Record")).Return(nil)

	logger := Named("anyLogger")
	logger.Trace("anyMessage")

	mockTestHandler.AssertCalled(t, "Handle", Record{
		Time:       current,
		Level:      LevelTrace,
		Message:    "anyMessage",
		LoggerName: "anyLogger",
		Caller: Caller{
			defined:  true,
			file:     filename,
			line:     493,
			function: "github.com/procyon-projects/logy.TestLogger_Trace",
		},
	})
}

func TestLogger_TraceWithStackTrace(t *testing.T) {
	current := time.Now()
	now = func() time.Time {
		return current
	}
	err := errors.New("anyError")

	mockTestHandler.On("Handle", mock.AnythingOfType("Record")).Return(nil)

	logger := Named("anyLogger")
	logger.Trace("anyMessage", err)

	mockTestHandler.AssertCalled(t, "Handle", Record{
		Time:       current,
		Level:      LevelTrace,
		Message:    "anyMessage",
		LoggerName: "anyLogger",
		Error:      err,
		StackTrace: "github.com/procyon-projects/logy.TestLogger_TraceWithStackTrace()\\n    /Users/burakkoken/GolandProjects/slog/logger_test.go:519",
		Caller: Caller{
			defined:  true,
			file:     filename,
			line:     519,
			function: "github.com/procyon-projects/logy.TestLogger_TraceWithStackTrace",
		},
	})
}

func TestLogger_TraceWithArguments(t *testing.T) {
	current := time.Now()
	now = func() time.Time {
		return current
	}

	mockTestHandler.On("Handle", mock.AnythingOfType("Record")).Return(nil)

	logger := Named("anyLogger")
	logger.Trace("anyMessage {}", "anyValue")

	mockTestHandler.AssertCalled(t, "Handle", Record{
		Time:       current,
		Level:      LevelTrace,
		Message:    "anyMessage anyValue",
		LoggerName: "anyLogger",
		Caller: Caller{
			defined:  true,
			file:     filename,
			line:     546,
			function: "github.com/procyon-projects/logy.TestLogger_TraceWithArguments",
		},
	})
}

func TestLogger_E(t *testing.T) {
	current := time.Now()
	now = func() time.Time {
		return current
	}

	ctx := WithValue(context.Background(), "anyKey", "anyValue")

	mockTestHandler.On("Handle", mock.AnythingOfType("Record")).Return(nil)

	logger := Named("anyLogger")
	logger.E(ctx, "anyMessage")

	mockTestHandler.AssertCalled(t, "Handle", Record{
		Time:       current,
		Level:      LevelError,
		Message:    "anyMessage",
		Context:    ctx,
		LoggerName: "anyLogger",
		Caller: Caller{
			defined:  true,
			file:     filename,
			line:     573,
			function: "github.com/procyon-projects/logy.TestLogger_E",
		},
	})
}

func TestLogger_EWithStackTrace(t *testing.T) {
	current := time.Now()
	now = func() time.Time {
		return current
	}
	err := errors.New("anyError")

	ctx := WithValue(context.Background(), "anyKey", "anyValue")

	mockTestHandler.On("Handle", mock.AnythingOfType("Record")).Return(nil)

	logger := Named("anyLogger")
	logger.E(ctx, "anyMessage", err)

	mockTestHandler.AssertCalled(t, "Handle", Record{
		Time:       current,
		Level:      LevelError,
		Message:    "anyMessage",
		Context:    ctx,
		LoggerName: "anyLogger",
		Error:      err,
		StackTrace: "github.com/procyon-projects/logy.TestLogger_EWithStackTrace()\\n    /Users/burakkoken/GolandProjects/slog/logger_test.go:602",
		Caller: Caller{
			defined:  true,
			file:     filename,
			line:     602,
			function: "github.com/procyon-projects/logy.TestLogger_EWithStackTrace",
		},
	})
}

func TestLogger_EWithArguments(t *testing.T) {
	current := time.Now()
	now = func() time.Time {
		return current
	}

	ctx := WithValue(context.Background(), "anyKey", "anyValue")

	mockTestHandler.On("Handle", mock.AnythingOfType("Record")).Return(nil)

	logger := Named("anyLogger")
	logger.E(ctx, "anyMessage {}", "anyValue")

	mockTestHandler.AssertCalled(t, "Handle", Record{
		Time:       current,
		Level:      LevelError,
		Message:    "anyMessage anyValue",
		Context:    ctx,
		LoggerName: "anyLogger",
		Caller: Caller{
			defined:  true,
			file:     filename,
			line:     632,
			function: "github.com/procyon-projects/logy.TestLogger_EWithArguments",
		},
	})
}

func TestLogger_Error(t *testing.T) {
	current := time.Now()
	now = func() time.Time {
		return current
	}

	mockTestHandler.On("Handle", mock.AnythingOfType("Record")).Return(nil)

	logger := Named("anyLogger")
	logger.Error("anyMessage")

	mockTestHandler.AssertCalled(t, "Handle", Record{
		Time:       current,
		Level:      LevelError,
		Message:    "anyMessage",
		LoggerName: "anyLogger",
		Caller: Caller{
			defined:  true,
			file:     filename,
			line:     658,
			function: "github.com/procyon-projects/logy.TestLogger_Error",
		},
	})
}

func TestLogger_ErrorWithStackTrace(t *testing.T) {
	current := time.Now()
	now = func() time.Time {
		return current
	}
	err := errors.New("anyError")

	mockTestHandler.On("Handle", mock.AnythingOfType("Record")).Return(nil)

	logger := Named("anyLogger")
	logger.Error("anyMessage", err)

	mockTestHandler.AssertCalled(t, "Handle", Record{
		Time:       current,
		Level:      LevelError,
		Message:    "anyMessage",
		LoggerName: "anyLogger",
		Error:      err,
		StackTrace: "github.com/procyon-projects/logy.TestLogger_ErrorWithStackTrace()\\n    /Users/burakkoken/GolandProjects/slog/logger_test.go:684",
		Caller: Caller{
			defined:  true,
			file:     filename,
			line:     684,
			function: "github.com/procyon-projects/logy.TestLogger_ErrorWithStackTrace",
		},
	})
}

func TestLogger_ErrorWithArguments(t *testing.T) {
	current := time.Now()
	now = func() time.Time {
		return current
	}

	mockTestHandler.On("Handle", mock.AnythingOfType("Record")).Return(nil)

	logger := Named("anyLogger")
	logger.Error("anyMessage {}", "anyValue")

	mockTestHandler.AssertCalled(t, "Handle", Record{
		Time:       current,
		Level:      LevelError,
		Message:    "anyMessage anyValue",
		LoggerName: "anyLogger",
		Caller: Caller{
			defined:  true,
			file:     filename,
			line:     711,
			function: "github.com/procyon-projects/logy.TestLogger_ErrorWithArguments",
		},
	})
}

func TestLogger_W(t *testing.T) {
	current := time.Now()
	now = func() time.Time {
		return current
	}

	ctx := WithValue(context.Background(), "anyKey", "anyValue")

	mockTestHandler.On("Handle", mock.AnythingOfType("Record")).Return(nil)

	logger := Named("anyLogger")
	logger.W(ctx, "anyMessage")

	mockTestHandler.AssertCalled(t, "Handle", Record{
		Time:       current,
		Level:      LevelWarn,
		Message:    "anyMessage",
		Context:    ctx,
		LoggerName: "anyLogger",
		Caller: Caller{
			defined:  true,
			file:     filename,
			line:     738,
			function: "github.com/procyon-projects/logy.TestLogger_W",
		},
	})
}

func TestLogger_WWithStackTrace(t *testing.T) {
	current := time.Now()
	now = func() time.Time {
		return current
	}
	err := errors.New("anyError")

	ctx := WithValue(context.Background(), "anyKey", "anyValue")

	mockTestHandler.On("Handle", mock.AnythingOfType("Record")).Return(nil)

	logger := Named("anyLogger")
	logger.W(ctx, "anyMessage", err)

	mockTestHandler.AssertCalled(t, "Handle", Record{
		Time:       current,
		Level:      LevelWarn,
		Message:    "anyMessage",
		Context:    ctx,
		LoggerName: "anyLogger",
		Error:      err,
		StackTrace: "github.com/procyon-projects/logy.TestLogger_WWithStackTrace()\\n    /Users/burakkoken/GolandProjects/slog/logger_test.go:767",
		Caller: Caller{
			defined:  true,
			file:     filename,
			line:     767,
			function: "github.com/procyon-projects/logy.TestLogger_WWithStackTrace",
		},
	})
}

func TestLogger_WWithArguments(t *testing.T) {
	current := time.Now()
	now = func() time.Time {
		return current
	}

	ctx := WithValue(context.Background(), "anyKey", "anyValue")

	mockTestHandler.On("Handle", mock.AnythingOfType("Record")).Return(nil)

	logger := Named("anyLogger")
	logger.W(ctx, "anyMessage {}", "anyValue")

	mockTestHandler.AssertCalled(t, "Handle", Record{
		Time:       current,
		Level:      LevelWarn,
		Message:    "anyMessage anyValue",
		Context:    ctx,
		LoggerName: "anyLogger",
		Caller: Caller{
			defined:  true,
			file:     filename,
			line:     797,
			function: "github.com/procyon-projects/logy.TestLogger_WWithArguments",
		},
	})
}

func TestLogger_Warn(t *testing.T) {
	current := time.Now()
	now = func() time.Time {
		return current
	}

	mockTestHandler.On("Handle", mock.AnythingOfType("Record")).Return(nil)

	logger := Named("anyLogger")
	logger.Warn("anyMessage")

	mockTestHandler.AssertCalled(t, "Handle", Record{
		Time:       current,
		Level:      LevelWarn,
		Message:    "anyMessage",
		LoggerName: "anyLogger",
		Caller: Caller{
			defined:  true,
			file:     filename,
			line:     823,
			function: "github.com/procyon-projects/logy.TestLogger_Warn",
		},
	})
}

func TestLogger_WarnWithArguments(t *testing.T) {
	current := time.Now()
	now = func() time.Time {
		return current
	}

	mockTestHandler.On("Handle", mock.AnythingOfType("Record")).Return(nil)

	logger := Named("anyLogger")
	logger.Warn("anyMessage {}", "anyValue")

	mockTestHandler.AssertCalled(t, "Handle", Record{
		Time:       current,
		Level:      LevelWarn,
		Message:    "anyMessage anyValue",
		LoggerName: "anyLogger",
		Caller: Caller{
			defined:  true,
			file:     filename,
			line:     848,
			function: "github.com/procyon-projects/logy.TestLogger_WarnWithArguments",
		},
	})
}

func TestLogger_WarnWithStackTrace(t *testing.T) {
	current := time.Now()
	now = func() time.Time {
		return current
	}
	err := errors.New("anyError")

	mockTestHandler.On("Handle", mock.AnythingOfType("Record")).Return(nil)

	logger := Named("anyLogger")
	logger.Warn("anyMessage", err)

	mockTestHandler.AssertCalled(t, "Handle", Record{
		Time:       current,
		Level:      LevelWarn,
		Message:    "anyMessage",
		LoggerName: "anyLogger",
		Error:      err,
		StackTrace: "github.com/procyon-projects/logy.TestLogger_WarnWithStackTrace()\\n    /Users/burakkoken/GolandProjects/slog/logger_test.go:874",
		Caller: Caller{
			defined:  true,
			file:     filename,
			line:     874,
			function: "github.com/procyon-projects/logy.TestLogger_WarnWithStackTrace",
		},
	})
}
