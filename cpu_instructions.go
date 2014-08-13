package i6502

// Add Memory to Accumulator with Carry
func (c *Cpu) ADC(in Instruction) {
	operand := c.resolveOperand(in)
	carryIn := c.getStatusInt(sCarry)

	if c.getStatus(sDecimal) {
		c.adcDecimal(c.A, operand, carryIn)
	} else {
		c.adcNormal(c.A, operand, carryIn)
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
	c.A = result
	// Update sZero and sNegative
	c.setArithmeticFlags(c.A)
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
	c.setNegative(false) // BCD is never negative
	c.setOverflow(false) // BCD never sets overflow

	c.A = result
}
