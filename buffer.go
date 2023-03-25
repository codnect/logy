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
	*b = t.AppendFormat(*b, time.RFC3339)
}

func (b *buffer) WriteTimeLayout(t time.Time, layout string) {
	*b = t.AppendFormat(*b, layout)
}

func (b *buffer) WriteIntWidth(i, width int) {
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
