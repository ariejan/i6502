package memory

import (
	"io/ioutil"
)

// Rom provies 16kB of Read Only Memory, typcially loaded from file.
type Rom struct {
	data [0x4000]byte
}

func (rom *Rom) Size() int {
	return len(rom.data)
}

func (rom *Rom) Read(address uint16) byte {
	return rom.data[address]
}

func (rom *Rom) Write(address uint16, value byte) {
	panic("Cannot write to ROM!")
}

// Load ROM from a binary file. The data is placed
// at the beginning of ROM.
func LoadRomFromFile(path string) (*Rom, error) {
	// Read data from file
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Create the rom instance
	rom := &Rom{}

	// Load data into ROM
	for i, b := range data {
		rom.data[i] = b
	}

	return rom, nil
}
