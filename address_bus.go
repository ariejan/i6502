package i6502

import (
	"fmt"
)

type AddressBus struct {
	addressables []addressable // Different components
}

type addressable struct {
	memory Memory // Actual memory
	start  uint16 // First address in address space
	end    uint16 // Last address in address space
}

func NewAddressBus() (*AddressBus, error) {
	return &AddressBus{addressables: make([]addressable, 0)}, nil
}

func (a *addressable) String() string {
	return fmt.Sprintf("\t0x%04X-%04X\n", a.start, a.end)
}

func (a *AddressBus) String() string {
	output := "Address Bus:\n"

	for _, addressable := range a.addressables {
		output += addressable.String()
	}

	return output
}

func (a *AddressBus) AddressablesCount() int {
	return len(a.addressables)
}

func (a *AddressBus) Attach(memory Memory, offset uint16) {
	start := offset
	end := offset + memory.Size() - 1
	addressable := addressable{memory: memory, start: start, end: end}

	a.addressables = append(a.addressables, addressable)
}

func (a *AddressBus) addressableForAddress(address uint16) (*addressable, error) {
	for _, addressable := range a.addressables {
		if addressable.start <= address && addressable.end >= address {
			return &addressable, nil
		}
	}

	return nil, fmt.Errorf("No addressable memory found at 0x%04X", address)
}

func (a *AddressBus) Read(address uint16) byte {
	addressable, err := a.addressableForAddress(address)
	if err != nil {
		panic(err)
	}

	return addressable.memory.Read(address - addressable.start)
}

func (a *AddressBus) Read16(address uint16) uint16 {
	lo := uint16(a.Read(address))
	hi := uint16(a.Read(address + 1))

	return (hi << 8) | lo
}

func (a *AddressBus) Write(address uint16, data byte) {
	addressable, err := a.addressableForAddress(address)
	if err != nil {
		panic(err)
	}

	addressable.memory.Write(address-addressable.start, data)
}

func (a *AddressBus) Write16(address uint16, data uint16) {
	a.Write(address, byte(data))
	a.Write(address+1, byte(data>>8))
}
