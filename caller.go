package logy

import "strings"

var (
	_pool = newPool()
)

type Caller struct {
	Defined  bool
	PC       uintptr
	File     string
	Line     int
	Function string
}

func (c Caller) Path() string {
	if !c.Defined {
		return "undefined"
	}
	buf := _pool.Get()
	buf.AppendString(c.Function)
	str := buf.String()
	buf.Free()
	return str
}

func (c Caller) Name() string {
	if !c.Defined {
		return "undefined"
	}

	buf := _pool.Get()

	lastDot := strings.LastIndexByte(c.Function, '.')
	if lastDot < 0 {
		lastDot = 0
	} else {
		lastDot = lastDot + 1
	}

	buf.AppendString(c.Function[lastDot:])
	name := buf.String()
	buf.Free()

	return name
}

func (c Caller) Package() string {
	if !c.Defined {
		return "undefined"
	}

	buf := _pool.Get()

	lastSlash := strings.LastIndexByte(c.Function, '/')
	if lastSlash < 0 {
		lastSlash = 0
	}
	lastDot := strings.IndexByte(c.Function[lastSlash:], '.') + lastSlash

	buf.AppendString(c.Function[:lastDot])
	pkg := buf.String()
	buf.Free()

	return pkg
}

func (c Caller) TrimmedPath() string {
	if !c.Defined {
		return "undefined"
	}

	idx := strings.LastIndexByte(c.Function, '/')
	if idx == -1 {
		return c.Path()
	}

	idx = strings.LastIndexByte(c.Function[:idx], '/')
	if idx == -1 {
		return c.Path()
	}

	buf := _pool.Get()
	buf.AppendString(c.Function[idx+1:])
	caller := buf.String()
	buf.Free()
	return caller
}
