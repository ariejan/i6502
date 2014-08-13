package i6502

type Cpu struct {
	bus *AddressBus // The address bus

	PC uint16 // 16-bit program counter
	P  byte   // Status Register
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
