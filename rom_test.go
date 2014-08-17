package i6502

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test8kRoms(t *testing.T) {
	rom, err := NewRom("test/8kb.rom")

	assert.Nil(t, err)
	assert.Equal(t, 0x2000, rom.Size())
	assert.Equal(t, 0x01, rom.Read(0x0000))
	assert.Equal(t, 0xFF, rom.Read(0x2000-1))
}

func TestRomWritePanic(t *testing.T) {
	rom, _ := NewRom("test/8kb.rom")

	// Writing to rom should panic
	assert.Panics(t, func() {
		rom.Write(0x1337, 0x42)
	}, "Writing to Rom should panic")
}

func Test16kRom(t *testing.T) {
	rom, err := NewRom("test/16kb.rom")

	assert.Nil(t, err)
	assert.Equal(t, 0x4000, rom.Size())
	assert.Equal(t, 0x01, rom.Read(0x0000))
	assert.Equal(t, 0xFF, rom.Read(0x4000-1))
}

func TestRomNotFound(t *testing.T) {
	rom, err := NewRom("test/does-not-exists.rom")
	assert.NotNil(t, err)
	assert.Nil(t, rom)
}
