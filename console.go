package logy

import (
	"os"
)

type ConsoleHandler struct {
	commonHandler
}

func newConsoleHandler() *ConsoleHandler {
	handler := &ConsoleHandler{}
	handler.initializeHandler()

	handler.setTarget(TargetStderr)

	handler.SetEnabled(true)
	handler.SetLevel(LevelDebug)

	handler.SetColorEnabled(true)
	return handler
}

func (h *ConsoleHandler) SetColorEnabled(enabled bool) {
	h.color.Store(enabled)
}

func (h *ConsoleHandler) IsColorEnabled() bool {
	return h.color.Load().(bool)
}

func (h *ConsoleHandler) setTarget(target Target) {
	var consoleWriter *syncWriter
	if h.writer != nil {
		consoleWriter = h.writer.(*syncWriter)
	} else {
		consoleWriter = newSyncWriter(nil)
		h.writer = consoleWriter
	}

	defer consoleWriter.mu.Unlock()
	consoleWriter.mu.Lock()

	switch target {
	case TargetStdout:
		h.target.Store(target)
		consoleWriter.writer = os.Stdout
	case TargetStderr:
		h.target.Store(target)
		consoleWriter.writer = os.Stderr
	default:
		h.target.Store(target)
		consoleWriter.writer = &discarder{}
	}
}

func (h *ConsoleHandler) Target() Target {
	return h.target.Load().(Target)
}

func (h *ConsoleHandler) OnConfigure(config Config) error {
	h.SetEnabled(config.Console.Enabled)
	h.SetLevel(config.Console.Level)
	h.SetFormat(config.Console.Format)

	h.setTarget(config.Console.Target)
	h.SetColorEnabled(config.Console.Color)

	h.applyJsonConfig(config.Console.Json)
	return nil
}
