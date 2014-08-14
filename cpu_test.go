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

func TestSEC(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x38}, 0x0300)

	assert.False(t, cpu.getCarry())
	cpu.Step()
	assert.True(t, cpu.getCarry())
}

func TestSED(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0xF8}, 0x0300)

	assert.False(t, cpu.getDecimal())
	cpu.Step()
	assert.True(t, cpu.getDecimal())
}

func TestSEI(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x78}, 0x0300)
	cpu.setInterrupt(false)

	assert.False(t, cpu.getInterrupt())
	cpu.Step()
	assert.True(t, cpu.getInterrupt())
}

func TestCLC(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x18}, 0x0300)
	cpu.setCarry(true)

	assert.True(t, cpu.getStatus(sCarry))
	cpu.Step()
	assert.False(t, cpu.getStatus(sCarry))
}

func TestCLD(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0xD8}, 0x0300)
	cpu.setDecimal(true)

	assert.True(t, cpu.getDecimal())
	cpu.Step()
	assert.False(t, cpu.getDecimal())
}

func TestCLI(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x58}, 0x0300)
	cpu.setInterrupt(true)

	assert.True(t, cpu.getInterrupt())
	cpu.Step()
	assert.False(t, cpu.getInterrupt())
}

func TestCLV(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0xB8}, 0x0300)
	cpu.setOverflow(true)

	assert.True(t, cpu.getOverflow())
	cpu.Step()
	assert.False(t, cpu.getOverflow())
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

//// SBC

func TestSBCImmediate(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0xE9, 0x01}, 0x0300)
	cpu.A = 0x42
	cpu.setCarry(true)

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0x41, cpu.A)
	assert.True(t, cpu.getStatus(sCarry))
}

func TestSBCWithoutCarry(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0xE9, 0x01}, 0x0300)
	cpu.A = 0x42
	cpu.setCarry(false)

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0x40, cpu.A)
	assert.True(t, cpu.getStatus(sCarry))
}

func TestSBCNegativeNoCarry(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0xE9, 0x43}, 0x0300)
	cpu.A = 0x42
	cpu.setCarry(true)

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0xff, cpu.A)
	assert.False(t, cpu.getCarry())
	assert.True(t, cpu.getNegative())
}

func TestSBCDecimal(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0xE9, 0x03}, 0x0300)
	cpu.setStatus(sDecimal, true)
	cpu.A = 0x32

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0x29, cpu.A)
}

func TestSBCZero(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0xE9, 0x42}, 0x0300)
	cpu.A = 0x42
	cpu.setCarry(true)

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0x00, cpu.A)
	assert.True(t, cpu.getZero())
}

func TestSBCZeropage(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0xE5, 0x53}, 0x0300)
	cpu.setCarry(true)
	cpu.A = 0x42
	cpu.bus.Write(0x53, 0x12)

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0x30, cpu.A)
}

func TestZeropageX(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0xF5, 0x53}, 0x0300)
	cpu.setCarry(true)
	cpu.A = 0x42
	cpu.X = 0x01
	cpu.bus.Write(0x54, 0x12)

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0x30, cpu.A)
}

func TestSBCAbsolute(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0xED, 0x00, 0x80}, 0x0300)
	cpu.setCarry(true)
	cpu.A = 0x42
	cpu.bus.Write(0x8000, 0x12)

	cpu.Step()

	assert.Equal(t, 0x0303, cpu.PC)
	assert.Equal(t, 0x30, cpu.A)
}

func TestSBCAbsoluteX(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0xFD, 0x00, 0x80}, 0x0300)
	cpu.setCarry(true)
	cpu.A = 0x42
	cpu.X = 0x02
	cpu.bus.Write(0x8002, 0x12)

	cpu.Step()

	assert.Equal(t, 0x0303, cpu.PC)
	assert.Equal(t, 0x30, cpu.A)
}

func TestSBCAbsoluteY(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0xF9, 0x00, 0x80}, 0x0300)
	cpu.setCarry(true)
	cpu.A = 0x42
	cpu.Y = 0x02
	cpu.bus.Write(0x8002, 0x12)

	cpu.Step()

	assert.Equal(t, 0x0303, cpu.PC)
	assert.Equal(t, 0x30, cpu.A)
}

