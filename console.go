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

	handler.SetTarget(TargetStderr)
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

func (h *ConsoleHandler) SetTarget(target Target) {
	switch target {
	case TargetStdout:
		h.target.Store(target)
		h.setWriter(os.Stdout)
	case TargetStderr:
		h.target.Store(target)
		h.setWriter(os.Stdout)
	case TargetDiscard:
		h.setWriter(&discarder{})
		h.target.Store(target)
	}
}

func (h *ConsoleHandler) Target() Target {
	return h.target.Load().(Target)
}

func (h *ConsoleHandler) onConfigure(config *ConsoleConfig) error {
	h.SetEnabled(config.Enable)
	h.SetLevel(config.Level)
	h.SetFormat(config.Format)

	h.SetTarget(config.Target)
	h.SetColorEnabled(config.Color)

	h.applyJsonConfig(config.Json)
	return nil
}
