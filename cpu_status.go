package i6502

// Status Register bits
const (
	sCarry = iota
	sZero
	sInterrupt
	sDecimal
	sBreak
	_
	sOverflow
	sNegative
)

func (c *Cpu) getStatusInt(bit uint8) uint8 {
	return (c.P >> bit) & 1
}

func (c *Cpu) getStatus(bit uint8) bool {
	return c.getStatusInt(bit) == 1
}

func (c *Cpu) setStatus(bit uint8, state bool) {
	if state {
		c.P |= 1 << bit
	} else {
		c.P &^= 1 << bit
	}
}

func (c *Cpu) setCarry(state bool) {
	c.setStatus(sCarry, state)
}

func (c *Cpu) setOverflow(state bool) {
	c.setStatus(sOverflow, state)
}

func (c *Cpu) setZero(state bool) {
	c.setStatus(sZero, state)
}

func (c *Cpu) setNegative(state bool) {
	c.setStatus(sNegative, state)
}

func (c *Cpu) setArithmeticFlags(value uint8) {
	// Set sZero if value is 0
	c.setStatus(sZero, value == 0)

	// Set sNegative if the 8th bit is 1, we're dealing with
	// uint8's internally, and using two's complement to identify
	// negatives
	c.setStatus(sNegative, (value>>7) == 1)
}