func TestSBCIndirectX(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0xE1, 0x80}, 0x0300)
	cpu.setCarry(true)
	cpu.A = 0x42
	cpu.X = 0x02
	cpu.bus.Write16(0x82, 0xC000)
	cpu.bus.Write(0xC000, 0x12)

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0x30, cpu.A)
}

func TestSBCIndirectY(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0xF1, 0x80}, 0x0300)
	cpu.setCarry(true)
	cpu.A = 0x42
	cpu.Y = 0x02
	cpu.bus.Write16(0x80, 0xC000)
	cpu.bus.Write(0xC002, 0x12)

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0x30, cpu.A)
}

//// INX

func TestINX(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0xE8}, 0x0300)
	cpu.X = 0x42

	cpu.Step()

	assert.Equal(t, 0x0301, cpu.PC)
	assert.Equal(t, 0x43, cpu.X)
}

func TestINXRollover(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0xE8}, 0x0300)
	cpu.X = 0xFF

	cpu.Step()

	assert.Equal(t, 0x0301, cpu.PC)
	assert.Equal(t, 0x00, cpu.X)
}

//// INY

func TestINY(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0xC8}, 0x0300)
	cpu.Y = 0x42

	cpu.Step()

	assert.Equal(t, 0x0301, cpu.PC)
	assert.Equal(t, 0x43, cpu.Y)
}

func TestINYRollover(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0xC8}, 0x0300)
	cpu.Y = 0xFF

	cpu.Step()

	assert.Equal(t, 0x0301, cpu.PC)
	assert.Equal(t, 0x00, cpu.Y)
}

//// INC

func TestINCZeropage(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0xE6, 0x42}, 0x0300)
	cpu.bus.Write(0x42, 0x01)

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0x02, cpu.bus.Read(0x42))
}

func TestINCZeropageX(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0xF6, 0x42}, 0x0300)
	cpu.X = 0x01
	cpu.bus.Write(0x43, 0x01)

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0x02, cpu.bus.Read(0x43))
}

func TestINCAbsolute(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0xEE, 0x00, 0x80}, 0x0300)
	cpu.bus.Write(0x8000, 0x01)

	cpu.Step()

	assert.Equal(t, 0x0303, cpu.PC)
	assert.Equal(t, 0x02, cpu.bus.Read(0x8000))
}

func TestINCAbsoluteX(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0xFE, 0x00, 0x80}, 0x0300)
	cpu.X = 0x02
	cpu.bus.Write(0x8002, 0x01)

	cpu.Step()

	assert.Equal(t, 0x0303, cpu.PC)
	assert.Equal(t, 0x02, cpu.bus.Read(0x8002))
}

//// DEX

func TestDEX(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0xCA}, 0x0300)
	cpu.X = 0x42

	cpu.Step()

	assert.Equal(t, 0x0301, cpu.PC)
	assert.Equal(t, 0x41, cpu.X)
}

func TestDEXRollover(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0xCA}, 0x0300)
	cpu.X = 0x00

	cpu.Step()

	assert.Equal(t, 0x0301, cpu.PC)
	assert.Equal(t, 0xFF, cpu.X)
}

//// DEY

func TestDEY(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x88}, 0x0300)
	cpu.Y = 0x42

	cpu.Step()

	assert.Equal(t, 0x0301, cpu.PC)
	assert.Equal(t, 0x41, cpu.Y)
}

func TestDEYRollover(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x88}, 0x0300)
	cpu.Y = 0x00

	cpu.Step()

	assert.Equal(t, 0x0301, cpu.PC)
	assert.Equal(t, 0xFF, cpu.Y)
}

//// DEC

func TestDECZeropage(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0xC6, 0x42}, 0x0300)
	cpu.bus.Write(0x42, 0x01)

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0x00, cpu.bus.Read(0x42))
}

func TestDECZeropageX(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0xD6, 0x42}, 0x0300)
	cpu.X = 0x01
	cpu.bus.Write(0x43, 0x01)

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0x00, cpu.bus.Read(0x43))
}

func TestDECAbsolute(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0xCE, 0x00, 0x80}, 0x0300)
	cpu.bus.Write(0x8000, 0x01)

	cpu.Step()

	assert.Equal(t, 0x0303, cpu.PC)
	assert.Equal(t, 0x00, cpu.bus.Read(0x8000))
}

func TestDECAbsoluteX(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0xDE, 0x00, 0x80}, 0x0300)
	cpu.X = 0x02
	cpu.bus.Write(0x8002, 0x01)

	cpu.Step()

	assert.Equal(t, 0x0303, cpu.PC)
	assert.Equal(t, 0x00, cpu.bus.Read(0x8002))
}

