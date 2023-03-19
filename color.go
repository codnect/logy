package logy

import (
	"strconv"
)

type logColor int

const (
	colorBlack logColor = iota + 30
	colorRed
	colorGreen
	colorYellow
	colorBlue
	colorMagenta
	colorCyan
	colorWhite
	colorDefault logColor = 39
)

var (
	supportColor = supportsColor()

	levelColors = []logColor{
		colorRed,
		colorYellow,
		colorGreen,
		colorBlue,
		colorMagenta,
	}
)

func SupportsColor() bool {
	return supportColor
}

func (c logColor) String() string {
	return strconv.FormatInt(int64(c), 10)
}

func (c logColor) print(buf *buffer, value string) {
	if supportColor {
		buf.WriteString("\x1b[")
		buf.WriteString(c.String())
		buf.WriteByte('m')
	}

	buf.WriteString(value)

	if supportColor {
		buf.WriteString("\x1b[0m")
	}
}

func (c logColor) start(buf *buffer) {
	if supportColor {
		*buf = append(*buf, "\x1b["...)
		*buf = append(*buf, c.String()...)
		*buf = append(*buf, 'm')
	}
}

func (c logColor) end(buf *buffer) {
	if supportColor {
		*buf = append(*buf, "\x1b[0m"...)
	}
}
