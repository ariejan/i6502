package i6502

// Add Memory to Accumulator with Carry
func (c *Cpu) ADC(in Instruction) {
	operand := c.resolveOperand(in)
	carryIn := c.getCarryInt()

	if c.getDecimal() {
		c.adcDecimal(c.A, operand, carryIn)
	} else {
		c.adcNormal(c.A, operand, carryIn)
	}
}

// Substract memory from Accummulator with carry
func (c *Cpu) SBC(in Instruction) {
	operand := c.resolveOperand(in)
	carryIn := c.getCarryInt()

	if c.getDecimal() {
		c.sbcDecimal(c.A, operand, carryIn)
	} else {
		c.adcNormal(c.A, ^operand, carryIn)
	}
}

func (c *Cpu) INC(in Instruction) {
	address := c.memoryAddress(in)
	value := c.bus.Read(address) + 1

	c.bus.Write(address, value)
	c.setArithmeticFlags(value)
}

func (c *Cpu) DEC(in Instruction) {
	address := c.memoryAddress(in)
	value := c.bus.Read(address) - 1

	c.bus.Write(address, value)
	c.setArithmeticFlags(value)
}

func (c *Cpu) LDA(in Instruction) {
	value := c.resolveOperand(in)
	c.setA(value)
}

func (c *Cpu) LDX(in Instruction) {
	value := c.resolveOperand(in)
	c.setX(value)
}

func (c *Cpu) LDY(in Instruction) {
	value := c.resolveOperand(in)
	c.setY(value)
}

func (c *Cpu) ORA(in Instruction) {
	value := c.resolveOperand(in)
	c.setA(c.A | value)
}

func (c *Cpu) AND(in Instruction) {
	value := c.resolveOperand(in)
	c.setA(c.A & value)
}

func (c *Cpu) EOR(in Instruction) {
	value := c.resolveOperand(in)
	c.setA(c.A ^ value)
}

func (c *Cpu) STA(in Instruction) {
	address := c.memoryAddress(in)
	c.bus.Write(address, c.A)
}

func (c *Cpu) STX(in Instruction) {
	address := c.memoryAddress(in)
	c.bus.Write(address, c.X)
}

func (c *Cpu) STY(in Instruction) {
	address := c.memoryAddress(in)
	c.bus.Write(address, c.Y)
}

func (c *Cpu) ASL(in Instruction) {
	switch in.addressingId {
	case accumulator:
		c.setCarry((c.A >> 7) == 1)
		c.A <<= 1
		c.setArithmeticFlags(c.A)
	default:
		address := c.memoryAddress(in)
		value := c.bus.Read(address)
		c.setCarry((value >> 7) == 1)
		value <<= 1
		c.bus.Write(address, value)
		c.setArithmeticFlags(value)
	}
}

func (c *Cpu) LSR(in Instruction) {
	switch in.addressingId {
	case accumulator:
		c.setCarry((c.A & 0x01) == 1)
		c.A >>= 1
		c.setArithmeticFlags(c.A)
	default:
		address := c.memoryAddress(in)
		value := c.bus.Read(address)
		c.setCarry((value & 0x01) == 1)
		value >>= 1
		c.bus.Write(address, value)
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

	carryIn = (carryIn + 1) % 1

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
