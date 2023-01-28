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
	rootLogger    = newLogger("", LevelInfo, nil)
	defaultLogger atomic.Value

	cache         = map[string]*Logger{}
	loggerCacheMu sync.RWMutex
)

func init() {
	defer loggerCacheMu.Unlock()
	loggerCacheMu.Lock()
	cache[""] = rootLogger
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
	return level >= l.Level()
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

func (l *Logger) T(ctx context.Context, msg string, args ...any) {
	l.logDepth(1, ctx, LevelTrace, msg, args...)
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

func (l *Logger) Trace(msg string, args ...any) {
	l.logDepth(1, nil, LevelTrace, msg, args...)
}

func (l *Logger) Attrs(level Level, msg string, attrs ...Attribute) {
	l.logAttrsDepth(1, nil, level, msg, attrs...)
}

func (l *Logger) Log(level Level, msg string, args ...any) {
	l.logDepth(1, nil, level, msg, args...)
}

func (l *Logger) onConfigure(config *Config) {
	if l.name == "" {
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

func (l *Logger) logDepth(depth int, ctx context.Context, level Level, msg string, args ...any) error {
	if !l.IsLoggable(level) {
		return nil
	}

	record := l.makeRecord(depth, ctx, level, msg)

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

func (l *Logger) logAttrsDepth(depth int, ctx context.Context, level Level, msg string, attrs ...Attribute) error {
	if !l.IsLoggable(level) {
		return nil
	}

	record := l.makeRecord(depth, ctx, level, msg)

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

func (l *Logger) makeRecord(depth int, ctx context.Context, level Level, msg string) Record {
	var pcs [1]uintptr
	runtime.Callers(depth+3, pcs[:])

	/*frames := runtime.CallersFrames(pcs[:1])
	frame, _ := frames.Next()

	*/
	return Record{
		Time:       time.Now(),
		Message:    msg,
		Level:      level,
		Context:    ctx,
		LoggerName: l.name,
		/*Caller: Caller{
			Defined:  frame.PC != 0,
			PC:       frame.PC,
			File:     frame.File,
			Line:     frame.Line,
			Function: frame.Function,
		},*/
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

func T(ctx context.Context, msg string, args ...any) {
	Default().logDepth(1, ctx, LevelTrace, msg, args...)
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

func Trace(msg string, args ...any) {
	Default().logDepth(1, nil, LevelTrace, msg, args...)
}

func Attrs(level Level, msg string, attrs ...Attribute) {
	Default().logAttrsDepth(1, nil, level, msg, attrs...)
}

func Log(level Level, msg string, args ...any) {
	Default().logDepth(1, nil, level, msg, args...)
}
