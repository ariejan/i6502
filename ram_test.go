package i6502

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRamSize(t *testing.T) {
	assert := assert.New(t)

	ram, _ := NewRam(0x8000) // 32 kB
	assert.Equal(0x8000, ram.Size())
}

func TestRamReadWrite(t *testing.T) {
	assert := assert.New(t)
	ram, _ := NewRam(0x8000) // 32 kB

	// Ram zeroed out initially
	for i := 0; i < 0x8000; i++ {
		assert.Equal(0x00, ram.data[i])
	}

	ram.Write(0x1000, 0x42)
	assert.Equal(0x42, ram.Read(0x1000))
}
