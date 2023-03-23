package logy

import "strings"

type Caller struct {
	defined  bool
	file     string
	line     int
	function string
}

func (c Caller) Defined() bool {
	return c.defined
}

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

func (c Caller) File() string {
	if !c.defined {
		return ""
	}

	lastSlash := strings.LastIndexByte(c.file, '/')
	return c.file[lastSlash+1:]
}

func (c Caller) Line() int {
	if !c.defined {
		return -1
	}

	return c.line
}

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
