package utils

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadVariableOffset(t *testing.T) {
	reaader := bytes.NewReader([]byte {
		158,
		46,
	})

	assert.Equal(t, int64(4014), ReadVariableOffset(reaader), "Should have correct offset value")

	reaader = bytes.NewReader([]byte {
		26, 
	})
	assert.Equal(t, int64(26), ReadVariableOffset(reaader), "Should have correct offset value")
}


func TestReadVariableSize(t *testing.T) {
	reaader := bytes.NewReader([]byte {
		193,
		1,
	})
	assert.Equal(t, int64(193), ReadVariableSize(reaader), "Should have correct size value")

	reaader = bytes.NewReader([]byte {
		191,
		1,
	})
	assert.Equal(t, int64(191), ReadVariableSize(reaader), "Should have correct size value")
}