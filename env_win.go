//go:build windows
// +build windows

package logy

import (
	"golang.org/x/sys/windows"
	"os"
)

var (
	winVersion, _, buildNumber = windows.RtlGetNtVersionNumbers()
)

func detectIfSpecialTermColorSupports(termVal string) bool {
	if os.Getenv("ConEmuANSI") == "ON" {
		return true
	}

	if buildNumber < 10586 || winVersion < 10 {
		if os.Getenv("ANSICON") != "" {
			return true
		}

		return false
	}

	return true
}
