package logy

import (
	"os"
)

type ConsoleHandler struct {
	commonHandler
}

func NewConsoleHandler() *ConsoleHandler {
	handler := &ConsoleHandler{}
	handler.initializeHandler()

	handler.setTarget(TargetStderr)
	handler.setWriter(os.Stderr)

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
	switch target {
	case TargetStdout:
		h.target.Store(target)
		h.setWriter(newSyncWriter(os.Stdout))
	case TargetStderr:
		h.target.Store(target)
		h.setWriter(newSyncWriter(os.Stdout))
	case TargetDiscard:
		h.setWriter(&discarder{})
		h.target.Store(target)
	}
}

func (h *ConsoleHandler) Target() Target {
	return h.target.Load().(Target)
}

func (h *ConsoleHandler) OnConfigure(config Config) error {
	h.SetEnabled(config.Console.Enable)
	h.SetLevel(config.Console.Level)
	h.SetFormat(config.Console.Format)

	h.setTarget(config.Console.Target)
	h.SetColorEnabled(config.Console.Color)

	h.applyJsonConfig(config.Console.Json)
	return nil
}
