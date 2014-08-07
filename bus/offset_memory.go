package bus

import (
	"github.com/ariejan/i6502/memory"
)

// The AddressDecoder routes a full 16-bit address to the
// appropriate relatie Memory component address.
type OffsetMemory struct {
	Offset uint16
	memory.Memory
}

func (offsetMemory OffsetMemory) Read(address uint16) byte {
	return offsetMemory.Memory.Read(address - offsetMemory.Offset)
}

func (offsetMemory OffsetMemory) Write(address uint16, value byte) {
	offsetMemory.Memory.Write(address-offsetMemory.Offset, value)
}
