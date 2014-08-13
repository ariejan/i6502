package i6502

type Cpu struct {
	bus *AddressBus // The address bus

	PC uint16 // 16-bit program counter
	P  byte   // Status Register
}

const (
	ResetVector = 0xFFFC
)

func NewCpu(bus *AddressBus) (*Cpu, error) {
	return &Cpu{bus: bus}, nil
}

func (c *Cpu) HasAddressBus() bool {
	return c.bus != nil
}

func (c *Cpu) Reset() {
	c.PC = c.bus.Read16(ResetVector)
	c.P = 0x34
}
