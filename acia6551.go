package i6502

const (
	aciaData = iota
	aciaStatus
	aciaCommand
	aciaControl
)

/*
ACIA 6551 Serial IO

This Asynchronous Communications Interface Adapater can be
directly attached to the 6502's address and data busses.

It provides serial IO.

The supplied Rx and Tx channels can be used to read and wirte
data to the ACIA 6551.
*/
type Acia6551 struct {
	rx byte
	tx byte

	commandData byte
	controlData byte

	rxFull  bool
	txEmpty bool

	rxIrqEnabled bool
	txIrqEnabled bool

	overrun bool
}

func NewAcia6551() (*Acia6551, error) {
	acia := &Acia6551{}
	acia.Reset()

	return acia, nil
}

func (a *Acia6551) Size() uint16 {
	// We have a only 4 addresses, Data, Status, Command and Control
	return 0x04
}

// Emulates a hardware reset
func (a *Acia6551) Reset() {
	a.rx = 0
	a.rxFull = false

	a.tx = 0
	a.txEmpty = true

	a.rxIrqEnabled = false
	a.txIrqEnabled = false

	a.overrun = false

	a.setControl(0)
	a.setCommand(0)
}

func (a *Acia6551) setControl(data byte) {
	a.controlData = data
}

func (a *Acia6551) setCommand(data byte) {
	a.commandData = data

	a.rxIrqEnabled = (data & 0x02) != 0
	a.txIrqEnabled = ((data & 0x04) != 0) && ((data & 0x08) != 1)
}

func (a *Acia6551) statusRegister() byte {
	status := byte(0)

	if a.rxFull {
		status |= 0x08
	}

	if a.txEmpty {
		status |= 0x10
	}

	if a.overrun {
		status |= 0x04
	}

	return status
}

// Implements io.Reader, for external programs to read TX'ed data from
// the serial output.
func (a *Acia6551) Read(p []byte) (n int, err error) {
	a.txEmpty = true
	copy(p, []byte{a.tx})
	// TODO: Handle txInterrupt
	return 1, nil
}

// Implements io.Writer, for external programs to write to the
// ACIA's RX
func (a *Acia6551) Write(p []byte) (n int, err error) {
	for _, b := range p {
		a.rxWrite(b)
	}

	return len(p), nil
}

// Used by the AddressBus to read data from the ACIA 6551
func (a *Acia6551) ReadByte(address uint16) byte {
	switch address {
	case aciaData:
		return a.rxRead()
	case aciaStatus:
		return a.statusRegister()
	case aciaCommand:
		return a.commandData
	case aciaControl:
		return a.controlData
	}

	return 0x00
}

// Used by the AddressBus to write data to the ACIA 6551
func (a *Acia6551) WriteByte(address uint16, data byte) {
	switch address {
	case aciaData:
		a.txWrite(data)
	case aciaStatus:
		a.Reset()
	case aciaCommand:
		a.setCommand(data)
	case aciaControl:
		a.setControl(data)
	}
}

func (a *Acia6551) rxRead() byte {
	a.overrun = false
	a.rxFull = false
	return a.rx
}

func (a *Acia6551) rxWrite(data byte) {
	// Oh no, overrun. Set the appropriate status
	if a.rxFull {
		a.overrun = true
	}

	a.rx = data
	a.rxFull = true

	// TODO: Interrupts
}

func (a *Acia6551) txWrite(data byte) {
	a.tx = data
	a.txEmpty = false
}
