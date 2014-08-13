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

//// NOP

func TestNOP(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0xEA}, 0x0300)
	cpu.Step()
	assert.Equal(t, 0x0301, cpu.PC)
}

//// ADC

func TestADCImmediate(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x69, 0x53}, 0x0300)
	cpu.A = 0x42

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0x95, cpu.A)
	assert.False(t, cpu.getStatus(sCarry))
}

func TestADCWithCarry(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x69, 0x53}, 0x0300)
	cpu.A = 0xC0

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0x13, cpu.A)
	assert.True(t, cpu.getStatus(sCarry))
}

func TestADCWithCarryOver(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x69, 0x04}, 0x0300)

	cpu.setStatus(sCarry, true)
	cpu.A = 0x05

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0x0A, cpu.A)
	assert.False(t, cpu.getStatus(sCarry))
}

func TestADCWithOverflow(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x69, 0xD0}, 0x0300)

	cpu.A = 0x90

	// 0x90 + 0xD0 = 0x160
	// 208 + 144 = 352 => unsigned carry
	// -48 + -112 = 96 => signed overflow
	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0x60, cpu.A)
	assert.True(t, cpu.getStatus(sCarry))
	assert.True(t, cpu.getStatus(sOverflow))
}

func TestADCZero(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x69, 0x00}, 0x0300)
	cpu.A = 0x00

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0x00, cpu.A)
	assert.True(t, cpu.getStatus(sZero))
}

func TestADCNegative(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x69, 0xF7}, 0x0300)
	cpu.A = 0x00

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0xF7, cpu.A)
	assert.True(t, cpu.getStatus(sNegative))
}

func TestADCDecimal(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x69, 0x28}, 0x0300)
	cpu.setStatus(sDecimal, true)
	cpu.A = 0x19

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0x47, cpu.A)
}

func TestADCZeropage(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x65, 0x53}, 0x0300)
	cpu.A = 0x42
	cpu.bus.Write(0x53, 0x12)

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0x54, cpu.A)
}

func TestADCZeropageX(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x75, 0x53}, 0x0300)
	cpu.A = 0x42
	cpu.X = 0x01
	cpu.bus.Write(0x54, 0x12)

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0x54, cpu.A)
}

func TestADCAbsolute(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x6D, 0x00, 0x80}, 0x0300)
	cpu.A = 0x42
	cpu.bus.Write(0x8000, 0x12)

	cpu.Step()

	assert.Equal(t, 0x0303, cpu.PC)
	assert.Equal(t, 0x54, cpu.A)
}

func TestADCAbsoluteX(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x7D, 0x00, 0x80}, 0x0300)
	cpu.A = 0x42
	cpu.X = 0x02
	cpu.bus.Write(0x8002, 0x12)

	cpu.Step()

	assert.Equal(t, 0x0303, cpu.PC)
	assert.Equal(t, 0x54, cpu.A)
}

func TestADCAbsoluteY(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x79, 0x00, 0x80}, 0x0300)
	cpu.A = 0x42
	cpu.Y = 0x02
	cpu.bus.Write(0x8002, 0x12)

	cpu.Step()

	assert.Equal(t, 0x0303, cpu.PC)
	assert.Equal(t, 0x54, cpu.A)
}

func TestADCIndirectX(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x61, 0x80}, 0x0300)
	cpu.A = 0x42
	cpu.X = 0x02
	cpu.bus.Write16(0x82, 0xC000)
	cpu.bus.Write(0xC000, 0x12)

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0x54, cpu.A)
}

func TestADCIndirectY(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x71, 0x80}, 0x0300)
	cpu.A = 0x42
	cpu.Y = 0x02
	cpu.bus.Write16(0x80, 0xC000)
	cpu.bus.Write(0xC002, 0x12)

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0x54, cpu.A)
}

//
