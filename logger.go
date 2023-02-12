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

const (
	RootLoggerName = "github.com.procyon-projects.procyon.app"
)

var (
	rootLogger    = newLogger(RootLoggerName, LevelTrace, nil)
	defaultLogger atomic.Value

	cache         = map[string]*Logger{}
	loggerCacheMu sync.RWMutex
)

func init() {
	defer loggerCacheMu.Unlock()
	loggerCacheMu.Lock()
	cache[RootLoggerName] = rootLogger
	defaultLogger.Store(rootLogger)
}

type Logger struct {
	name     string
	level    atomic.Value
	handlers map[string]Handler

	parent   *Logger
	children map[string]*Logger

	mu sync.RWMutex
}

func Default() *Logger {
	return defaultLogger.Load().(*Logger)
}

func SetDefault(logger *Logger) {
	if logger == nil {
		panic("logger cannot be nil")
	}

	defaultLogger.Store(logger)

	log.SetOutput(newWriter(logger))
	log.SetFlags(0)
}

func New() *Logger {
	rpc := make([]uintptr, 1)
	runtime.Callers(1, rpc[:])
	frame, _ := runtime.CallersFrames(rpc).Next()

	lastSlash := strings.LastIndexByte(frame.Function, '/')
	if lastSlash < 0 {
		lastSlash = 0
	}

	lastDot := strings.IndexByte(frame.Function[lastSlash:], '.') + lastSlash
	return getLogger(frame.Function[:lastDot], "")
}

func Of[T any]() *Logger {
	typ := reflect.TypeOf((*T)(nil)).Elem()

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	name := typ.Name()
	return getLogger(typ.PkgPath(), name)
}

func Named(name string) *Logger {
	return getLogger(name, "")
}

func getLogger(pkg string, typeName string) *Logger {
	defer loggerCacheMu.Unlock()
	loggerCacheMu.Lock()

	loggerName := pkg
	if typeName != "" {
		loggerName = fmt.Sprintf("%s.%s", pkg, typeName)
	}

	if logger, ok := cache[loggerName]; ok {
		return logger
	}

	names := strings.Split(pkg, "/")
	if typeName != "" {
		names = append(names, typeName)
	}

	logger := rootLogger
	loggerName = ""

	for index, value := range names {
		if loggerName == "" {
			loggerName = value
		} else if len(names)-1 == index && typeName != "" {
			loggerName = fmt.Sprintf("%s.%s", pkg, typeName)
		} else {
			loggerName = fmt.Sprintf("%s/%s", loggerName, value)
		}

		childLogger, ok := logger.getChildLogger(loggerName)
		if !ok {
			logger = logger.createChildLogger(loggerName)
		} else {
			logger = childLogger
		}

		if _, exists := cache[loggerName]; !exists {
			cache[loggerName] = logger
		}
	}

	return logger
}

func newLogger(name string, level Level, parent *Logger) *Logger {
	logger := &Logger{
		name:     name,
		handlers: map[string]Handler{},
		parent:   parent,
		children: map[string]*Logger{},
	}

	logger.level.Store(level)
	return logger
}

func (l *Logger) createChildLogger(name string) *Logger {
	defer l.mu.Unlock()
	l.mu.Lock()

	logger := newLogger(name, LevelInfo, l)
	l.children[name] = logger
	return logger
}

func (l *Logger) getChildLogger(name string) (*Logger, bool) {
	defer l.mu.Unlock()
	l.mu.Lock()

	if logger, ok := l.children[name]; ok {
		return logger, true
	}

	return nil, false
}

func (l *Logger) Name() string {
	return l.name
}

func (l *Logger) SetLevel(level Level) {
	l.level.Store(level)
}

func (l *Logger) Level() Level {
	return l.level.Load().(Level)
}

func (l *Logger) IsLoggable(level Level) bool {
	return level <= l.Level()
}

func (l *Logger) I(ctx context.Context, msg string, args ...any) {
	l.logDepth(1, ctx, LevelInfo, msg, args...)
}

func (l *Logger) Info(msg string, args ...any) {
	l.logDepth(1, nil, LevelInfo, msg, args...)
}

func (l *Logger) IsInfoEnabled() bool {
	return LevelInfo <= l.Level()
}

func (l *Logger) E(ctx context.Context, msg string, args ...any) {
	l.logDepth(1, ctx, LevelError, msg, args...)
}

func (l *Logger) Error(msg string, args ...any) {
	l.logDepth(1, nil, LevelError, msg, args...)
}

func (l *Logger) IsErrorEnabled() bool {
	return LevelError <= l.Level()
}

func (l *Logger) W(ctx context.Context, msg string, args ...any) {
	l.logDepth(1, ctx, LevelWarn, msg, args...)
}

func (l *Logger) Warn(msg string, args ...any) {
	l.logDepth(1, nil, LevelWarn, msg, args...)
}

func (l *Logger) IsWarnEnabled() bool {
	return LevelWarn <= l.Level()
}

func (l *Logger) D(ctx context.Context, msg string, args ...any) {
	l.logDepth(1, ctx, LevelDebug, msg, args...)
}

func (l *Logger) Debug(msg string, args ...any) {
	l.logDepth(1, nil, LevelDebug, msg, args...)
}

func (l *Logger) IsDebugEnabled() bool {
	return LevelDebug <= l.Level()
}

func (l *Logger) T(ctx context.Context, msg string, args ...any) {
	l.logDepth(1, ctx, LevelTrace, msg, args...)
}

func (l *Logger) Trace(msg string, args ...any) {
	l.logDepth(1, nil, LevelTrace, msg, args...)
}

func (l *Logger) IsTraceEnabled() bool {
	return LevelTrace <= l.Level()
}

