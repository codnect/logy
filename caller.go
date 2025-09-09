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

import "strings"

// Caller represents the caller information of a log entry.
type Caller struct {
	// defined indicates whether the caller information is defined.
	defined bool
	// file is the file name where the log entry was created.
	file string
	// line is the line number in the file where the log entry was created.
	line int
	// function is the name of the function where the log entry was created.
	function string
}

// Defined method returns whether the caller information is defined.
func (c Caller) Defined() bool {
	return c.defined
}

// Name method returns the name of the function.
func (c Caller) Name() string {
	if !c.defined {
		return ""
	}

	lastDot := strings.LastIndexByte(c.function, '.')

	if lastDot < 0 {
		lastDot = 0
	} else {
		lastDot = lastDot + 1
	}

	return c.function[lastDot:]
}

// File method returns the name of the file.
func (c Caller) File() string {
	if !c.defined {
		return ""
	}

	lastSlash := strings.LastIndexByte(c.file, '/')
	return c.file[lastSlash+1:]
}

// Line method returns the line number in the file.
func (c Caller) Line() int {
	if !c.defined {
		return -1
	}

	return c.line
}

// Package method returns the package name of the function.
func (c Caller) Package() string {
	if !c.defined {
		return ""
	}

	lastSlash := strings.LastIndexByte(c.function, '/')

	if lastSlash < 0 {
		lastSlash = 0
	}

	lastDot := strings.IndexByte(c.function[lastSlash:], '.') + lastSlash

	if lastDot == -1 {
		return c.function
	}

	return c.function[:lastDot]
}

// Path method returns the path of the file without the file name.
func (c Caller) Path() string {
	if !c.defined {
		return ""
	}

	lastSlash := strings.LastIndexByte(c.file, '/')

	if lastSlash == -1 {
		return c.file
	}

	return c.file[:lastSlash]
}
