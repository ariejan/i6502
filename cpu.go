package i6502

import "fmt"

type Cpu struct {
	bus *AddressBus // The address bus

	PC uint16 // 16-bit program counter
	P  byte   // Status Register
	SP byte   // Stack Pointer

	A byte // Accumulator
	X byte // X index register
	Y byte // Y index register
}

const (
	ResetVector = 0xFFFC // 0xFFFC-FFFD
)

// Create an new Cpu instance with the specified AddressBus
func NewCpu(bus *AddressBus) (*Cpu, error) {
	return &Cpu{bus: bus}, nil
}

func (c *Cpu) HasAddressBus() bool {
	return c.bus != nil
}

// Reset the CPU, emulating the RESB pin.
func (c *Cpu) Reset() {
	c.PC = c.bus.Read16(ResetVector)
	c.P = 0x34
}

// Load the specified program data at the given memory location
// and point the Program Counter to the beginning of the program
func (c *Cpu) LoadProgram(data []byte, location uint16) {
	for i, b := range data {
		c.bus.Write(location+uint16(i), b)
	}

	c.PC = location
}

// Execute the instruction pointed to by the Program Counter (PC)
func (c *Cpu) Step() {
	instruction := c.readNextInstruction()
	c.PC += uint16(instruction.Size)
	// fmt.Println(instruction)
	c.execute(instruction)
}

func (c *Cpu) execute(instruction Instruction) {
	switch instruction.opcodeId {
	case nop:
		break
	case adc:
		c.ADC(instruction)
	case sbc:
		c.SBC(instruction)
	case sec:
		c.setCarry(true)
	case sed:
		c.setDecimal(true)
	case sei:
		c.setInterrupt(true)
	case clc:
		c.setCarry(false)
	case cld:
		c.setDecimal(false)
	case cli:
		c.setInterrupt(false)
	case clv:
		c.setOverflow(false)
	case inx:
		c.setX(c.X + 1)
	case iny:
		c.setY(c.Y + 1)
	case inc:
		c.INC(instruction)
	case dex:
		c.setX(c.X - 1)
	case dey:
		c.setY(c.Y - 1)
	case dec:
		c.DEC(instruction)
	case lda:
		c.LDA(instruction)
	case ldx:
		c.LDX(instruction)
	case ldy:
		c.LDY(instruction)
	case ora:
		c.ORA(instruction)
	case and:
		c.AND(instruction)
	case eor:
		c.EOR(instruction)
	case sta:
		c.STA(instruction)
	case stx:
		c.STX(instruction)
	case sty:
		c.STY(instruction)
	case tax:
		c.setX(c.A)
	case tay:
		c.setY(c.A)
	case txa:
		c.setA(c.X)
	case tya:
		c.setA(c.Y)
	case tsx:
		c.setX(c.SP)
	case txs:
		c.setSP(c.X)
	case asl:
		c.ASL(instruction)
	case lsr:
		c.LSR(instruction)
	default:
		panic(fmt.Errorf("Unimplemented instruction: %s", instruction))
	}
}

func (c *Cpu) resolveOperand(in Instruction) uint8 {
	switch in.addressingId {
	case immediate:
		return in.Op8
	default:
		return c.bus.Read(c.memoryAddress(in))
	}
}

func (c *Cpu) memoryAddress(in Instruction) uint16 {
	switch in.addressingId {
	case absolute:
		return in.Op16
	case absoluteX:
		return in.Op16 + uint16(c.X)
	case absoluteY:
		return in.Op16 + uint16(c.Y)
	case indirect:
		return c.bus.Read16(in.Op16)
	case indirectX:
		return c.bus.Read16(uint16(in.Op8 + c.X))
	case indirectY:
		return c.bus.Read16(uint16(in.Op8)) + uint16(c.Y)
	case relative:
		panic("Relative addressing not yet implemented.")
	case zeropage:
		return uint16(in.Op8)
	case zeropageX:
		return uint16(in.Op8 + c.X)
	case zeropageY:
		return uint16(in.Op8 + c.Y)
	default:
		panic(fmt.Errorf("Unhandled addressing mode. Are you sure you are running a 6502 ROM?"))
	}
}
