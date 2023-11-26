package logy

import (
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
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
			line:     80,
			function: "codnect.io/logy.TestLogger_I",
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
		StackTrace: fmt.Sprintf("codnect.io/logy.TestLogger_IWithStackTrace()\\n    %s:109", filename),
		Caller: Caller{
			defined:  true,
			file:     filename,
			line:     109,
			function: "codnect.io/logy.TestLogger_IWithStackTrace",
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
			line:     139,
			function: "codnect.io/logy.TestLogger_IWithArguments",
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
			line:     165,
			function: "codnect.io/logy.TestLogger_Info",
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
		StackTrace: fmt.Sprintf("codnect.io/logy.TestLogger_InfoWithStackTrace()\\n    %s:191", filename),
		Caller: Caller{
			defined:  true,
			file:     filename,
			line:     191,
			function: "codnect.io/logy.TestLogger_InfoWithStackTrace",
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
			line:     218,
			function: "codnect.io/logy.TestLogger_InfoWithArguments",
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
			line:     245,
			function: "codnect.io/logy.TestLogger_D",
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
		StackTrace: fmt.Sprintf("codnect.io/logy.TestLogger_DWithStackTrace()\\n    %s:274", filename),
		Caller: Caller{
			defined:  true,
			file:     filename,
			line:     274,
			function: "codnect.io/logy.TestLogger_DWithStackTrace",
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
			line:     304,
			function: "codnect.io/logy.TestLogger_DWithArguments",
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
			line:     330,
			function: "codnect.io/logy.TestLogger_Debug",
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
		StackTrace: fmt.Sprintf("codnect.io/logy.TestLogger_DebugWithStackTrace()\\n    %s:356", filename),
		Caller: Caller{
			defined:  true,
			file:     filename,
			line:     356,
			function: "codnect.io/logy.TestLogger_DebugWithStackTrace",
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
			line:     383,
			function: "codnect.io/logy.TestLogger_DebugWithArguments",
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
			line:     410,
			function: "codnect.io/logy.TestLogger_T",
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
		StackTrace: fmt.Sprintf("codnect.io/logy.TestLogger_TWithStackTrace()\\n    %s:439", filename),
		Caller: Caller{
			defined:  true,
			file:     filename,
			line:     439,
			function: "codnect.io/logy.TestLogger_TWithStackTrace",
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
			line:     469,
			function: "codnect.io/logy.TestLogger_TWithArguments",
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
			line:     495,
			function: "codnect.io/logy.TestLogger_Trace",
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
		StackTrace: fmt.Sprintf("codnect.io/logy.TestLogger_TraceWithStackTrace()\\n    %s:521", filename),
		Caller: Caller{
			defined:  true,
			file:     filename,
			line:     521,
			function: "codnect.io/logy.TestLogger_TraceWithStackTrace",
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
			line:     548,
			function: "codnect.io/logy.TestLogger_TraceWithArguments",
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
			line:     575,
			function: "codnect.io/logy.TestLogger_E",
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
		StackTrace: fmt.Sprintf("codnect.io/logy.TestLogger_EWithStackTrace()\\n    %s:604", filename),
		Caller: Caller{
			defined:  true,
			file:     filename,
			line:     604,
			function: "codnect.io/logy.TestLogger_EWithStackTrace",
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
			line:     634,
			function: "codnect.io/logy.TestLogger_EWithArguments",
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
			line:     660,
			function: "codnect.io/logy.TestLogger_Error",
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
		StackTrace: fmt.Sprintf("codnect.io/logy.TestLogger_ErrorWithStackTrace()\\n    %s:686", filename),
		Caller: Caller{
			defined:  true,
			file:     filename,
			line:     686,
			function: "codnect.io/logy.TestLogger_ErrorWithStackTrace",
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
			line:     713,
			function: "codnect.io/logy.TestLogger_ErrorWithArguments",
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
			line:     740,
			function: "codnect.io/logy.TestLogger_W",
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
		StackTrace: fmt.Sprintf("codnect.io/logy.TestLogger_WWithStackTrace()\\n    %s:769", filename),
		Caller: Caller{
			defined:  true,
			file:     filename,
			line:     769,
			function: "codnect.io/logy.TestLogger_WWithStackTrace",
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
			line:     799,
			function: "codnect.io/logy.TestLogger_WWithArguments",
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
			line:     825,
			function: "codnect.io/logy.TestLogger_Warn",
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
			line:     850,
			function: "codnect.io/logy.TestLogger_WarnWithArguments",
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
		StackTrace: fmt.Sprintf("codnect.io/logy.TestLogger_WarnWithStackTrace()\\n    %s:876", filename),
		Caller: Caller{
			defined:  true,
			file:     filename,
			line:     876,
			function: "codnect.io/logy.TestLogger_WarnWithStackTrace",
		},
	})
}

func TestLogger_IsDebugEnabled(t *testing.T) {
	logger := &Logger{}
	logger.SetLevel(LevelWarn)
	assert.False(t, logger.IsDebugEnabled())

	logger.SetLevel(LevelAll)
	assert.True(t, logger.IsDebugEnabled())
}

func TestLogger_IsInfoEnabled(t *testing.T) {
	logger := &Logger{}
	logger.SetLevel(LevelError)
	assert.False(t, logger.IsInfoEnabled())

	logger.SetLevel(LevelAll)
	assert.True(t, logger.IsInfoEnabled())
}

func TestLogger_IsErrorEnabled(t *testing.T) {
	logger := &Logger{}
	logger.SetLevel(LevelOff)
	assert.False(t, logger.IsErrorEnabled())

	logger.SetLevel(LevelAll)
	assert.True(t, logger.IsErrorEnabled())
}

func TestLogger_IsWarnEnabled(t *testing.T) {
	logger := &Logger{}
	logger.SetLevel(LevelError)
	assert.False(t, logger.IsWarnEnabled())

	logger.SetLevel(LevelAll)
	assert.True(t, logger.IsWarnEnabled())
}

func TestLogger_IsTraceEnabled(t *testing.T) {
	logger := &Logger{}
	logger.SetLevel(LevelError)
	assert.False(t, logger.IsTraceEnabled())

	logger.SetLevel(LevelAll)
	assert.True(t, logger.IsTraceEnabled())
}

func TestOf(t *testing.T) {
	log := Of[testHandler]()
	assert.Equal(t, "codnect.io/logy.testHandler", log.Name())
}

func TestGet(t *testing.T) {
	log := Get()
	assert.Equal(t, "codnect.io/logy", log.Name())
}

func TestNamed(t *testing.T) {
	log := Named("anyLogger")
	assert.Equal(t, "anyLogger", log.Name())
}
