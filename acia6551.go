package i6502

const (
	aciaData = iota
	aciaStatus
	aciaCommand
	aciaControl
)

/*
ACIA 6551 Serial IO
*/
type Acia6551 struct {
	Rx chan byte // Reading (Acia Input) line
	Tx chan byte // Transmitting (Acia Output) line

	rxData byte
	txData byte

	commandData byte
	controlData byte

	rxFull  bool
	txEmpty bool

	rxIrqEnabled bool
	txIrqEnabled bool

	overrun bool
}

func NewAcia6551(rx chan byte, tx chan byte) (*Acia6551, error) {
	acia := &Acia6551{Tx: tx, Rx: rx}
	acia.Reset()

	go func() {
		// Handle rx data channel
	}()

	go func() {
		// Handle tx data channel
	}()

	return acia, nil
}

func (a *Acia6551) Size() uint16 {
	// We have a only 4 addresses, Data, Status, Command and Control
	return 0x04
}

// Emulates a hardware reset
func (a *Acia6551) Reset() {
	a.rxData = 0
	a.rxFull = false

	a.txData = 0
	a.txEmpty = true

	a.rxIrqEnabled = false
	a.txIrqEnabled = false

	a.overrun = false

	a.setControl(0)
	a.setCommand(0)
}

func (a *Acia6551) setControl(data byte) {
}

func (a *Acia6551) setCommand(data byte) {
}

func (a *Acia6551) Read(address uint16) byte {
	switch address {
	case aciaData:
		// Read Rx
	case aciaStatus:
		// Read Status reg.
	case aciaCommand:
		// Read command
	case aciaControl:
		// Read control
	}

	return 0x00
}

func (a *Acia6551) Write(address uint16, data byte) {
	switch address {
	case aciaData:
		// Write Tx
	case aciaStatus:
		// Reset
	case aciaCommand:
		// Write command
	case aciaControl:
		// Write control
	}
}
