package i6502

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyAddressBus(t *testing.T) {
	assert := assert.New(t)

	bus, err := NewAddressBus()

	assert.Nil(err)

	if assert.NotNil(bus) {
		assert.EqualValues(0, len(bus.addressables))
	}
}

func TestAttachToAddressBus(t *testing.T) {
	assert := assert.New(t)

	bus, _ := NewAddressBus()
	ram, _ := NewRam(0x10000)

	bus.Attach(ram, 0x0000)
	assert.EqualValues(1, len(bus.addressables))
}

func TestBusReadWrite(t *testing.T) {
	assert := assert.New(t)

	bus, _ := NewAddressBus()
	ram, _ := NewRam(0x8000)
	ram2, _ := NewRam(0x8000)
	bus.Attach(ram, 0x0000)
	bus.Attach(ram2, 0x8000)

	// 8-bit Writing
	bus.WriteByte(0x1234, 0xFA)
	assert.EqualValues(0xFA, ram.ReadByte(0x1234))

	// 16-bit Writing
	bus.Write16(0x1000, 0xAB42)
	assert.EqualValues(0x42, ram.ReadByte(0x1000))
	assert.EqualValues(0xAB, ram.ReadByte(0x1001))

	// 8-bit Reading
	ram.WriteByte(0x5522, 0xDA)
	assert.EqualValues(0xDA, bus.ReadByte(0x5522))

	// 16-bit Reading
	ram.WriteByte(0x4440, 0x7F)
	ram.WriteByte(0x4441, 0x56)
	assert.EqualValues(0x567F, bus.Read16(0x4440))

	//// Test addressing memory not mounted at 0x0000

	// Read from relative addressable Ram2: $C123
	ram2.WriteByte(0x4123, 0xEF)
	assert.EqualValues(0xEF, bus.ReadByte(0xC123))

	bus.WriteByte(0x8001, 0x12)
	assert.EqualValues(0x12, ram2.ReadByte(0x0001))
}
