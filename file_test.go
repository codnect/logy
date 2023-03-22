package logy

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFileHandler_FileName(t *testing.T) {
	cfg := Config{
		File: &FileConfig{
			Name: "anyName",
			Path: "anyPath",
		},
	}

	handler := newFileHandler(true)
	handler.OnConfigure(cfg)
	assert.Equal(t, cfg.File.Name, handler.FileName())
}

func TestFileHandler_FilePath(t *testing.T) {
	cfg := Config{
		File: &FileConfig{
			Name: "anyName",
			Path: "anyPath",
		},
	}

	handler := newFileHandler(true)
	handler.OnConfigure(cfg)
	assert.Equal(t, cfg.File.Path, handler.FilePath())
}

func TestFileHandler_OnConfigure(t *testing.T) {
	cfg := Config{
		File: &FileConfig{
			Enabled: true,
			Level:   LevelTrace,
			Format:  "anyFormat",
			Name:    "anyName",
			Path:    "anyPath",
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

	handler := newFileHandler(true)
	handler.OnConfigure(cfg)

	assert.Equal(t, cfg.File.Enabled, handler.IsEnabled())
	assert.Equal(t, cfg.File.Level, handler.Level())
	assert.Equal(t, cfg.File.Format, handler.Format())
	assert.Equal(t, cfg.File.Name, handler.FileName())
	assert.Equal(t, cfg.File.Path, handler.FilePath())

	assert.Equal(t, cfg.File.Json.Enabled, handler.IsJsonEnabled())
	assert.Equal(t, cfg.File.Json.ExcludedKeys, handler.ExcludedKeys())
	assert.Equal(t, cfg.File.Json.KeyOverrides, handler.KeyOverrides())
	assert.Equal(t, cfg.File.Json.AdditionalFields, handler.AdditionalFields())
}
