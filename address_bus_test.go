package i6502

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEmptyAddressBus(t *testing.T) {
	assert := assert.New(t)

	bus, err := NewAddressBus()

	assert.Nil(err)

	if assert.NotNil(bus) {
		assert.Equal(0, bus.AddressablesCount())
	}
}

func TestAttachToAddressBus(t *testing.T) {
	assert := assert.New(t)

	bus, _ := NewAddressBus()
	ram, _ := NewRam(0x10000)

	bus.Attach(ram, 0x0000)
	assert.Equal(1, bus.AddressablesCount())
}

func TestBusReadWrite(t *testing.T) {
	assert := assert.New(t)

	bus, _ := NewAddressBus()
	ram, _ := NewRam(0x8000)
	ram2, _ := NewRam(0x8000)
	bus.Attach(ram, 0x0000)
	bus.Attach(ram2, 0x8000)

	// 8-bit Writing
	bus.Write(0x1234, 0xFA)
	assert.Equal(0xFA, ram.Read(0x1234))

	// 16-bit Writing
	bus.Write16(0x1000, 0xAB42)
	assert.Equal(0x42, ram.Read(0x1000))
	assert.Equal(0xAB, ram.Read(0x1001))

	// 8-bit Reading
	ram.Write(0x5522, 0xDA)
	assert.Equal(0xDA, bus.Read(0x5522))

	// 16-bit Reading
	ram.Write(0x4440, 0x7F)
	ram.Write(0x4441, 0x56)
	assert.Equal(0x567F, bus.Read16(0x4440))

	//// Test addressing memory not mounted at 0x0000

	// Read from relative addressable Ram2: $C123
	ram2.Write(0x4123, 0xEF)
	assert.Equal(0xEF, bus.Read(0xC123))

	bus.Write(0x8001, 0x12)
	assert.Equal(0x12, ram2.Read(0x0001))
}
