package logy

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestSupportsColor_IfItIsWslDistro(t *testing.T) {
	os.Clearenv()
	os.Setenv("WSL_DISTRO_NAME", "anyWslDistroName")
	assert.False(t, supportsColor())
}

func TestSupportsColor_ReturnTrueIfItIsJetBrainsTerminal(t *testing.T) {
	os.Clearenv()
	os.Setenv("TERM", "notScreen")
	os.Setenv("TERMINAL_EMULATOR", "JetBrains-JediTerm")
	assert.True(t, supportsColor())
}

func TestSupportsColor_IfItIsSpecialTerm(t *testing.T) {
	os.Clearenv()
	os.Setenv("TERM", "screen")
	os.Setenv("TERMINAL_EMULATOR", "")
	os.Setenv("COLORTERM", "")
	os.Setenv("TERM_PROGRAM", "")
	os.Setenv("FORCE_COLOR", "")
	assert.True(t, supportsColor())
}

func TestSupportsColor_IfItIsAppleTerminal(t *testing.T) {
	os.Clearenv()
	os.Setenv("TERM", "screen")
	os.Setenv("TERMINAL_EMULATOR", "")
	os.Setenv("COLORTERM", "")
	os.Setenv("TERM_PROGRAM", "Apple_Terminal")
	os.Setenv("FORCE_COLOR", "")
	assert.True(t, supportsColor())
}

func TestSupportsColor_IfItIsTerminus(t *testing.T) {
	os.Clearenv()
	os.Setenv("TERM", "screen")
	os.Setenv("TERMINAL_EMULATOR", "")
	os.Setenv("COLORTERM", "")
	os.Setenv("TERM_PROGRAM", "Terminus")
	os.Setenv("FORCE_COLOR", "")
	assert.True(t, supportsColor())
}

func TestSupportsColor_IfItIsHyper(t *testing.T) {
	os.Clearenv()
	os.Setenv("TERM", "screen")
	os.Setenv("TERMINAL_EMULATOR", "")
	os.Setenv("COLORTERM", "")
	os.Setenv("TERM_PROGRAM", "Hyper")
	os.Setenv("FORCE_COLOR", "")
	assert.True(t, supportsColor())
}

func TestSupportsColor_IfItIsITerm(t *testing.T) {
	os.Clearenv()
	os.Setenv("TERM", "screen")
	os.Setenv("TERMINAL_EMULATOR", "")
	os.Setenv("COLORTERM", "")
	os.Setenv("TERM_PROGRAM", "iTerm.app")
	os.Setenv("FORCE_COLOR", "")
	assert.True(t, supportsColor())
}

func TestSupportsColor_IfItIsTrueColor(t *testing.T) {
	os.Clearenv()
	os.Setenv("TERM", "screen")
	os.Setenv("TERMINAL_EMULATOR", "")
	os.Setenv("COLORTERM", "truecolor")
	os.Setenv("TERM_PROGRAM", "")
	os.Setenv("FORCE_COLOR", "")
	assert.True(t, supportsColor())
}
