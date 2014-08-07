package memory

import (
	"io/ioutil"
)

// Ram provides 32kB of Random Access Memory
type Ram struct {
	data [0x8000]byte
}

func (ram *Ram) Size() int {
	return len(ram.data)
}

func (ram *Ram) Read(address uint16) byte {
	return ram.data[address]
}

func (ram *Ram) Write(address uint16, value byte) {
	ram.data[address] = value
}

func CreateRam() *Ram {
	return &Ram{}
}

func (ram *Ram) Dump(path string) {
	err := ioutil.WriteFile(path, ram.data[:], 0640)
	if err != nil {
		panic(err)
	}
}
