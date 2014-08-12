package devices

type Via6522 struct {
}

func (v *Via6522) Size() int {
	return 0x10
}

func (v *Via6522) Read(address uint16) byte {
	return 0x00
}

func (v *Via6522) Write(address uint16, data byte) {
	// NOP
}

func NewVia6522(interruptChan chan bool) *Via6522 {
	via := &Via6522{}
	return via
}
