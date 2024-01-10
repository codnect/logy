package logy

import (
	"io"
	"os"
	"path/filepath"
	"sync/atomic"
)

type FileHandler struct {
	commonHandler
	name      atomic.Value
	path      atomic.Value
	underTest atomic.Value
}

func newFileHandler(underTest bool) *FileHandler {
	handler := &FileHandler{}
	handler.initializeHandler()

	handler.SetEnabled(false)
	handler.SetLevel(LevelInfo)
	handler.setWriter(newSyncWriter(nil, true))
	handler.underTest.Store(underTest)
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
	return os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
}

func (h *FileHandler) OnConfigure(config Config) error {
	h.SetEnabled(config.File.Enabled)
	h.SetLevel(config.File.Level)
	h.SetFormat(config.File.Format)

	h.setFileName(config.File.Name)
	h.setFilePath(config.File.Path)

	underTest := h.underTest.Load().(bool)

	var (
		file io.Writer
		err  error
	)

	discarded := false

	if !underTest && h.IsEnabled() {
		file, err = h.createLogFile(config.File.Path, config.File.Name)
		if err != nil {
			h.SetEnabled(false)
			discarded = true
		}
	} else {
		discarded = true
	}

	fileWriter := h.writer.(*syncWriter)

	defer fileWriter.mu.Unlock()
	fileWriter.mu.Lock()
	fileWriter.writer = file
	fileWriter.setDiscarded(discarded)

	h.applyJsonConfig(config.File.Json)
	return err
}
