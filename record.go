package slog

import (
	"context"
	"time"
)

type Record struct {
	Time    time.Time
	Level   Level
	Message string
	Context context.Context

	depth int
}
