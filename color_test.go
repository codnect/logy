package logy

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSupportsColor(t *testing.T) {
	assert.Equal(t, supportColor, SupportsColor())
}

func TestLogColor_String(t *testing.T) {
	for _, color := range levelColors {
		assert.Equal(t, fmt.Sprintf("%d", color.value()), color.String())
	}
}

func TestLogColor_print(t *testing.T) {
	buf := newBuffer()
	for _, color := range levelColors {
		color.print(buf, "anyValue")
		assert.Equal(t, fmt.Sprintf("\x1b[%dm%s\x1b[0m", color.value(), "anyValue"), buf.String())
		buf.Reset()
	}
}

func TestColor_startAndEnd(t *testing.T) {
	buf := newBuffer()
	for _, color := range levelColors {
		color.start(buf)
		buf.WriteString("anyValue")
		color.end(buf)
		assert.Equal(t, fmt.Sprintf("\x1b[%dm%s\x1b[0m", color.value(), "anyValue"), buf.String())
		buf.Reset()
	}
}
