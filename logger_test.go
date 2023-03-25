package logy

import (
	"context"
	"github.com/stretchr/testify/mock"
	"io"
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
	mockTestHandler = &testHandler{}
)

func init() {
	RegisterHandler("testHandler", mockTestHandler)
	LoadConfig(&Config{
		Level:    LevelAll,
		Handlers: Handlers{"testHandler"},
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
	})
}
