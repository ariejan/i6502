package i6502

type Memory interface {
	Size() uint16
	Read(address uint16) byte
	Write(address uint16, data byte)
}
