package i6502

import "fmt"

type Cpu struct {
	PC uint16 // 16-bit program counter
	P  byte   // Status Register
	SP byte   // Stack Pointer

	A byte // Accumulator
	X byte // X index register
	Y byte // Y index register

	Bus *AddressBus // The address bus
}

const (
	ResetVector = 0xFFFC // 0xFFFC-FFFD
	IrqVector   = 0xFFFE // 0xFFFE-FFFF

	StackBase = 0x0100 // One page 0x0100-01FF
)

// Create an new Cpu instance with the specified AddressBus
func NewCpu(bus *AddressBus) (*Cpu, error) {
	return &Cpu{Bus: bus}, nil
}

func (c *Cpu) String() string {
	str := ">>> CPU [  A ] [  X ] [  Y ] [ SP ] [  PC  ] NVxBDIZC\n>>>      0x%02X   0x%02X   0x%02X   0x%02X   0x%04X  %08b\n"
	return fmt.Sprintf(str, c.A, c.X, c.Y, c.SP, c.PC, c.P)
}

func (c *Cpu) hasAddressBus() bool {
	return c.Bus != nil
}

// Reset the CPU, emulating the RESB pin.
func (c *Cpu) Reset() {
	c.PC = c.Bus.Read16(ResetVector)
	c.P = 0x34

	// Not specified, but let's clean up
	c.A = 0x00
	c.X = 0x00
	c.Y = 0x00
	c.SP = 0xFF
}

// Simulate the IRQ pin
func (c *Cpu) Interrupt() {
	c.handleIrq(c.PC)
}

func (c *Cpu) handleIrq(PC uint16) {
	c.stackPush(byte(PC >> 8))
	c.stackPush(byte(PC))
	c.stackPush(c.P)

	c.setIrqDisable(true)

	c.PC = c.Bus.Read16(IrqVector)
}

// Load the specified program data at the given memory location
// and point the Program Counter to the beginning of the program
func (c *Cpu) LoadProgram(data []byte, location uint16) {
	for i, b := range data {
		c.Bus.Write(location+uint16(i), b)
	}

	c.PC = location
}

// Execute the instruction pointed to by the Program Counter (PC)
func (c *Cpu) Step() {
	instruction := c.readNextInstruction()
	c.PC += uint16(instruction.Size)
	c.execute(instruction)
}

func (c *Cpu) execute(instruction Instruction) {
	switch instruction.opcodeId {
	case nop:
		break
	case adc:
		c.adc(instruction)
	case sbc:
		c.sbc(instruction)
	case sec:
		c.setCarry(true)
	case sed:
		c.setDecimal(true)
	case sei:
		c.setIrqDisable(true)
	case clc:
		c.setCarry(false)
	case cld:
		c.setDecimal(false)
	case cli:
		c.setIrqDisable(false)
	case clv:
		c.setOverflow(false)
	case inx:
		c.setX(c.X + 1)
	case iny:
		c.setY(c.Y + 1)
	case inc:
		c.inc(instruction)
	case dex:
		c.setX(c.X - 1)
	case dey:
		c.setY(c.Y - 1)
	case dec:
		c.dec(instruction)
	case lda:
		value := c.resolveOperand(instruction)
		c.setA(value)
	case ldx:
		value := c.resolveOperand(instruction)
		c.setX(value)
	case ldy:
		value := c.resolveOperand(instruction)
		c.setY(value)
	case ora:
		value := c.resolveOperand(instruction)
		c.setA(c.A | value)
	case and:
		value := c.resolveOperand(instruction)
		c.setA(c.A & value)
	case eor:
		value := c.resolveOperand(instruction)
		c.setA(c.A ^ value)
	case sta:
		address := c.memoryAddress(instruction)
		c.Bus.Write(address, c.A)
	case stx:
		address := c.memoryAddress(instruction)
		c.Bus.Write(address, c.X)
	case sty:
		address := c.memoryAddress(instruction)
		c.Bus.Write(address, c.Y)
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
		c.SP = c.X
	case asl:
		c.asl(instruction)
	case lsr:
		c.lsr(instruction)
	case rol:
		c.rol(instruction)
	case ror:
		c.ror(instruction)
	case cmp:
		value := c.resolveOperand(instruction)
		c.setCarry(c.A >= value)
		c.setArithmeticFlags(c.A - value)
	case cpx:
		value := c.resolveOperand(instruction)
		c.setCarry(c.X >= value)
		c.setArithmeticFlags(c.X - value)
	case cpy:
		value := c.resolveOperand(instruction)
		c.setCarry(c.Y >= value)
		c.setArithmeticFlags(c.Y - value)
	case brk:
		c.setBreak(true)
		c.handleIrq(c.PC + 1)
	case bcc:
		if !c.getCarry() {
			c.branch(instruction)
		}
	case bcs:
		if c.getCarry() {
			c.branch(instruction)
		}
	case bne:
		if !c.getZero() {
			c.branch(instruction)
		}
	case beq:
		if c.getZero() {
			c.branch(instruction)
		}
	case bpl:
		if !c.getNegative() {
			c.branch(instruction)
		}
	case bmi:
		if c.getNegative() {
			c.branch(instruction)
		}
	case bvc:
		if !c.getOverflow() {
			c.branch(instruction)
		}
	case bvs:
		if c.getOverflow() {
			c.branch(instruction)
		}
	case bit:
		value := c.resolveOperand(instruction)
		c.setNegative((value & 0x80) != 0)
		c.setOverflow((value & 0x40) != 0)
		c.setZero((c.A & value) == 0)
	case php:
		c.stackPush(c.P | 0x30)
	case plp:
		c.setP(c.stackPop())
	case pha:
		c.stackPush(c.A)
	case pla:
		value := c.stackPop()
		c.setA(value)
	case jmp:
		c.PC = c.memoryAddress(instruction)
	case jsr:
		c.stackPush(byte((c.PC - 1) >> 8))
		c.stackPush(byte(c.PC - 1))
		c.PC = c.memoryAddress(instruction)
	case rts:
		c.PC = (uint16(c.stackPop()) | uint16(c.stackPop())<<8) + 1
	case rti:
		c.setP(c.stackPop())
		c.PC = uint16(c.stackPop()) | uint16(c.stackPop())<<8
	default:
		panic(fmt.Errorf("Unimplemented instruction: %s", instruction))
	}
}

