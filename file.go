package logy

import (
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
)

type FileHandler struct {
	writer           atomic.Value
	enabled          atomic.Value
	level            atomic.Value
	path             atomic.Value
	format           atomic.Value
	json             atomic.Value
	excludedKeys     atomic.Value
	additionalFields atomic.Value
	writerMu         sync.RWMutex
}

func NewFileHandler() *FileHandler {
	handler := &FileHandler{}

	handler.enabled.Store(false)
	handler.level.Store(LevelInfo)
	handler.json.Store(false)
	handler.excludedKeys.Store(map[string]struct{}{})
	handler.additionalFields.Store(map[string]JsonAdditionalField{})
	return handler
}

func (h *FileHandler) Handle(record Record) error {
	var (
		buf  []byte
		json bool
	)

	json = h.json.Load().(bool)

	if json {
		excludedKeys := h.excludedKeys.Load().(map[string]struct{})
		additionalFields := h.additionalFields.Load().(map[string]JsonAdditionalField)
		formatJson(&buf, record, excludedKeys, additionalFields)
	} else {
		format := h.format.Load().(string)
		formatText(&buf, format, record, false)
	}

	defer h.writerMu.Unlock()
	h.writerMu.Lock()

	file := h.writer.Load().(*os.File)
	_, err := file.Write(buf)

	return err
}

func (h *FileHandler) SetLevel(level Level) {
	h.level.Store(level)
}

func (h *FileHandler) Level() Level {
	return h.level.Load().(Level)
}

func (h *FileHandler) SetEnabled(enabled bool) {
	h.enabled.Store(enabled)
}

func (h *FileHandler) IsEnabled() bool {
	return h.enabled.Load().(bool)
}

func (h *FileHandler) IsLoggable(record Record) bool {
	if !h.IsEnabled() {
		return false
	}

	return record.Level <= h.Level()
}

func (h *FileHandler) onConfigure(config *FileConfig) error {
	h.enabled.Store(config.Enable)
	h.level.Store(config.Level)
	h.path.Store(config.Path)
	h.format.Store(config.Format)

	logFilePath := filepath.FromSlash(filepath.Join(config.Path, config.Name))

	logFile, err := os.OpenFile(logFilePath, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	h.writer.Store(logFile)

	if config.Json != nil {
		h.json.Store(true)

		if len(config.Json.ExcludeKeys) != 0 {
			excludedKeys := map[string]struct{}{}

			for _, key := range config.Json.ExcludeKeys {
				excludedKeys[key] = struct{}{}
			}

			h.excludedKeys.Store(excludedKeys)
		}

		if len(config.Json.AdditionalFields) != 0 {
			h.additionalFields.Store(config.Json.AdditionalFields)
		}
	}

	return nil
}