//// LDA

func TestLDAImmediate(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0xA9, 0x42}, 0x0300)

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0x42, cpu.A)
}

func TestLDANegative(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0xA9, 0xAE}, 0x0300)

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0xAE, cpu.A)
	assert.True(t, cpu.getNegative())
}

func TestLDAZero(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0xA9, 0x00}, 0x0300)

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0x00, cpu.A)
	assert.True(t, cpu.getZero())
}

func TestLDAZeropage(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0xA5, 0x42}, 0x0300)
	cpu.bus.Write(0x42, 0xF8)

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0xF8, cpu.A)
}

func TestLDAZeropageX(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0xB5, 0x41}, 0x0300)
	cpu.X = 0x01
	cpu.bus.Write(0x42, 0xF8)

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0xF8, cpu.A)
}

func TestLDAAbsolute(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0xAD, 0x00, 0x80}, 0x0300)
	cpu.bus.Write16(0x8000, 0xF8)

	cpu.Step()

	assert.Equal(t, 0x0303, cpu.PC)
	assert.Equal(t, 0xF8, cpu.A)
}

func TestLDAAbsoluteX(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0xBD, 0x00, 0x80}, 0x0300)
	cpu.X = 0x02
	cpu.bus.Write16(0x8002, 0xF8)

	cpu.Step()

	assert.Equal(t, 0x0303, cpu.PC)
	assert.Equal(t, 0xF8, cpu.A)
}

func TestLDAAbsoluteY(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0xB9, 0x00, 0x80}, 0x0300)
	cpu.Y = 0x02
	cpu.bus.Write16(0x8002, 0xF8)

	cpu.Step()

	assert.Equal(t, 0x0303, cpu.PC)
	assert.Equal(t, 0xF8, cpu.A)
}

func TestLDAIndirectX(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0xA1, 0x80}, 0x0300)
	cpu.X = 0x02
	cpu.bus.Write16(0x82, 0xC000)
	cpu.bus.Write(0xC000, 0xF8)

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0xF8, cpu.A)
}

func TestLDAIndirectY(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0xB1, 0x80}, 0x0300)
	cpu.Y = 0x02
	cpu.bus.Write16(0x80, 0xC000)
	cpu.bus.Write(0xC002, 0xF8)

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0xF8, cpu.A)
}

//// LDX

func TestLDXImmediate(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0xA2, 0x42}, 0x0300)

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0x42, cpu.X)
}

func TestLDXNegative(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0xA2, 0xAE}, 0x0300)

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0xAE, cpu.X)
	assert.True(t, cpu.getNegative())
}

func TestLDXZero(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0xA2, 0x00}, 0x0300)

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0x00, cpu.X)
	assert.True(t, cpu.getZero())
}

func TestLDXZeropage(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0xA6, 0x42}, 0x0300)
	cpu.bus.Write(0x42, 0xF8)

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0xF8, cpu.X)
}

func TestLDXZeropageY(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0xB6, 0x41}, 0x0300)
	cpu.Y = 0x01
	cpu.bus.Write(0x42, 0xF8)

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0xF8, cpu.X)
}

func TestLDXAbsolute(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0xAE, 0x00, 0x80}, 0x0300)
	cpu.bus.Write16(0x8000, 0xF8)

	cpu.Step()

	assert.Equal(t, 0x0303, cpu.PC)
	assert.Equal(t, 0xF8, cpu.X)
}

func TestLDXAbsoluteY(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0xBE, 0x00, 0x80}, 0x0300)
	cpu.Y = 0x02
	cpu.bus.Write16(0x8002, 0xF8)

	cpu.Step()

	assert.Equal(t, 0x0303, cpu.PC)
	assert.Equal(t, 0xF8, cpu.X)
}

//// LDY

func TestLDYImmediate(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0xA0, 0x42}, 0x0300)

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0x42, cpu.Y)
}

func TestLDYNegative(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0xA0, 0xAE}, 0x0300)

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0xAE, cpu.Y)
	assert.True(t, cpu.getNegative())
}

func TestLDYZero(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0xA0, 0x00}, 0x0300)

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0x00, cpu.Y)
	assert.True(t, cpu.getZero())
}

