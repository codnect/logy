package logy

import (
	"strconv"
	"sync"
	"time"
)

type buffer []byte

var bufPool = sync.Pool{
	New: func() any {
		b := make([]byte, 0, 1024)
		return (*buffer)(&b)
	},
}

func newBuffer() *buffer {
	return bufPool.Get().(*buffer)
}

func (b *buffer) WritePadding(n int) {
	if n <= 0 {
		return
	}

	for i := 0; i < n; i++ {
		*b = append(*b, ' ')
	}
}

func (b *buffer) Write(p []byte) (int, error) {
	*b = append(*b, p...)
	return len(p), nil
}

func (b *buffer) WriteByte(c byte) {
	*b = append(*b, c)
}

func (b *buffer) WriteString(s string) {
	*b = append(*b, s...)
}

func (b *buffer) WriteInt(i int64) {
	*b = strconv.AppendInt(*b, i, 10)
}

func (b *buffer) WriteTime(t time.Time) {
	b.WriteInt(t.UnixNano())
}

func (b *buffer) WriteTimeAsString(t time.Time) {
	year, month, day := t.Date()
	b.WritePosIntWidth(year, 2)
	b.WriteByte('-')

	b.WritePosIntWidth(int(month), 2)
	b.WriteByte('-')

	b.WritePosIntWidth(day, 2)
	b.WriteByte('T')

	hour, min, sec := t.Clock()
	b.WritePosIntWidth(hour, 2)
	b.WriteByte(':')
	b.WritePosIntWidth(min, 2)
	b.WriteByte(':')
	b.WritePosIntWidth(sec, 2)

	b.WriteByte('.')
	b.WritePosIntWidth(t.Nanosecond()/1e3, 6)
}

func (b *buffer) WritePosInt(i int) {
	b.WritePosIntWidth(i, 0)
}

func (b *buffer) WritePosIntWidth(i, width int) {
	if i < 0 {
		panic("negative int")
	}

	var bb [20]byte
	bp := len(bb) - 1
	for i >= 10 || width > 1 {
		width--
		q := i / 10
		bb[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	bb[bp] = byte('0' + i)
	b.Write(bb[bp:])
}

func (b *buffer) WriteUint(i uint64) {
	*b = strconv.AppendUint(*b, i, 10)
}

func (b *buffer) WriteBool(v bool) {
	*b = strconv.AppendBool(*b, v)
}

func (b *buffer) WriteFloat(f float64, bitSize int) {
	*b = strconv.AppendFloat(*b, f, 'f', -1, bitSize)
}

func (b *buffer) Len() int {
	return len(*b)
}

func (b *buffer) Cap() int {
	return cap(*b)
}

func (b *buffer) Bytes() []byte {
	return *b
}

func (b *buffer) String() string {
	return string(*b)
}

func (b *buffer) Reset() {
	*b = (*b)[:0]
}

func (b *buffer) Free() {
	const maxBufferSize = 16 << 10
	if cap(*b) <= maxBufferSize {
		*b = (*b)[:0]
		bufPool.Put(b)
	}
}
