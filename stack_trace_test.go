package logy

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCaptureFirstStacktrace(t *testing.T) {
	var (
		buf []byte
	)

	stack := captureStacktrace(0, stackTraceFirst)
	frame, more := stack.next()

	formatFrame(&buf, 0, frame)

	i := 1
	for frame, more = stack.next(); more; frame, more = stack.next() {
		formatFrame(&buf, i, frame)
	}

	assert.Equal(t, "github.com/procyon-projects/logy.TestCaptureFirstStacktrace()\\n    /Users/burakkoken/GolandProjects/slog/stack_trace_test.go:13", string(buf))
}

func TestCaptureFullStacktrace(t *testing.T) {
	var (
		buf []byte
	)

	stack := captureStacktrace(0, stackTraceFull)
	frame, more := stack.next()

	formatFrame(&buf, 0, frame)

	i := 1
	for frame, more = stack.next(); more; frame, more = stack.next() {
		formatFrame(&buf, i, frame)
	}

	assert.Equal(t, "github.com/procyon-projects/logy.TestCaptureFullStacktrace()\\n    /Users/burakkoken/GolandProjects/slog/stack_trace_test.go:31", string(buf))
}
