package logy

import (
	"runtime"
	"strconv"
	"sync"
)

var stackTracePool = sync.Pool{
	New: func() interface{} {
		return &stackTrace{
			storage: make([]uintptr, 64),
		}
	},
}

type stackTrace struct {
	pcs    []uintptr
	frames *runtime.Frames

	storage []uintptr
}

type stacktraceDepth int

const (
	stackTraceFirst stacktraceDepth = iota
	stackTraceFull
)

func captureStacktrace(skip int, depth stacktraceDepth) *stackTrace {
	stack := stackTracePool.Get().(*stackTrace)

	switch depth {
	case stackTraceFirst:
		stack.pcs = stack.storage[:1]
	case stackTraceFull:
		stack.pcs = stack.storage
	}

	numFrames := runtime.Callers(
		skip+2,
		stack.pcs,
	)

	if depth == stackTraceFull {
		pcs := stack.pcs
		for numFrames == len(pcs) {
			pcs = make([]uintptr, len(pcs)*2)
			numFrames = runtime.Callers(skip+2, pcs)
		}

		stack.storage = pcs
		stack.pcs = pcs[:numFrames-1]
	} else {
		stack.pcs = stack.pcs[:numFrames]
	}

	stack.frames = runtime.CallersFrames(stack.pcs)
	return stack
}

func (s *stackTrace) free() {
	s.frames = nil
	s.pcs = nil
	stackTracePool.Put(s)
}

func (s *stackTrace) count() int {
	return len(s.pcs)
}

func (s *stackTrace) next() (_ runtime.Frame, more bool) {
	return s.frames.Next()
}

func formatFrame(buf *[]byte, index int, frame runtime.Frame) {
	if index != 0 {
		*buf = append(*buf, "\\n"...)
	}

	*buf = append(*buf, frame.Function...)
	*buf = append(*buf, '(')
	*buf = append(*buf, ')')
	*buf = append(*buf, "\\n"...)
	*buf = append(*buf, "    "...)
	*buf = append(*buf, frame.File...)
	*buf = append(*buf, ':')
	*buf = append(*buf, strconv.FormatInt(int64(frame.Line), 10)...)
}
