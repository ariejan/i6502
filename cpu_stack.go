package i6502

const (
	StackBase = 0x0100
)

func (c *Cpu) stackPush(data byte) {
	c.bus.Write(StackBase+uint16(c.SP), data)
	c.SP -= 1
}

func (c *Cpu) stackPeek() byte {
	return c.bus.Read(StackBase + uint16(c.SP+1))
}

func (c *Cpu) stackPop() byte {
	c.SP += 1
	return c.bus.Read(StackBase + uint16(c.SP))
}
