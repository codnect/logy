package logy

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCaller_Defined(t *testing.T) {
	caller := Caller{
		defined: true,
	}

	assert.True(t, caller.Defined())

	caller = Caller{
		defined: false,
	}

	assert.False(t, caller.Defined())
}

func TestCaller_NameShouldReturnEmptyStringIfItIsNotDefined(t *testing.T) {
	caller := Caller{
		defined: false,
	}

	assert.Empty(t, caller.Name())
}

func TestCaller_FileShouldReturnEmptyStringIfItIsNotDefined(t *testing.T) {
	caller := Caller{
		defined: false,
	}

	assert.Empty(t, caller.File())
}

func TestCaller_File(t *testing.T) {
	caller := Caller{
		defined: true,
		file:    "/Users/burakkoken/GolandProjects/slog/test/main.go",
	}

	assert.Equal(t, "main.go", caller.File())
}

func TestCaller_PathShouldReturnEmptyStringIfItIsNotDefined(t *testing.T) {
	caller := Caller{
		defined: false,
	}

	assert.Empty(t, caller.Path())
}

func TestCaller_Path(t *testing.T) {
	caller := Caller{
		defined: true,
		file:    "/Users/burakkoken/GolandProjects/slog/test/main.go",
	}

	assert.Equal(t, "/Users/burakkoken/GolandProjects/slog/test", caller.Path())
}

func TestCaller_PackageShouldReturnEmptyStringIfItIsNotDefined(t *testing.T) {
	caller := Caller{
		defined: false,
	}

	assert.Empty(t, caller.Package())
}

func TestCaller_Package(t *testing.T) {
	caller := Caller{
		defined:  true,
		function: "/test/any.main",
	}

	assert.Equal(t, "/test/any", caller.Package())
}
