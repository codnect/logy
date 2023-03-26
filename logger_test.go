package logy

import (
	"context"
	"errors"
	"fmt"
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
			line:     79,
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
		StackTrace: fmt.Sprintf("github.com/procyon-projects/logy.TestLogger_IWithStackTrace()\\n    %s:108", filename),
		Caller: Caller{
			defined:  true,
			file:     filename,
			line:     108,
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
			line:     138,
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
			line:     164,
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
		StackTrace: fmt.Sprintf("github.com/procyon-projects/logy.TestLogger_InfoWithStackTrace()\\n    %s:190", filename),
		Caller: Caller{
			defined:  true,
			file:     filename,
			line:     190,
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
			line:     217,
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
			line:     244,
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
		StackTrace: fmt.Sprintf("github.com/procyon-projects/logy.TestLogger_DWithStackTrace()\\n    %s:273", filename),
		Caller: Caller{
			defined:  true,
			file:     filename,
			line:     273,
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
			line:     303,
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
			line:     329,
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
		StackTrace: fmt.Sprintf("github.com/procyon-projects/logy.TestLogger_DebugWithStackTrace()\\n    %s:355", filename),
		Caller: Caller{
			defined:  true,
			file:     filename,
			line:     355,
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
			line:     382,
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
			line:     409,
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
		StackTrace: fmt.Sprintf("github.com/procyon-projects/logy.TestLogger_TWithStackTrace()\\n    %s:438", filename),
		Caller: Caller{
			defined:  true,
			file:     filename,
			line:     438,
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
			line:     468,
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
			line:     494,
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
		StackTrace: fmt.Sprintf("github.com/procyon-projects/logy.TestLogger_TraceWithStackTrace()\\n    %s:520", filename),
		Caller: Caller{
			defined:  true,
			file:     filename,
			line:     520,
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
			line:     547,
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
			line:     574,
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
		StackTrace: fmt.Sprintf("github.com/procyon-projects/logy.TestLogger_EWithStackTrace()\\n    %s:603", filename),
		Caller: Caller{
			defined:  true,
			file:     filename,
			line:     603,
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
			line:     633,
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
			line:     659,
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
		StackTrace: fmt.Sprintf("github.com/procyon-projects/logy.TestLogger_ErrorWithStackTrace()\\n    %s:685", filename),
		Caller: Caller{
			defined:  true,
			file:     filename,
			line:     685,
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
			line:     712,
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
			line:     739,
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
		StackTrace: fmt.Sprintf("github.com/procyon-projects/logy.TestLogger_WWithStackTrace()\\n    %s:768", filename),
		Caller: Caller{
			defined:  true,
			file:     filename,
			line:     768,
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
			line:     798,
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
			line:     824,
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
			line:     849,
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
		StackTrace: "github.com/procyon-projects/logy.TestLogger_WarnWithStackTrace()\\n    /Users/burakkoken/GolandProjects/slog/logger_test.go:875",
		Caller: Caller{
			defined:  true,
			file:     filename,
			line:     875,
			function: "github.com/procyon-projects/logy.TestLogger_WarnWithStackTrace",
		},
	})
}
