package logy

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestWriter struct {
	message string
}

func (tw *TestWriter) Write(p []byte) (n int, err error) {
	tw.message = string(p)
	return 0, nil
}

func TestCommonHandler_Handle(t *testing.T) {
}

func TestCommonHandler_SetLevel(t *testing.T) {
	handler := newConsoleHandler()
	handler.SetLevel(LevelTrace)
	assert.Equal(t, LevelTrace, handler.level.Load().(Level))
}

func TestCommonHandler_Level(t *testing.T) {
	handler := newConsoleHandler()
	handler.SetLevel(LevelTrace)
	assert.Equal(t, LevelTrace, handler.Level())
}

func TestCommonHandler_SetEnabled(t *testing.T) {
	handler := newConsoleHandler()
	handler.SetEnabled(true)
	assert.True(t, handler.enabled.Load().(bool))

	handler.SetEnabled(false)
	assert.False(t, handler.enabled.Load().(bool))
}

func TestCommonHandler_IsEnabled(t *testing.T) {
	handler := newConsoleHandler()
	handler.SetEnabled(true)
	assert.True(t, handler.IsEnabled())

	handler.SetEnabled(false)
	assert.False(t, handler.IsEnabled())
}

func TestCommonHandler_IsLoggable(t *testing.T) {
	handler := newConsoleHandler()
	handler.SetEnabled(false)
	assert.False(t, handler.IsLoggable(Record{}))

	handler.SetEnabled(true)
	handler.SetLevel(LevelOff)
	assert.False(t, handler.IsLoggable(Record{
		Level: LevelTrace,
	}))
}

func TestCommonHandler_SetFormat(t *testing.T) {
	handler := newConsoleHandler()
	handler.SetFormat("anyFormat")
	assert.Equal(t, "anyFormat", handler.format.Load().(string))
}

func TestCommonHandler_Format(t *testing.T) {
	handler := newConsoleHandler()
	handler.SetFormat("anyFormat")
	assert.Equal(t, "anyFormat", handler.Format())
}

func TestCommonHandler_SetJsonEnabled(t *testing.T) {
	handler := newConsoleHandler()
	handler.SetJsonEnabled(true)
	assert.True(t, handler.json.Load().(bool))

	handler.SetJsonEnabled(false)
	assert.False(t, handler.json.Load().(bool))
}

func TestCommonHandler_IsJsonEnabled(t *testing.T) {
	handler := newConsoleHandler()
	handler.SetJsonEnabled(true)
	assert.True(t, handler.IsJsonEnabled())

	handler.SetJsonEnabled(false)
	assert.False(t, handler.IsJsonEnabled())
}

func TestCommonHandler_SetExcludedKeys(t *testing.T) {
	handler := newConsoleHandler()
	excludedKeys := ExcludedKeys{
		"timestamp", "anyKey1",
	}
	handler.SetExcludedKeys(excludedKeys)
	assert.Equal(t, allKeysEnabled^timestampKeyEnabled, handler.enabledKeys.Load().(int))
	assert.Equal(t, ExcludedKeys{"anyKey1"}, handler.excludedFields.Load().(ExcludedKeys))
}

func TestCommonHandler_ExcludedKeys(t *testing.T) {
	handler := newConsoleHandler()
	excludedKeys := ExcludedKeys{
		"timestamp", "anyKey1",
	}
	handler.SetExcludedKeys(excludedKeys)
	assert.Equal(t, allKeysEnabled^timestampKeyEnabled, handler.enabledKeys.Load().(int))
	assert.Equal(t, excludedKeys, handler.ExcludedKeys())
}

func TestCommonHandler_SetKeyOverrides(t *testing.T) {
	handler := newConsoleHandler()
	keyOverrides := KeyOverrides{
		"anyKey1": "anyValue1",
	}
	handler.SetKeyOverrides(keyOverrides)
	assert.Equal(t, keyOverrides, handler.keyOverrides.Load().(KeyOverrides))
}

func TestCommonHandler_KeyOverrides(t *testing.T) {
	handler := newConsoleHandler()
	keyOverrides := KeyOverrides{
		"anyKey1": "anyValue1",
	}
	handler.SetKeyOverrides(keyOverrides)
	assert.Equal(t, keyOverrides, handler.KeyOverrides())
}

func TestCommonHandler_SetAdditionalFields(t *testing.T) {
	handler := newConsoleHandler()
	additionalFields := AdditionalFields{
		"anyKey1": "anyValue1",
	}
	handler.SetAdditionalFields(additionalFields)
	assert.Equal(t, "\"anyKey1\":\"anyValue1\"", handler.additionalFieldsJson.Load().(string))
	assert.Equal(t, additionalFields, handler.additionalFields.Load().(AdditionalFields))
}

func TestCommonHandler_AdditionalFields(t *testing.T) {
	handler := newConsoleHandler()
	additionalFields := AdditionalFields{
		"anyKey1": "anyValue1",
	}
	handler.SetAdditionalFields(additionalFields)
	assert.Equal(t, additionalFields, handler.AdditionalFields())
}