func (c *Cpu) readNextInstruction() Instruction {
	// Read the opcode
	opcode := c.Bus.Read(c.PC)

	optype, ok := opTypes[opcode]
	if !ok {
		panic(fmt.Sprintf("Unknown or unimplemented opcode 0x%02X\n%s", opcode, c.String()))
	}

	instruction := Instruction{OpType: optype, Address: c.PC}
	switch instruction.Size {
	case 1: // Zero operand instruction
	case 2: // 8-bit operand
		instruction.Op8 = c.Bus.Read(c.PC + 1)
	case 3: // 16-bit operand
		instruction.Op16 = c.Bus.Read16(c.PC + 1)
	}

	return instruction
}

func (c *Cpu) branch(in Instruction) {
	relative := int8(in.Op8) // Signed!
	if relative >= 0 {
		c.PC += uint16(relative)
	} else {
		c.PC -= -uint16(relative)
	}
}

func (c *Cpu) resolveOperand(in Instruction) uint8 {
	switch in.addressingId {
	case immediate:
		return in.Op8
	default:
		return c.Bus.Read(c.memoryAddress(in))
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
		return c.Bus.Read16(in.Op16)
	case indirectX:
		return c.Bus.Read16(uint16(in.Op8 + c.X))
	case indirectY:
		return c.Bus.Read16(uint16(in.Op8)) + uint16(c.Y)
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

// Add Memory to Accumulator with Carry
func (c *Cpu) adc(in Instruction) {
	operand := c.resolveOperand(in)
	carryIn := c.getCarryInt()

	if c.getDecimal() {
		c.adcDecimal(c.A, operand, carryIn)
	} else {
		c.adcNormal(c.A, operand, carryIn)
	}
}

// Substract memory from Accummulator with carry
func (c *Cpu) sbc(in Instruction) {
	operand := c.resolveOperand(in)
	carryIn := c.getCarryInt()

	// fmt.Printf("SBC: A: 0x%02X V: 0x%02X C: %b D: %v\n", c.A, operand, carryIn, c.getDecimal())

	if c.getDecimal() {
		c.sbcDecimal(c.A, operand, carryIn)
	} else {
		c.adcNormal(c.A, ^operand, carryIn)
	}
}

func (c *Cpu) inc(in Instruction) {
	address := c.memoryAddress(in)
	value := c.Bus.Read(address) + 1

	c.Bus.Write(address, value)
	c.setArithmeticFlags(value)
}

func (c *Cpu) dec(in Instruction) {
	address := c.memoryAddress(in)
	value := c.Bus.Read(address) - 1

	c.Bus.Write(address, value)
	c.setArithmeticFlags(value)
}

func (c *Cpu) asl(in Instruction) {
	switch in.addressingId {
	case accumulator:
		c.setCarry((c.A >> 7) == 1)
		c.A <<= 1
		c.setArithmeticFlags(c.A)
	default:
		address := c.memoryAddress(in)
		value := c.Bus.Read(address)
		c.setCarry((value >> 7) == 1)
		value <<= 1
		c.Bus.Write(address, value)
		c.setArithmeticFlags(value)
	}
}

func (c *Cpu) lsr(in Instruction) {
	switch in.addressingId {
	case accumulator:
		c.setCarry((c.A & 0x01) == 1)
		c.A >>= 1
		c.setArithmeticFlags(c.A)
	default:
		address := c.memoryAddress(in)
		value := c.Bus.Read(address)
		c.setCarry((value & 0x01) == 1)
		value >>= 1
		c.Bus.Write(address, value)
		c.setArithmeticFlags(value)
	}
}

func (c *Cpu) rol(in Instruction) {
	carry := c.getCarryInt()

	switch in.addressingId {
	case accumulator:
		c.setCarry((c.A & 0x80) != 0)
		c.A = c.A<<1 | carry
		c.setArithmeticFlags(c.A)
	default:
		address := c.memoryAddress(in)
		value := c.Bus.Read(address)
		c.setCarry((value & 0x80) != 0)
		value = value<<1 | carry
		c.Bus.Write(address, value)
		c.setArithmeticFlags(value)
	}
}

func (c *Cpu) ror(in Instruction) {
	carry := c.getCarryInt()

	switch in.addressingId {
	case accumulator:
		c.setCarry(c.A&0x01 == 1)
		c.A = c.A>>1 | carry<<7
		c.setArithmeticFlags(c.A)
	default:
		address := c.memoryAddress(in)
		value := c.Bus.Read(address)
		c.setCarry(value&0x01 == 1)
		value = value>>1 | carry<<7
		c.Bus.Write(address, value)
		c.setArithmeticFlags(value)
	}
}

// Performs regular, 8-bit addition
func (c *Cpu) adcNormal(a uint8, b uint8, carryIn uint8) {
	result16 := uint16(a) + uint16(b) + uint16(carryIn)
	result := uint8(result16)
	carryOut := (result16 & 0x100) != 0
	overflow := (a^result)&(b^result)&0x80 != 0

	// Set the carry flag if we exceed 8-bits
	c.setCarry(carryOut)
	// Set the overflow bit
	c.setOverflow(overflow)
	// Store the resulting value (8-bits)
	c.setA(result)
}

// Performs addition in decimal mode
func (c *Cpu) adcDecimal(a uint8, b uint8, carryIn uint8) {
	var carryB uint8 = 0

	low := (a & 0x0F) + (b & 0x0F) + carryIn
	if (low & 0xFF) > 9 {
		low += 6
	}
	if low > 15 {
		carryB = 1
	}

	high := (a >> 4) + (b >> 4) + carryB
	if (high & 0xFF) > 9 {
		high += 6
	}

	result := (low & 0x0F) | (high<<4)&0xF0

	c.setCarry(high > 15)
	c.setZero(result == 0)
	c.setNegative(false) // BCD never sets negative
	c.setOverflow(false) // BCD never sets overflow

	c.A = result
}

func (c *Cpu) sbcDecimal(a uint8, b uint8, carryIn uint8) {
	var carryB uint8 = 0

	if carryIn == 0 {
		carryIn = 1
	} else {
		carryIn = 0
	}

	low := (a & 0x0F) - (b & 0x0F) - carryIn
	if (low & 0x10) != 0 {
		low -= 6
	}
	if (low & 0x10) != 0 {
		carryB = 1
	}

	high := (a >> 4) - (b >> 4) - carryB
	if (high & 0x10) != 0 {
		high -= 6
	}

	result := (low & 0x0F) | (high << 4)

	c.setCarry((high & 0xFF) < 15)
	c.setZero(result == 0)
	c.setNegative(false) // BCD never sets negative
	c.setOverflow(false) // BCD never sets overflow

	c.A = result
}

func (c *Cpu) stackPush(data byte) {
	c.Bus.Write(StackBase+uint16(c.SP), data)
	c.SP -= 1
}

func (c *Cpu) stackPeek() byte {
	return c.Bus.Read(StackBase + uint16(c.SP+1))
}

func (c *Cpu) stackPop() byte {
	c.SP += 1
	return c.Bus.Read(StackBase + uint16(c.SP))
}
