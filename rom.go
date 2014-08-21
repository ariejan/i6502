package i6502

import (
	"fmt"
	"io/ioutil"
)

/*
Read-Only Memory
*/
type Rom struct {
	data []byte
}

/*
Create a new Rom component, using the content of `path`. The file automatically
specifies the size of Rom.
*/
func NewRom(path string) (*Rom, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return &Rom{data: data}, nil
}

func (r *Rom) Size() uint16 {
	return uint16(len(r.data))
}

func (r *Rom) ReadByte(address uint16) byte {
	return r.data[address]
}

func (r *Rom) WriteByte(address uint16, data byte) {
	panic(fmt.Errorf("Trying to write to ROM at 0x%04X", address))
}
