package i6502

// Status Register bits
const (
	sCarry = iota
	sZero
	sIrqDisable
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

func (c *Cpu) setZero(state bool) {
	c.setStatus(sZero, state)
}

func (c *Cpu) setIrqDisable(state bool) {
	c.setStatus(sIrqDisable, state)
}

func (c *Cpu) setDecimal(state bool) {
	c.setStatus(sDecimal, state)
}

func (c *Cpu) setBreak(state bool) {
	c.setStatus(sBreak, state)
}

func (c *Cpu) setOverflow(state bool) {
	c.setStatus(sOverflow, state)
}

func (c *Cpu) setNegative(state bool) {
	c.setStatus(sNegative, state)
}

func (c *Cpu) getCarry() bool {
	return c.getStatus(sCarry)
}

func (c *Cpu) getCarryInt() uint8 {
	return c.getStatusInt(sCarry)
}

func (c *Cpu) getZero() bool {
	return c.getStatus(sZero)
}

func (c *Cpu) getIrqDisable() bool {
	return c.getStatus(sIrqDisable)
}

func (c *Cpu) getDecimal() bool {
	return c.getStatus(sDecimal)
}

func (c *Cpu) getBreak() bool {
	return c.getStatus(sBreak)
}

func (c *Cpu) getOverflow() bool {
	return c.getStatus(sOverflow)
}

func (c *Cpu) getNegative() bool {
	return c.getStatus(sNegative)
}

func (c *Cpu) setArithmeticFlags(value uint8) {
	// Set sZero if value is 0
	c.setStatus(sZero, value == 0)

	// Set sNegative if the 8th bit is 1, we're dealing with
	// uint8's internally, and using two's complement to identify
	// negatives
	c.setStatus(sNegative, (value>>7) == 1)
}

func (c *Cpu) setA(value byte) {
	c.A = value
	c.setArithmeticFlags(c.A)
}

func (c *Cpu) setX(value byte) {
	c.X = value
	c.setArithmeticFlags(c.X)
}

func (c *Cpu) setY(value byte) {
	c.Y = value
	c.setArithmeticFlags(c.Y)
}

func (c *Cpu) setSP(value byte) {
	c.SP = value
	c.setArithmeticFlags(c.SP)
}

func (c *Cpu) setP(value byte) {
	c.P = value
	c.P |= 0x20
}
