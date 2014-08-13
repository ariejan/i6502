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

	cpu.Reset()

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

func TestCpuReset(t *testing.T) {
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

func TestProgramLoading(t *testing.T) {
	assert := assert.New(t)

	program := []byte{0xEA, 0xEB, 0xEC}

	cpu, bus, _ := NewRamMachine()
	cpu.LoadProgram(program, 0x0300)

	assert.Equal(0xEA, bus.Read(0x0300))
	assert.Equal(0xEB, bus.Read(0x0301))
	assert.Equal(0xEC, bus.Read(0x0302))

	assert.Equal(0x0300, cpu.PC)
}
