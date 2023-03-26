package logy

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAbbreviateLoggerName_ShouldNotAbbreviateIfLengthOfLoggerNameIsNotGreaterThanTargetLength(t *testing.T) {
	buf := newBuffer()
	abbreviateLoggerName(buf, "github.com/procyon-projects/logy/test", 40, false)
	assert.Equal(t, "github.com/procyon-projects/logy/test   ", buf.String())
}

func TestAbbreviateLoggerName_ShouldAbbreviateIfLengthOfLoggerNameIsGreaterThanTargetLength(t *testing.T) {
	buf := newBuffer()
	abbreviateLoggerName(buf, "github.com/procyon-projects/logy/test/any", 40, false)
	assert.Equal(t, "g.com/procyon-projects/logy/test/any    ", buf.String())
}