func TestLDYZeropage(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0xA4, 0x42}, 0x0300)
	cpu.bus.Write(0x42, 0xF8)

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0xF8, cpu.Y)
}

func TestLDYZeropageY(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0xB4, 0x41}, 0x0300)
	cpu.Y = 0x01
	cpu.bus.Write(0x42, 0xF8)

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0xF8, cpu.Y)
}

func TestLDYAbsolute(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0xAC, 0x00, 0x80}, 0x0300)
	cpu.bus.Write16(0x8000, 0xF8)

	cpu.Step()

	assert.Equal(t, 0x0303, cpu.PC)
	assert.Equal(t, 0xF8, cpu.Y)
}

func TestLDYAbsoluteY(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0xBC, 0x00, 0x80}, 0x0300)
	cpu.Y = 0x02
	cpu.bus.Write16(0x8002, 0xF8)

	cpu.Step()

	assert.Equal(t, 0x0303, cpu.PC)
	assert.Equal(t, 0xF8, cpu.Y)
}

//// ORA

func TestORAImmediate(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x09, 0x02}, 0x0300)
	cpu.A = 0x70

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0x72, cpu.A)
}

func TestORANegative(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x09, 0x02}, 0x0300)
	cpu.A = 0xF0

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0xF2, cpu.A)
	assert.True(t, cpu.getNegative())
}

func TestORAZero(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x09, 0x00}, 0x0300)
	cpu.A = 0x00

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0x00, cpu.A)
	assert.True(t, cpu.getZero())
}

func TestORAZeropage(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x05, 0x42}, 0x0300)
	cpu.A = 0xF0
	cpu.bus.Write(0x0042, 0x02)

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0xF2, cpu.A)
}

func TestORAZeropageX(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x15, 0x40}, 0x0300)
	cpu.A = 0xF0
	cpu.X = 0x02
	cpu.bus.Write(0x0042, 0x02)

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0xF2, cpu.A)
}

func TestORAAbsolute(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x0D, 0x00, 0x80}, 0x0300)
	cpu.A = 0xF0
	cpu.bus.Write(0x8000, 0x02)

	cpu.Step()

	assert.Equal(t, 0x0303, cpu.PC)
	assert.Equal(t, 0xF2, cpu.A)
}

func TestORAAbsoluteX(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x1D, 0x00, 0x80}, 0x0300)
	cpu.A = 0xF0
	cpu.X = 0x02
	cpu.bus.Write(0x8002, 0x02)

	cpu.Step()

	assert.Equal(t, 0x0303, cpu.PC)
	assert.Equal(t, 0xF2, cpu.A)
}

func TestORAAbsoluteY(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x19, 0x00, 0x80}, 0x0300)
	cpu.A = 0xF0
	cpu.Y = 0x02
	cpu.bus.Write(0x8002, 0x02)

	cpu.Step()

	assert.Equal(t, 0x0303, cpu.PC)
	assert.Equal(t, 0xF2, cpu.A)
}

func TestORAIndirectX(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x01, 0x40}, 0x0300)
	cpu.A = 0xF0
	cpu.X = 0x02
	cpu.bus.Write16(0x42, 0xC000)
	cpu.bus.Write(0xC000, 0x02)

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0xF2, cpu.A)
}

func TestORAIndirectY(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x11, 0x40}, 0x0300)
	cpu.A = 0xF0
	cpu.Y = 0x02
	cpu.bus.Write16(0x40, 0xC000)
	cpu.bus.Write(0xC002, 0x02)

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0xF2, cpu.A)
}

//// AND

func TestANDImmediate(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x29, 0x0f}, 0x0300)
	cpu.A = 0x42

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0x02, cpu.A)
}

func TestANDNegative(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x29, 0xF0}, 0x0300)
	cpu.A = 0xA3

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0xA0, cpu.A)
	assert.True(t, cpu.getNegative())
}

func TestANDZero(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x29, 0xf0}, 0x0300)
	cpu.A = 0x0f

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0x00, cpu.A)
	assert.True(t, cpu.getZero())
}

func TestANDZeropage(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x25, 0x42}, 0x0300)
	cpu.A = 0xE9
	cpu.bus.Write(0x0042, 0x0f)

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0x09, cpu.A)
}

func TestANDZeropageX(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x35, 0x40}, 0x0300)
	cpu.A = 0x42
	cpu.X = 0x02
	cpu.bus.Write(0x0042, 0x0F)

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0x02, cpu.A)
}

