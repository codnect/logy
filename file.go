package logy

import (
	"os"
	"path/filepath"
	"sync/atomic"
)

type FileHandler struct {
	commonHandler
	name atomic.Value
	path atomic.Value
}

func NewFileHandler() *FileHandler {
	handler := &FileHandler{}
	handler.initializeHandler()

	handler.SetEnabled(false)
	handler.SetLevel(LevelInfo)
	return handler
}

func (h *FileHandler) setFileName(name string) {
	h.name.Store(name)
}

func (h *FileHandler) FileName() string {
	return h.name.Load().(string)
}

func (h *FileHandler) setFilePath(path string) {
	h.path.Store(path)
}

func (h *FileHandler) FilePath() string {
	return h.path.Load().(string)
}

func (h *FileHandler) createLogFile(dir, name string) (*os.File, error) {
	path := filepath.FromSlash(filepath.Join(dir, name))
	return os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0755)
}

func (h *FileHandler) OnConfigure(config Config) error {
	h.SetEnabled(config.File.Enabled)
	h.SetLevel(config.File.Level)
	h.SetFormat(config.File.Format)

	h.setFileName(config.File.Name)
	h.setFilePath(config.File.Path)

	file, err := h.createLogFile(config.File.Path, config.File.Name)
	if err != nil {
		h.SetEnabled(false)
		return err
	}

	h.setWriter(newSyncWriter(file))

	h.applyJsonConfig(config.File.Json)
	return nil
}
