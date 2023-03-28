package logy

import "testing"

func TestConfigFileReader(t *testing.T) {
	reader := newConfigFileReader()
	reader.ReadFile("anyFile")
}