func TestANDAbsolute(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x2D, 0x00, 0x80}, 0x0300)
	cpu.A = 0x42
	cpu.bus.Write(0x8000, 0x0F)

	cpu.Step()

	assert.Equal(t, 0x0303, cpu.PC)
	assert.Equal(t, 0x02, cpu.A)
}

func TestANDAbsoluteX(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x3D, 0x00, 0x80}, 0x0300)
	cpu.A = 0x42
	cpu.X = 0x02
	cpu.bus.Write(0x8002, 0x0F)

	cpu.Step()

	assert.Equal(t, 0x0303, cpu.PC)
	assert.Equal(t, 0x02, cpu.A)
}

func TestANDAbsoluteY(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x39, 0x00, 0x80}, 0x0300)
	cpu.A = 0x42
	cpu.Y = 0x02
	cpu.bus.Write(0x8002, 0x0F)

	cpu.Step()

	assert.Equal(t, 0x0303, cpu.PC)
	assert.Equal(t, 0x02, cpu.A)
}

func TestANDIndirectX(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x21, 0x80}, 0x0300)
	cpu.A = 0x42
	cpu.X = 0x02
	cpu.bus.Write16(0x82, 0xC000)
	cpu.bus.Write(0xC000, 0x0F)

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0x02, cpu.A)
}

func TestANDIndirectY(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x31, 0x80}, 0x0300)
	cpu.A = 0x42
	cpu.Y = 0x02
	cpu.bus.Write16(0x80, 0xC000)
	cpu.bus.Write(0xC002, 0x0F)

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0x02, cpu.A)
}

//// EOR

func TestEORImmediate(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x49, 0x7f}, 0x0300)
	cpu.A = 0x42

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0x3D, cpu.A)
}

func TestEORNegative(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x49, 0xf7}, 0x0300)
	cpu.A = 0x24

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0xD3, cpu.A)
	assert.True(t, cpu.getNegative())
}

func TestEORZero(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x49, 0xff}, 0x0300)
	cpu.A = 0xff

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0x00, cpu.A)
	assert.True(t, cpu.getZero())
}

func TestEORZeropage(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x45, 0x80}, 0x0300)
	cpu.A = 0x42
	cpu.bus.Write(0x0080, 0x7f)

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0x3D, cpu.A)
}

func TestEORZeropageX(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x55, 0x80}, 0x0300)
	cpu.A = 0x42
	cpu.X = 0x02
	cpu.bus.Write(0x0082, 0x7F)

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0x3D, cpu.A)
}

func TestEORAbsolute(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x4D, 0x00, 0x80}, 0x0300)
	cpu.A = 0x42
	cpu.bus.Write(0x8000, 0x7F)

	cpu.Step()

	assert.Equal(t, 0x0303, cpu.PC)
	assert.Equal(t, 0x3D, cpu.A)
}

func TestEORAbsoluteX(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x5D, 0x00, 0x80}, 0x0300)
	cpu.A = 0x42
	cpu.X = 0x02
	cpu.bus.Write(0x8002, 0x7F)

	cpu.Step()

	assert.Equal(t, 0x0303, cpu.PC)
	assert.Equal(t, 0x3D, cpu.A)
}

func TestEORAbsoluteY(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x59, 0x00, 0x80}, 0x0300)
	cpu.A = 0x42
	cpu.Y = 0x02
	cpu.bus.Write(0x8002, 0x7F)

	cpu.Step()

	assert.Equal(t, 0x0303, cpu.PC)
	assert.Equal(t, 0x3D, cpu.A)
}

func TestEORIndirectX(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x41, 0x80}, 0x0300)
	cpu.A = 0x42
	cpu.X = 0x02
	cpu.bus.Write16(0x82, 0xC000)
	cpu.bus.Write(0xC000, 0x7F)

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0x3D, cpu.A)
}

func TestEORIndirectY(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x51, 0x80}, 0x0300)
	cpu.A = 0x42
	cpu.Y = 0x02
	cpu.bus.Write16(0x80, 0xC000)
	cpu.bus.Write(0xC002, 0x7F)

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0x3D, cpu.A)
}

//// STA

func TestSTAZeropage(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x85, 0x80}, 0x0300)
	cpu.A = 0x42

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0x42, cpu.bus.Read(0x0080))
}

func TestSTAZeropageX(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x95, 0x80}, 0x0300)
	cpu.A = 0x42
	cpu.X = 0x02

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0x42, cpu.bus.Read(0x0082))
}

