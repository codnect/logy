// Copyright 2025 Codnect
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package logy

import (
	"context"
	"sync"
	"time"
)

// Predefined handler names.
// They are reserved and cannot be used for custom handlers.
const (
	// ConsoleHandlerName represents the console handler name.
	ConsoleHandlerName = "console"
	// FileHandlerName represents the file handler name.
	FileHandlerName = "file"
	// SyslogHandlerName represents the syslog handler name.
	SyslogHandlerName = "syslog"
)

var (
	// handlers map stores the registered handlers.
	handlers = map[string]Handler{}
	// handlerMu is a mutex to protect the handler map.
	handlerMu sync.RWMutex
)

// Record struct represents a log record.
// It contains information about the log event such as time, level, message, context,
// logger name, stack trace, error and caller.
// It is used by handlers to process and format log messages.
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

// Handler interface defines the methods that a log handler must implement.
// A handler is responsible for processing log records and outputting them to a specific destination.
type Handler interface {
	// Handle method processes a log record.
	Handle(record Record) error
	// SetLevel method sets the minimum log level for the handler.
	SetLevel(level Level)
	// Level method returns the current log level of the handler.
	Level() Level
	// SetEnabled method enables or disables the handler.
	SetEnabled(enabled bool)
	// IsEnabled method returns whether the handler is enabled.
	IsEnabled() bool
	// IsLoggable method checks whether the log record should be logged based on the handler's level and enabled status.
	IsLoggable(record Record) bool
}

// Register function registers a custom log handler with the given name.
func Register(name string, handler Handler) {
	handlerMu.Lock()
	defer handlerMu.Unlock()

	if name == ConsoleHandlerName || name == FileHandlerName || name == SyslogHandlerName {
		panic("logy: 'console', 'file' and 'syslog' handler names are reserved and cannot be registered")
	}

	if handler == nil {
		panic("logy: nil handler")
	}

	handlers[name] = handler
}
