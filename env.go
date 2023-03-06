package logy

import (
	"github.com/xo/terminfo"
	"os"
	"runtime"
	"strings"
)

var (
	detectedWSL bool
)

func supportsColor() bool {
	if val := os.Getenv("WSL_DISTRO_NAME"); val != "" {
		if detectWSL() {
			return true
		}
	}

	isWin := runtime.GOOS == "windows"
	termVal := os.Getenv("TERM")

	if termVal != "screen" {
		val := os.Getenv("TERMINAL_EMULATOR")
		if val == "JetBrains-JediTerm" {
			return true
		}
	}

	if detectIfSupportsColorFromEnv(termVal, isWin) {
		return true
	}

	return detectIfSpecialTermColorSupports(termVal)
}

func detectWSL() bool {
	if !detectedWSL {
		detectedWSL = true

		b := make([]byte, 1024)
		f, err := os.Open("/proc/version")
		if err == nil {
			_, _ = f.Read(b)
			if err = f.Close(); err != nil {
			}

			wslContents := string(b)
			return strings.Contains(wslContents, "Microsoft")
		}
	}

	return false
}

func detectIfSupportsColorFromEnv(termVal string, isWin bool) bool {
	colorTerm, termProg, forceColor := os.Getenv("COLORTERM"), os.Getenv("TERM_PROGRAM"), os.Getenv("FORCE_COLOR")
	switch {
	case strings.Contains(colorTerm, "truecolor") || strings.Contains(colorTerm, "24bit"):
		return true
	case colorTerm != "" || forceColor != "":
		return true
	case termProg == "Apple_Terminal":
		return true
	case termProg == "Terminus" || termProg == "Hyper":
		return true
	case termProg == "iTerm.app":
		return true
	}

	if !isWin && termVal != "" {
		ti, err := terminfo.Load(termVal)
		if err != nil {
			return false
		}

		v, ok := ti.Nums[terminfo.MaxColors]
		switch {
		case !ok || v <= 16:
			return false
		case ok && v >= 256:
			return true
		}
		return true
	}

	return false
}