func TestSTAAbsolute(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x8D, 0x00, 0x80}, 0x0300)
	cpu.A = 0x42

	cpu.Step()

	assert.Equal(t, 0x0303, cpu.PC)
	assert.Equal(t, 0x42, cpu.bus.Read(0x8000))
}

func TestSTAAbsoluteX(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x9D, 0x00, 0x80}, 0x0300)
	cpu.A = 0x42
	cpu.X = 0x02

	cpu.Step()

	assert.Equal(t, 0x0303, cpu.PC)
	assert.Equal(t, 0x42, cpu.bus.Read(0x8002))
}

func TestSTAAbsoluteY(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x99, 0x00, 0x80}, 0x0300)
	cpu.A = 0x42
	cpu.Y = 0x02
	cpu.bus.Write(0x8002, 0x7F)

	cpu.Step()

	assert.Equal(t, 0x0303, cpu.PC)
	assert.Equal(t, 0x42, cpu.bus.Read(0x8002))
}

func TestSTAIndirectX(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x81, 0x80}, 0x0300)
	cpu.A = 0x42
	cpu.X = 0x02
	cpu.bus.Write16(0x82, 0xC000)

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0x42, cpu.bus.Read(0xC000))
}

func TestSTAIndirectY(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x91, 0x80}, 0x0300)
	cpu.A = 0x42
	cpu.Y = 0x02
	cpu.bus.Write16(0x80, 0xC000)

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0x42, cpu.bus.Read(0xC002))
}

//// STX

func TestSTXZeropage(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x86, 0x80}, 0x0300)
	cpu.X = 0x42

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0x42, cpu.bus.Read(0x0080))
}

func TestSTXZeropageY(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x96, 0x80}, 0x0300)
	cpu.X = 0x42
	cpu.Y = 0x02

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0x42, cpu.bus.Read(0x0082))
}

func TestSTXAbsolute(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x8E, 0x00, 0x80}, 0x0300)
	cpu.X = 0x42

	cpu.Step()

	assert.Equal(t, 0x0303, cpu.PC)
	assert.Equal(t, 0x42, cpu.bus.Read(0x8000))
}

//// STY

func TestSTYZeropage(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x84, 0x80}, 0x0300)
	cpu.Y = 0x42

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0x42, cpu.bus.Read(0x0080))
}

func TestSTYZeropageX(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x94, 0x80}, 0x0300)
	cpu.X = 0x02
	cpu.Y = 0x42

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0x42, cpu.bus.Read(0x0082))
}

func TestSTYAbsolute(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x8C, 0x00, 0x80}, 0x0300)
	cpu.Y = 0x42

	cpu.Step()

	assert.Equal(t, 0x0303, cpu.PC)
	assert.Equal(t, 0x42, cpu.bus.Read(0x8000))
}

//// TAX

func TestTAX(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0xAA}, 0x0300)
	cpu.A = 0x42

	cpu.Step()

	assert.Equal(t, 0x0301, cpu.PC)
	assert.Equal(t, 0x42, cpu.X)
}

func TestTAXNegative(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0xAA}, 0x0300)
	cpu.A = 0xE0

	cpu.Step()

	assert.Equal(t, 0x0301, cpu.PC)
	assert.Equal(t, 0xE0, cpu.X)
	assert.True(t, cpu.getNegative())
}

func TestTAXZero(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0xAA}, 0x0300)
	cpu.A = 0x00

	cpu.Step()

	assert.Equal(t, 0x0301, cpu.PC)
	assert.Equal(t, 0x00, cpu.X)
	assert.True(t, cpu.getZero())
}

//// TAY

func TestTAY(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0xA8}, 0x0300)
	cpu.A = 0x42

	cpu.Step()

	assert.Equal(t, 0x0301, cpu.PC)
	assert.Equal(t, 0x42, cpu.Y)
}

func TestTAYNegative(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0xA8}, 0x0300)
	cpu.A = 0xE0

	cpu.Step()

	assert.Equal(t, 0x0301, cpu.PC)
	assert.Equal(t, 0xE0, cpu.Y)
	assert.True(t, cpu.getNegative())
}

func TestTAYZero(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0xA8}, 0x0300)
	cpu.A = 0x00

	cpu.Step()

	assert.Equal(t, 0x0301, cpu.PC)
	assert.Equal(t, 0x00, cpu.Y)
	assert.True(t, cpu.getZero())
}

