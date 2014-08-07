// Package bus represents the 16-bit address and 8-bit data bus for the
// 6502. Different memory modules can be used for different addresses.
package bus

import (
	"fmt"
	"github.com/ariejan/i6502/memory"
)

type busModule struct {
	memory memory.Memory
	name   string
	start  uint16
	end    uint16
}

// A 16-bit address and 8-bit data bus. It maps access to different
// attached modules, like Ram or Rom.
type Bus struct {
	modules []busModule
}

func (module busModule) String() string {
	return fmt.Sprintf("%s\t0x%04X-%04X\n", module.name, module.start, module.end)
}

func (bus *Bus) String() string {
	output := "\n\nAddress bus modules:\n\n"

	for _, module := range bus.modules {
		// output = append(output, module.String())
		output += module.String()
	}

	return output
}

// Creates a new Bus, no modules are attached by default.
func CreateBus() (*Bus, error) {
	return &Bus{modules: make([]busModule, 0)}, nil
}

func (bus *Bus) Attach(memory memory.Memory, name string, offset uint16) error {
	offsetMemory := OffsetMemory{Offset: offset, Memory: memory}
	end := offset + uint16(memory.Size()-1)

	module := busModule{memory: offsetMemory, name: name, start: offset, end: end}

	bus.modules = append(bus.modules, module)

	return nil
}

func (bus *Bus) backendFor(address uint16) (memory.Memory, error) {
	for _, module := range bus.modules {
		if address >= module.start && address <= module.end {
			return module.memory, nil
		}
	}

	return nil, fmt.Errorf("No module addressable at 0x%04X", address)
}

// Read an 8-bit value from the module mapped on the bus at the
// given 16-bit address.
func (bus *Bus) Read(address uint16) byte {
	memory, err := bus.backendFor(address)
	if err != nil {
		panic(err)
	}

	value := memory.Read(address)
	return value
}

// Write an 8-bit value to the module mapped on the bus through the
// 16-bit address.
func (bus *Bus) Write(address uint16, value byte) {
	memory, err := bus.backendFor(address)
	if err != nil {
		panic(err)
	}

	memory.Write(address, value)
}

// Helper method to read a 16-bit value from
// Note: LSB goes first: address [lo] + address+1 [hi]
func (bus *Bus) Read16(address uint16) uint16 {
	lo := uint16(bus.Read(address))
	hi := uint16(bus.Read(address + 1))
	return hi<<8 | lo
}

// Helper to write a 16-bit value to address and address + 1
// Note: LSB goes first: address [lo] + address+1 [hi]
func (bus *Bus) Write16(address uint16, value uint16) {
	bus.Write(address, byte(value))
	bus.Write(address+1, byte(value>>8))
}
