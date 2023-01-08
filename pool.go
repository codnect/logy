package logy

import "sync"

type pool struct {
	p *sync.Pool
}

func newPool() pool {
	return pool{p: &sync.Pool{
		New: func() interface{} {
			return &buffer{bs: make([]byte, 0, _size)}
		},
	}}
}

func (p pool) Get() *buffer {
	buf := p.p.Get().(*buffer)
	buf.Reset()
	buf.pool = p
	return buf
}

func (p pool) put(buf *buffer) {
	p.p.Put(buf)
}
