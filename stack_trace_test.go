package logy

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"runtime"
	"testing"
)

var (
	_, stackTraceFilename, _, _ = runtime.Caller(0)
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

	assert.Equal(t, fmt.Sprintf("github.com/procyon-projects/logy.TestCaptureFirstStacktrace()\\n    %s:19", stackTraceFilename), string(buf))
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

	assert.Equal(t, fmt.Sprintf("github.com/procyon-projects/logy.TestCaptureFullStacktrace()\\n    %s:37", stackTraceFilename), string(buf))
}
