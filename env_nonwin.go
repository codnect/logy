//go:build !windows
// +build !windows

package logy

func detectIfSpecialTermColorSupports(termVal string) bool {
	if termVal == "" {
		return false
	}

	return true
}
