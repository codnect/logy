package benchmark

import (
	"github.com/procyon-projects/slog"
	"log"
	"testing"
)

func TestL(t *testing.T) {
	l := slog.Default()
	slog.SetDefault(l)

	l.Info("hello")
	slog.Info("v2")

	log.Print("test")
}
