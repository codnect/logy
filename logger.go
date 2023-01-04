package slog

import (
	"context"
	"log"
	"sync"
	"time"
)

var (
	rootLogger    = New("")
	defaultLogger = rootLogger

	loggers = map[string]*Logger{}
	mu      sync.RWMutex
)

type Logger struct {
	level   Level
	handler Handler

	recordPool sync.Pool
	mu         sync.RWMutex
}

func SetDefault(logger *Logger) {
	if logger == nil {
		panic("logger cannot be nil")
	}

	defer mu.Unlock()
	mu.Lock()

	log.SetOutput(newWriter(logger))
	log.SetFlags(0)
}

func Default() *Logger {
	defer mu.Unlock()
	mu.Lock()

	return defaultLogger
}

func Of[T any]() *Logger {
	return nil
}

func New(name string) *Logger {
	defer mu.Unlock()
	mu.Lock()

	if logger, ok := loggers[name]; ok {
		return logger
	}

	logger := &Logger{
		recordPool: sync.Pool{
			New: func() any {
				return &Record{}
			},
		},
		level:   LevelInfo,
		handler: NewConsoleHandler(),
	}
	loggers[name] = logger
	return logger
}

func (l *Logger) SetLevel(level Level) {
	defer l.mu.Unlock()
	l.mu.Lock()
	l.level = level
}

func (l *Logger) IsLoggable(level Level) bool {
	defer l.mu.Unlock()
	l.mu.Lock()
	return level >= l.level
}

func (l *Logger) I(ctx context.Context, msg string, args ...string) {
	l.logDepth(0, ctx, LevelInfo, msg, args...)
}

func (l *Logger) E(ctx context.Context, err error, args ...string) {
	l.logDepth(0, ctx, LevelError, err.Error(), args...)
}

func (l *Logger) W(ctx context.Context, msg string, args ...string) {
	l.logDepth(0, ctx, LevelWarn, msg, args...)
}

func (l *Logger) D(ctx context.Context, msg string, args ...string) {
	l.logDepth(0, ctx, LevelDebug, msg, args...)
}

func (l *Logger) A(ctx context.Context, level Level, msg string, attrs ...*Attr) {

}

func (l *Logger) L(ctx context.Context, level Level, msg string, args ...string) {
	l.logDepth(0, ctx, level, msg, args...)
}

func (l *Logger) Info(msg string, args ...string) {
	l.logDepth(0, nil, LevelInfo, msg, args...)
}

func (l *Logger) Error(err error, args ...string) {
	l.logDepth(0, nil, LevelError, err.Error(), args...)
}

func (l *Logger) Warn(msg string, args ...string) {
	l.logDepth(0, nil, LevelWarn, msg, args...)
}

func (l *Logger) Debug(msg string, args ...string) {
	l.logDepth(0, nil, LevelDebug, msg, args...)
}

func (l *Logger) Attrs(level Level, msg string, attrs ...*Attr) {
	l.A(nil, level, msg, attrs...)
}

func (l *Logger) Log(level Level, msg string, args ...string) {
	l.logDepth(0, nil, level, msg, args...)
}

func (l *Logger) logDepth(depth int, ctx context.Context, level Level, msg string, args ...string) {
	if !l.IsLoggable(level) {
		return
	}

	record := l.recordPool.Get().(*Record)
	record.Time = time.Now()
	record.Level = level
	record.Message = msg
	record.Context = ctx
	record.depth = 0
	_ = l.logRecord(record)
}

func (l *Logger) logRecord(record *Record) error {
	l.mu.Lock()
	handler := l.handler
	l.mu.Unlock()

	return handler.Handle(record)
}

func I(ctx context.Context, msg string, args ...string) {
	Default().logDepth(0, ctx, LevelInfo, msg, args...)
}

func E(ctx context.Context, err error, args ...string) {
	Default().logDepth(0, ctx, LevelError, err.Error(), args...)
}

func W(ctx context.Context, msg string, args ...string) {
	Default().logDepth(0, ctx, LevelWarn, msg, args...)
}

func D(ctx context.Context, msg string, args ...string) {
	Default().logDepth(0, ctx, LevelWarn, msg, args...)
}

func A(ctx context.Context, level Level, msg string, attrs ...*Attr) {

}

func L(ctx context.Context, level Level, msg string, args ...string) {
	Default().logDepth(0, ctx, level, msg, args...)
}

func Info(msg string, args ...string) {
	Default().logDepth(0, nil, LevelInfo, msg, args...)
}

func Error(err error, args ...string) {
	Default().logDepth(0, nil, LevelError, err.Error(), args...)
}

func Warn(msg string, args ...string) {
	Default().logDepth(0, nil, LevelWarn, msg, args...)
}

func Debug(msg string, args ...string) {
	Default().logDepth(0, nil, LevelWarn, msg, args...)
}

func Attrs(level Level, msg string, attrs ...*Attr) {

}

func Log(level Level, msg string, args ...string) {
	Default().logDepth(0, nil, level, msg, args...)
}