//// TXA

func TestTXA(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x8A}, 0x0300)
	cpu.X = 0x42

	cpu.Step()

	assert.Equal(t, 0x0301, cpu.PC)
	assert.Equal(t, 0x42, cpu.A)
}

func TestTXANegative(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x8A}, 0x0300)
	cpu.X = 0xE0

	cpu.Step()

	assert.Equal(t, 0x0301, cpu.PC)
	assert.Equal(t, 0xE0, cpu.A)
	assert.True(t, cpu.getNegative())
}

func TestTXAZero(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x8A}, 0x0300)
	cpu.X = 0x00

	cpu.Step()

	assert.Equal(t, 0x0301, cpu.PC)
	assert.Equal(t, 0x00, cpu.A)
	assert.True(t, cpu.getZero())
}

//// TYA

func TestTYA(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x98}, 0x0300)
	cpu.Y = 0x42

	cpu.Step()

	assert.Equal(t, 0x0301, cpu.PC)
	assert.Equal(t, 0x42, cpu.A)
}

func TestTYANegative(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x98}, 0x0300)
	cpu.Y = 0xE0

	cpu.Step()

	assert.Equal(t, 0x0301, cpu.PC)
	assert.Equal(t, 0xE0, cpu.A)
	assert.True(t, cpu.getNegative())
}

func TestTYAZero(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x98}, 0x0300)
	cpu.Y = 0x00

	cpu.Step()

	assert.Equal(t, 0x0301, cpu.PC)
	assert.Equal(t, 0x00, cpu.A)
	assert.True(t, cpu.getZero())
}

//// TSX

func TestTSX(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0xBA}, 0x0300)
	cpu.SP = 0x42

	cpu.Step()

	assert.Equal(t, 0x0301, cpu.PC)
	assert.Equal(t, 0x42, cpu.X)
}

func TestTSXNegative(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0xBA}, 0x0300)
	cpu.SP = 0xE0

	cpu.Step()

	assert.Equal(t, 0x0301, cpu.PC)
	assert.Equal(t, 0xE0, cpu.X)
	assert.True(t, cpu.getNegative())
}

func TestTSXZero(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0xBA}, 0x0300)
	cpu.SP = 0x00

	cpu.Step()

	assert.Equal(t, 0x0301, cpu.PC)
	assert.Equal(t, 0x00, cpu.X)
	assert.True(t, cpu.getZero())
}

//// TXS

func TestTXS(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x9A}, 0x0300)
	cpu.X = 0x42

	cpu.Step()

	assert.Equal(t, 0x0301, cpu.PC)
	assert.Equal(t, 0x42, cpu.SP)
}

func TestTXSNegative(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x9A}, 0x0300)
	cpu.X = 0xE0

	cpu.Step()

	assert.Equal(t, 0x0301, cpu.PC)
	assert.Equal(t, 0xE0, cpu.SP)
	assert.True(t, cpu.getNegative())
}

func TestTXSZero(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x9A}, 0x0300)
	cpu.X = 0x00

	cpu.Step()

	assert.Equal(t, 0x0301, cpu.PC)
	assert.Equal(t, 0x00, cpu.SP)
	assert.True(t, cpu.getZero())
}

//// ASL

func TestASLaccumulator(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x0A}, 0x0300)
	cpu.A = 0x01

	cpu.Step()

	assert.Equal(t, 0x0301, cpu.PC)
	assert.Equal(t, 0x02, cpu.A)
}

func TestASLAccumulatorNegative(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x0A}, 0x0300)
	cpu.A = 0x40 // 0100 0000

	cpu.Step()

	assert.Equal(t, 0x0301, cpu.PC)
	assert.Equal(t, 0x80, cpu.A)
	assert.True(t, cpu.getNegative())
}

func TestASLAccumulatorZero(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x0A}, 0x0300)
	cpu.A = 0x80 // 1000 0000

	cpu.Step()

	assert.Equal(t, 0x0301, cpu.PC)
	assert.Equal(t, 0x00, cpu.A)
	assert.True(t, cpu.getZero())
}

func TestASLAccumulatorCarry(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x0A}, 0x0300)
	cpu.A = 0xAA // 1010 1010

	cpu.Step()

	assert.Equal(t, 0x0301, cpu.PC)
	assert.Equal(t, 0x54, cpu.A)
	assert.True(t, cpu.getCarry())
}

