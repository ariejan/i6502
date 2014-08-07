// Package memory provides different memory components
// and a common interface.
package memory

type Memory interface {
	Size() int
	Read(address uint16) byte
	Write(address uint16, value byte)
}
