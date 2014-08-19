package i6502

/*
Anything implementing the Memory interface can be attached to the AddressBus
and become accessible by the Cpu.
*/
type Memory interface {
	Size() uint16
	ReadByte(address uint16) byte
	WriteByte(address uint16, data byte)
}
