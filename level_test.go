package logy

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLevel_MarshalJSON(t *testing.T) {
	for levelText, level := range levelValues {
		actual, _ := level.MarshalJSON()
		assert.Equal(t, fmt.Sprintf("\"%s\"", levelText), string(actual))
	}
}

func TestLevel_UnmarshalJSON(t *testing.T) {
	for levelText, level := range levelValues {
		var actualLevel Level
		actualLevel.UnmarshalJSON([]byte(fmt.Sprintf("\"%s\"", levelText)))
		assert.Equal(t, level, actualLevel)
	}
}
