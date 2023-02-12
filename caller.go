package logy

import "strings"

type Caller struct {
	file     string
	line     int
	function string
}

func (c Caller) Name() string {
	lastDot := strings.LastIndexByte(c.function, '.')

	if lastDot < 0 {
		lastDot = 0
	} else {
		lastDot = lastDot + 1
	}

	return c.function[lastDot:]
}

func (c Caller) File() string {
	lastSlash := strings.LastIndexByte(c.file, '/')
	return c.file[lastSlash+1:]
}

func (c Caller) Line() int {
	return c.line
}

func (c Caller) Package() string {

	lastSlash := strings.LastIndexByte(c.function, '/')

	if lastSlash < 0 {
		lastSlash = 0
	}

	lastDot := strings.IndexByte(c.function[lastSlash:], '.') + lastSlash

	return c.function[:lastDot]
}

func (c Caller) Path() string {
	lastSlash := strings.LastIndexByte(c.file, '/')
	return c.file[:lastSlash]
}
