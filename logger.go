package logy

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var (
	rootLogger    = Named("")
	defaultLogger atomic.Value

	loggers = map[string]*Logger{}
	mu      sync.RWMutex
)

func init() {
	defaultLogger.Store(rootLogger)
}

type Logger struct {
	name    string
	level   Level
	handler Handler

	mu sync.RWMutex
}

func SetDefault(logger *Logger) {
	if logger == nil {
		panic("logger cannot be nil")
	}

	defaultLogger.Store(logger)

	log.SetOutput(newWriter(logger))
	log.SetFlags(0)
}

func Default() *Logger {
	return defaultLogger.Load().(*Logger)
}

func Of[T any]() *Logger {
	typ := reflect.TypeOf((*T)(nil)).Elem()
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	name := typ.Name()
	pkg := strings.ReplaceAll(typ.PkgPath(), "/", ".")
	return Named(fmt.Sprintf("%s.%s", pkg, name))
}

func Named(name string) *Logger {
	defer mu.Unlock()
	mu.Lock()

	if logger, ok := loggers[name]; ok {
		return logger
	}

	logger := &Logger{
		name:    name,
		level:   LevelError,
		handler: NewConsoleHandler(),
	}

	loggers[name] = logger
	return logger
}

func New() *Logger {
	return Named("")
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

func (l *Logger) I(ctx context.Context, msg string, args ...any) {
	l.logDepth(1, ctx, LevelInfo, msg, args...)
}

func (l *Logger) E(ctx context.Context, err error, args ...any) {
	l.logDepth(1, ctx, LevelError, err.Error(), args...)
}

func (l *Logger) W(ctx context.Context, msg string, args ...any) {
	l.logDepth(1, ctx, LevelWarn, msg, args...)
}

func (l *Logger) D(ctx context.Context, msg string, args ...any) {
	l.logDepth(1, ctx, LevelDebug, msg, args...)
}

func (l *Logger) A(ctx context.Context, level Level, msg string, attrs ...Attribute) {
	l.logAttrsDepth(1, ctx, level, msg, attrs...)
}

func (l *Logger) L(ctx context.Context, level Level, msg string, args ...any) {
	l.logDepth(1, ctx, level, msg, args...)
}

func (l *Logger) Info(msg string, args ...any) {
	l.logDepth(1, nil, LevelInfo, msg, args...)
}

func (l *Logger) Error(err error, args ...any) {
	l.logDepth(1, nil, LevelError, err.Error(), args...)
}

func (l *Logger) Warn(msg string, args ...any) {
	l.logDepth(1, nil, LevelWarn, msg, args...)
}

func (l *Logger) Debug(msg string, args ...any) {
	l.logDepth(1, nil, LevelDebug, msg, args...)
}

func (l *Logger) Attrs(level Level, msg string, attrs ...Attribute) {
	l.logAttrsDepth(1, nil, level, msg, attrs...)
}

func (l *Logger) Log(level Level, msg string, args ...any) {
	l.logDepth(1, nil, level, msg, args...)
}

func (l *Logger) logDepth(depth int, ctx context.Context, level Level, msg string, args ...any) error {
	/*if !l.handler.IsLoggable(level) {
		return
	}*/

	record := l.makeRecord(depth, ctx, msg, level)
	return l.handler.Handle(record)
}

func (l *Logger) logAttrsDepth(depth int, ctx context.Context, level Level, msg string, attrs ...Attribute) {
	/*if !l.handler.IsLoggable(level) {
		return
	}*/

	record := l.makeRecord(depth, ctx, msg, level)
	_ = l.handler.Handle(record)
}

func (l *Logger) makeRecord(depth int, ctx context.Context, msg string, level Level) Record {
	var pcs [1]uintptr
	runtime.Callers(depth+3, pcs[:])

	frames := runtime.CallersFrames(pcs[:1])
	frame, _ := frames.Next()

	return Record{
		Time:       time.Now(),
		Message:    msg,
		Level:      level,
		Context:    ctx,
		LoggerName: l.name,
		Caller: Caller{
			Defined:  frame.PC != 0,
			PC:       frame.PC,
			File:     frame.File,
			Line:     frame.Line,
			Function: frame.Function,
		},
	}
}

func I(ctx context.Context, msg string, args ...any) {
	Default().logDepth(1, ctx, LevelInfo, msg, args...)
}

func E(ctx context.Context, err error, args ...any) {
	Default().logDepth(1, ctx, LevelError, err.Error(), args...)
}

func W(ctx context.Context, msg string, args ...any) {
	Default().logDepth(1, ctx, LevelWarn, msg, args...)
}

func D(ctx context.Context, msg string, args ...any) {
	Default().logDepth(1, ctx, LevelWarn, msg, args...)
}

func A(ctx context.Context, level Level, msg string, attrs ...Attribute) {
	Default().logAttrsDepth(1, ctx, level, msg, attrs...)
}

func L(ctx context.Context, level Level, msg string, args ...any) {
	Default().logDepth(0, ctx, level, msg, args...)
}

func Info(msg string, args ...any) {
	Default().logDepth(1, nil, LevelInfo, msg, args...)
}

func Error(err error, args ...any) {
	Default().logDepth(1, nil, LevelError, err.Error(), args...)
}

func Warn(msg string, args ...any) {
	Default().logDepth(1, nil, LevelWarn, msg, args...)
}

func Debug(msg string, args ...any) {
	Default().logDepth(1, nil, LevelWarn, msg, args...)
}

func Attrs(level Level, msg string, attrs ...Attribute) {
	Default().logAttrsDepth(1, nil, level, msg, attrs...)
}

func Log(level Level, msg string, args ...any) {
	Default().logDepth(1, nil, level, msg, args...)
}