func TestASLzeropage(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x06, 0x80}, 0x0300)
	cpu.bus.Write(0x0080, 0x01)

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0x02, cpu.bus.Read(0x0080))
}

func TestASLzeropageNegative(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x06, 0x80}, 0x0300)
	cpu.bus.Write(0x0080, 0x40)

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0x80, cpu.bus.Read(0x0080))
	assert.True(t, cpu.getNegative())
}

func TestASLzeropageZero(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x06, 0x80}, 0x0300)
	cpu.bus.Write(0x0080, 0x80)

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0x00, cpu.bus.Read(0x0080))
	assert.True(t, cpu.getZero())
}

func TestASLzeropageCarry(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x06, 0x80}, 0x0300)
	cpu.bus.Write(0x0080, 0xAA)

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0x54, cpu.bus.Read(0x0080))
	assert.True(t, cpu.getCarry())
}

func TestASLzeropageX(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x16, 0x80}, 0x0300)
	cpu.X = 0x02
	cpu.bus.Write(0x0082, 0x01)

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0x02, cpu.bus.Read(0x0082))
}

func TestASLabsolute(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x0E, 0x00, 0x80}, 0x0300)
	cpu.bus.Write(0x8000, 0x01)

	cpu.Step()

	assert.Equal(t, 0x0303, cpu.PC)
	assert.Equal(t, 0x02, cpu.bus.Read(0x8000))
}

func TestASLabsoluteX(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x1E, 0x00, 0x80}, 0x0300)
	cpu.X = 0x02
	cpu.bus.Write(0x8002, 0x01)

	cpu.Step()

	assert.Equal(t, 0x0303, cpu.PC)
	assert.Equal(t, 0x02, cpu.bus.Read(0x8002))
}

//// LSR

func TestLSRaccumulator(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x4A}, 0x0300)
	cpu.A = 0x02

	cpu.Step()

	assert.Equal(t, 0x0301, cpu.PC)
	assert.Equal(t, 0x01, cpu.A)
}

func TestLSRAccumulatorZero(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x4A}, 0x0300)
	cpu.A = 0x01

	cpu.Step()

	assert.Equal(t, 0x0301, cpu.PC)
	assert.Equal(t, 0x00, cpu.A)
	assert.True(t, cpu.getZero())
}

func TestLSRAccumulatorCarry(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x4A}, 0x0300)
	cpu.A = 0x01

	cpu.Step()

	assert.Equal(t, 0x0301, cpu.PC)
	assert.Equal(t, 0x00, cpu.A)
	assert.True(t, cpu.getCarry())
}

func TestLSRzeropage(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x46, 0x80}, 0x0300)
	cpu.bus.Write(0x0080, 0x02)

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0x01, cpu.bus.Read(0x0080))
}

func TestLSRzeropageZero(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x46, 0x80}, 0x0300)
	cpu.bus.Write(0x0080, 0x01)

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0x00, cpu.bus.Read(0x0080))
	assert.True(t, cpu.getZero())
}

func TestLSRzeropageCarry(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x46, 0x80}, 0x0300)
	cpu.bus.Write(0x0080, 0x01)

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0x00, cpu.bus.Read(0x0080))
	assert.True(t, cpu.getCarry())
}

func TestLSRzeropageX(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x56, 0x80}, 0x0300)
	cpu.X = 0x02
	cpu.bus.Write(0x0082, 0x04)

	cpu.Step()

	assert.Equal(t, 0x0302, cpu.PC)
	assert.Equal(t, 0x02, cpu.bus.Read(0x0082))
}

func TestLSRabsolute(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x4E, 0x00, 0x80}, 0x0300)
	cpu.bus.Write(0x8000, 0x04)

	cpu.Step()

	assert.Equal(t, 0x0303, cpu.PC)
	assert.Equal(t, 0x02, cpu.bus.Read(0x8000))
}

func TestLSRabsoluteX(t *testing.T) {
	cpu, _, _ := NewRamMachine()
	cpu.LoadProgram([]byte{0x5E, 0x00, 0x80}, 0x0300)
	cpu.X = 0x02
	cpu.bus.Write(0x8002, 0x04)

	cpu.Step()

	assert.Equal(t, 0x0303, cpu.PC)
	assert.Equal(t, 0x02, cpu.bus.Read(0x8002))
}
