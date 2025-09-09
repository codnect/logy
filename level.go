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

// Level represents the severity level of a log message.
type Level int

// Log levels.
// Higher values indicate more verbose logging.
const (
	// LevelOff disables all logging.
	LevelOff Level = iota + 1
	// LevelError represents error-level messages.
	LevelError
	// LevelWarn represents warning-level messages.
	LevelWarn
	// LevelInfo represents informational messages.
	LevelInfo
	// LevelDebug represents debug-level messages.
	LevelDebug
	// LevelTrace represents trace-level messages.
	LevelTrace
	// LevelAll enables all logging.
	LevelAll
)