func (l *Logger) onConfigure(config *Config) {
	if l.name == RootLoggerName {
		l.SetLevel(config.Level)
		l.prepareHandlers(config.Handlers, false)
	} else {
		if cfg, exists := config.Package[l.name]; exists {
			l.SetLevel(cfg.Level)
			l.prepareHandlers(cfg.Handlers, cfg.UseParentHandlers)
		} else {
			l.SetLevel(l.parent.Level())
			l.prepareHandlers(nil, true)
		}
	}

	for _, child := range l.children {
		child.onConfigure(config)
	}
}

func (l *Logger) prepareHandlers(handlerNames []string, useParentHandlers bool) {
	l.handlers = make(map[string]Handler, 0)

	if useParentHandlers && l.parent != nil {
		for name, handler := range l.parent.handlers {
			l.handlers[name] = handler
		}
	}

	for _, handlerName := range handlerNames {
		if _, ok := l.handlers[strings.TrimSpace(handlerName)]; ok {
			continue
		}

		if handler, ok := handlers[strings.TrimSpace(handlerName)]; ok {
			l.handlers[handlerName] = handler
		}
	}
}

func (l *Logger) expand(msg string, args ...any) (string, int) {
	var buf []byte

	i := 0
	argIndex := 0
	for j := 0; j < len(msg); j++ {
		if msg[j] == '{' && j+1 < len(msg) {
			if buf == nil {
				buf = make([]byte, 0, 2*len(msg))
			}

			buf = append(buf, msg[i:j]...)

			if msg[j+1] == '}' {
				if len(args)-1 < argIndex {
					buf = append(buf, '{')
					buf = append(buf, '}')
				} else {
					appendValue(&buf, args[argIndex], false)
					argIndex++
				}
			} else {
				buf = append(buf, msg[j+1])
			}

			j += 1
			i = j + 1
		}
	}

	if buf == nil {
		return msg, argIndex
	}

	return string(buf) + msg[i:], argIndex
}

func (l *Logger) logDepth(depth int, ctx context.Context, level Level, msg string, args ...any) error {
	if !l.IsLoggable(level) {
		return nil
	}

	arg := 0
	msg, arg = l.expand(msg, args...)
	record := l.makeRecord(ctx, level, msg)

	if arg == len(args)-1 {
		err, isError := args[arg].(error)

		if isError {
			l.includeStackTrace(depth+1, err, &record)
		} else {
			l.includeCaller(depth+1, &record)
		}
	} else {
		l.includeCaller(depth+1, &record)
	}

	for _, handler := range l.handlers {
		if handler.IsLoggable(record) {
			err := handler.Handle(record)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (l *Logger) makeRecord(ctx context.Context, level Level, msg string) Record {
	record := Record{
		Time:       time.Now(),
		Message:    msg,
		Level:      level,
		Context:    ctx,
		LoggerName: l.name,
		Caller:     Caller{},
	}
	return record
}

func (l *Logger) includeCaller(depth int, record *Record) {
	stack := captureStacktrace(depth+1, stackTraceFirst)
	defer stack.free()

	frame, _ := stack.next()

	record.Caller.function = frame.Function
	record.Caller.line = frame.Line
	record.Caller.file = frame.File
}

func (l *Logger) includeStackTrace(depth int, err error, record *Record) {
	var (
		buf           []byte
		errorTypeName = "Error"
	)

	typ := reflect.TypeOf(err)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	if typ.Name() != "errorString" {
		errorTypeName = typ.String()
	}

	stack := captureStacktrace(depth+1, stackTraceFull)
	defer stack.free()

	frame, more := stack.next()

	buf = append(buf, errorTypeName...)
	buf = append(buf, ": "...)
	buf = append(buf, err.Error()...)
	buf = append(buf, "\\n"...)

	formatFrame(&buf, 0, frame)

	record.Caller.function = frame.Function
	record.Caller.line = frame.Line
	record.Caller.file = frame.File

	i := 1
	for frame, more = stack.next(); more; frame, more = stack.next() {
		formatFrame(&buf, i, frame)
	}

	record.StackTrace = string(buf)
}

func I(ctx context.Context, msg string, args ...any) {
	Default().logDepth(1, ctx, LevelInfo, msg, args...)
}

func Info(msg string, args ...any) {
	Default().logDepth(1, nil, LevelInfo, msg, args...)
}

func IsInfoEnabled() bool {
	return Default().IsInfoEnabled()
}

func E(ctx context.Context, err error, args ...any) {
	Default().logDepth(1, ctx, LevelError, err.Error(), args...)
}

func Error(err error, args ...any) {
	Default().logDepth(1, nil, LevelError, err.Error(), args...)
}

func IsErrorEnabled() bool {
	return Default().IsErrorEnabled()
}

func W(ctx context.Context, msg string, args ...any) {
	Default().logDepth(1, ctx, LevelWarn, msg, args...)
}

func Warn(msg string, args ...any) {
	Default().logDepth(1, nil, LevelWarn, msg, args...)
}

func IsWarnEnabled() bool {
	return Default().IsWarnEnabled()
}

func D(ctx context.Context, msg string, args ...any) {
	Default().logDepth(1, ctx, LevelWarn, msg, args...)
}

func Debug(msg string, args ...any) {
	Default().logDepth(1, nil, LevelWarn, msg, args...)
}

func IsDebugEnabled() bool {
	return Default().IsDebugEnabled()
}

func T(ctx context.Context, msg string, args ...any) {
	Default().logDepth(1, ctx, LevelTrace, msg, args...)
}

func Trace(msg string, args ...any) {
	Default().logDepth(1, nil, LevelTrace, msg, args...)
}

func IsTraceEnabled() bool {
	return Default().IsTraceEnabled()
}
