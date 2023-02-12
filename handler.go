package logy

import "sync"

var (
	handlers  = map[string]Handler{}
	handlerMu sync.RWMutex
)

type Handler interface {
	Handle(record Record) error
	SetLevel(level Level)
	Level() Level
	SetEnabled(enabled bool)
	IsEnabled() bool
	IsLoggable(record Record) bool
}

type ConfigurableHandler interface {
	OnConfigure(properties ConfigProperties)
}

func RegisterHandler(name string, handler Handler) {
	defer handlerMu.Unlock()
	handlerMu.Lock()

	handlers[name] = handler
}
