package logy

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestConsoleHandler_SetColorEnabled(t *testing.T) {
	handler := newConsoleHandler()
	handler.SetColorEnabled(true)
	assert.True(t, handler.color.Load().(bool))

	handler.SetColorEnabled(false)
	assert.False(t, handler.color.Load().(bool))
}

func TestConsoleHandler_IsColorEnabled(t *testing.T) {
	handler := newConsoleHandler()
	handler.SetColorEnabled(true)
	assert.True(t, handler.IsColorEnabled())

	handler.SetColorEnabled(false)
	assert.False(t, handler.IsColorEnabled())
}

func TestConsoleHandler_Target(t *testing.T) {
	handler := newConsoleHandler()
	handler.setTarget(TargetStderr)
	assert.Equal(t, TargetStderr, handler.Target())
	assert.Equal(t, os.Stderr, handler.writer.(*syncWriter).writer)

	handler.setTarget(TargetStdout)
	assert.Equal(t, TargetStdout, handler.Target())
	assert.Equal(t, os.Stdout, handler.writer.(*syncWriter).writer)

	handler.setTarget(TargetDiscard)
	assert.Equal(t, TargetDiscard, handler.Target())
	assert.IsType(t, &os.File{}, handler.writer.(*syncWriter).writer)
}

func TestConsoleHandler_OnConfigure(t *testing.T) {
	cfg := Config{
		Console: &ConsoleConfig{
			Enabled: true,
			Level:   LevelTrace,
			Format:  "anyFormat",
			Target:  TargetStdout,
			Color:   true,
			Json: &JsonConfig{
				Enabled:      true,
				ExcludedKeys: ExcludedKeys{"timestamp", "anyExcludedKey1", "anyExcludedKey2"},
				KeyOverrides: KeyOverrides{
					"anyKey3": "anyValue3",
				},
				AdditionalFields: AdditionalFields{
					"anyAdditionalField1": "anyAdditionalValue1",
					"anyAdditionalField2": 41,
				},
			},
		},
	}

	handler := newConsoleHandler()
	handler.OnConfigure(cfg)

	assert.Equal(t, cfg.Console.Enabled, handler.IsEnabled())
	assert.Equal(t, cfg.Console.Level, handler.Level())
	assert.Equal(t, cfg.Console.Format, handler.Format())
	assert.Equal(t, cfg.Console.Target, handler.Target())
	assert.Equal(t, cfg.Console.Color, handler.IsColorEnabled())

	assert.Equal(t, cfg.Console.Json.Enabled, handler.IsJsonEnabled())
	assert.Equal(t, cfg.Console.Json.ExcludedKeys, handler.ExcludedKeys())
	assert.Equal(t, cfg.Console.Json.KeyOverrides, handler.KeyOverrides())
	assert.Equal(t, cfg.Console.Json.AdditionalFields, handler.AdditionalFields())
}
