package logy

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestBuffer_WritePadding(t *testing.T) {
	buf := newBuffer()

	buf.WritePadding(0)
	assert.Equal(t, "", buf.String())

	buf.WritePadding(5)
	assert.Equal(t, "     ", buf.String())
}

func TestBuffer_Write(t *testing.T) {
	buf := newBuffer()

	buf.WriteByte('a')
	assert.Equal(t, "a", buf.String())
}

func TestBuffer_WriteString(t *testing.T) {
	buf := newBuffer()

	buf.WriteString("test")
	assert.Equal(t, "test", buf.String())
}

func TestBuffer_WriteInt(t *testing.T) {
	buf := newBuffer()

	buf.WriteInt(41)
	assert.Equal(t, "41", buf.String())
}

func TestBuffer_WriteTime(t *testing.T) {
	buf := newBuffer()

	now := time.Now()
	expectedTimeString := now.Format(time.RFC3339)

	buf.WriteTime(now)
	assert.Equal(t, expectedTimeString, buf.String())
}

func TestBuffer_WriteTimeLayout(t *testing.T) {
	buf := newBuffer()

	now := time.Now()
	expectedTimeString := now.Format(time.RFC3339)

	buf.WriteTimeLayout(now, time.RFC3339)
	assert.Equal(t, expectedTimeString, buf.String())
}

func TestBuffer_WriteIntWidth(t *testing.T) {
	buf := newBuffer()

	buf.WriteIntWidth(10, 4)
	assert.Equal(t, "0010", buf.String())
}

func TestBuffer_WriteIntWidthPanicsIfItIsInvokedWithNegativeNumber(t *testing.T) {
	buf := newBuffer()

	assert.Panics(t, func() {
		buf.WriteIntWidth(-1, 0)
	})
}

func TestBuffer_WriteUint(t *testing.T) {
	buf := newBuffer()

	buf.WriteUint(141)
	assert.Equal(t, "141", buf.String())
}

func TestBuffer_WriteBool(t *testing.T) {
	buf := newBuffer()

	buf.WriteBool(true)
	assert.Equal(t, "true", buf.String())
	buf.Reset()

	buf.WriteBool(false)
	assert.Equal(t, "false", buf.String())
}

func TestBuffer_WriteFloat(t *testing.T) {
	buf := newBuffer()

	buf.WriteFloat(41.5, 32)
	assert.Equal(t, "41.5", buf.String())
}

func TestBuffer_Len(t *testing.T) {
	buf := newBuffer()

	buf.WriteString("test")
	assert.Equal(t, 4, buf.Len())
}

func TestBuffer_Cap(t *testing.T) {
	buf := newBuffer()

	buf.WriteString("test")
	assert.Equal(t, 1024, buf.Cap())
}

func TestBuffer_Bytes(t *testing.T) {
	buf := newBuffer()

	buf.WriteString("test")
	assert.Equal(t, []byte{'t', 'e', 's', 't'}, buf.Bytes())
}

func TestBuffer_String(t *testing.T) {
	buf := newBuffer()

	buf.WriteString("test")
	assert.Equal(t, "test", buf.String())
}

func TestBuffer_Reset(t *testing.T) {
	buf := newBuffer()

	buf.WriteString("test")
	buf.Reset()
	assert.Empty(t, buf.String())
}

func TestBuffer_Free(t *testing.T) {
	buf := newBuffer()
	buf.Free()
}
