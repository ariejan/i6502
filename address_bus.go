package i6502

import (
	"fmt"
)

/*
The AddressBus contains a list of all attached memory components,
like Ram, Rom and IO. It takes care of mapping the global 16-bit
address space of the Cpu to the relative memory addressing of
each component.
*/
type AddressBus struct {
	addressables []addressable // Different components
}

type addressable struct {
	memory Memory // Actual memory
	start  uint16 // First address in address space
	end    uint16 // Last address in address space
}

func (a *addressable) String() string {
	return fmt.Sprintf("\t0x%04X-%04X\n", a.start, a.end)
}

// Creates a new, empty 16-bit AddressBus
func NewAddressBus() (*AddressBus, error) {
	return &AddressBus{addressables: make([]addressable, 0)}, nil
}

// Returns a string with details about the AddressBus and attached memory
func (a *AddressBus) String() string {
	output := "Address Bus:\n"

	for _, addressable := range a.addressables {
		output += addressable.String()
	}

	return output
}

/*
Attach the given Memory at the specified memory offset.

To attach 16kB ROM at 0xC000-FFFF, you simple attach the Rom at
address 0xC000, the Size of the Memory determines the end-address.

    rom, _ := i6502.NewRom(0x4000)
    bus.Attach(rom, 0xC000)
*/
func (a *AddressBus) Attach(memory Memory, offset uint16) {
	start := offset
	end := offset + memory.Size() - 1
	addressable := addressable{memory: memory, start: start, end: end}

	a.addressables = append(a.addressables, addressable)
}

/*
Read an 8-bit value from Memory attached at the 16-bit address.

This will panic if you try to read from an address that has no Memory attached.
*/
func (a *AddressBus) ReadByte(address uint16) byte {
	addressable, err := a.addressableForAddress(address)
	if err != nil {
		panic(err)
	}

	return addressable.memory.ReadByte(address - addressable.start)
}

/*
Convenience method to quickly read a 16-bit value from address and address + 1.

Note that we first read the LOW byte from address and then the HIGH byte from address + 1.
*/
func (a *AddressBus) Read16(address uint16) uint16 {
	lo := uint16(a.ReadByte(address))
	hi := uint16(a.ReadByte(address + 1))

	return (hi << 8) | lo
}

/*
Write an 8-bit value to the Memory at the 16-bit address.

This will panic if you try to write to an address that has no Memory attached or
Memory that is read-only, like Rom.
*/
func (a *AddressBus) WriteByte(address uint16, data byte) {
	addressable, err := a.addressableForAddress(address)
	if err != nil {
		panic(err)
	}

	addressable.memory.WriteByte(address-addressable.start, data)
}

/*
Convenience method to quickly write a 16-bit value to address and address + 1.

Note that the LOW byte will be stored in address and the high byte in address + 1.
*/
func (a *AddressBus) Write16(address uint16, data uint16) {
	a.WriteByte(address, byte(data))
	a.WriteByte(address+1, byte(data>>8))
}

// Returns the addressable for the specified address, or an error if no addressable exists.
func (a *AddressBus) addressableForAddress(address uint16) (*addressable, error) {
	for _, addressable := range a.addressables {
		if addressable.start <= address && addressable.end >= address {
			return &addressable, nil
		}
	}

	return nil, fmt.Errorf("No addressable memory found at 0x%04X", address)
}
