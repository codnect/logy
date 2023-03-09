package logy

import (
	"context"
	"time"
)

type Record struct {
	Time       time.Time
	Level      Level
	Message    string
	Context    context.Context
	LoggerName string
	StackTrace string
	Error      error
	Caller     Caller
}
