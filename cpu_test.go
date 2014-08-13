package i6502

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// Creates a new machine, returning the different parts
func NewRamMachine() (*Cpu, *AddressBus, *Ram) {
	ram, _ := NewRam(0x100000) // Full 64kB
	bus, _ := NewAddressBus()
	bus.Attach(ram, 0x0000)
	cpu, _ := NewCpu(bus)

	return cpu, bus, ram
}

func TestNewCpu(t *testing.T) {
	cpu, err := NewCpu(nil)

	assert.NotNil(t, cpu)
	assert.Nil(t, err)
}

func TestCpuAddressBus(t *testing.T) {
	assert := assert.New(t)

	cpu, _, _ := NewRamMachine()
	assert.True(cpu.HasAddressBus())
}

func TestCpuState(t *testing.T) {
	assert := assert.New(t)

	cpu, _, _ := NewRamMachine()

	cpu.bus.Write(0xFFFC, 0x34)
	cpu.bus.Write(0xFFFD, 0x12)

	cpu.Reset()

	// **1101** is specified, but we are satisfied with
	// 00110100 here.
	assert.Equal(0x34, cpu.P)

	// Read PC from $FFFC-FFFD
	assert.Equal(0x1234, cpu.PC)
}
