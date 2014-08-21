package i6502

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRamAsMemory(t *testing.T) {
	assert.Implements(t, (*Memory)(nil), new(Ram))
}

func TestRamSize(t *testing.T) {
	ram, _ := NewRam(0x8000) // 32 kB
	assert.Equal(t, 0x8000, ram.Size())
}

func TestRamReadWrite(t *testing.T) {
	ram, _ := NewRam(0x8000) // 32 kB

	// Ram zeroed out initially
	for i := 0; i < 0x8000; i++ {
		assert.Equal(t, 0x00, ram.data[i])
	}

	ram.WriteByte(0x1000, 0x42)
	assert.Equal(t, 0x42, ram.ReadByte(0x1000))
}
