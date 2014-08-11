package devices

import (
	"fmt"
	"time"
)

const (
	aciaData = iota
	aciaStatus
	aciaCommand
	aciaControl
)

var baudRateSelectors = [...]int{0, 50, 75, 110, 135, 150, 300, 600, 1200, 1800, 2400, 3600, 4800, 7200, 9600, 19200}

type Acia6551 struct {
	// Registers
	rx              byte
	tx              byte
	commandRegister byte
	controlRegister byte

	// Other required bits and pieces
	lastTxWrite   int64
	lastRxRead    int64
	overrun       bool
	baudRate      int
	baudRateDelay int64
	rxFull        bool
	txEmpty       bool

	rxInterruptEnabled bool
	txInterruptEnabled bool

	InterruptChan chan bool

	RxChan chan byte
	TxChan chan byte
}

func NewAcia6551(cpu *Cpu) *Acia6551 {
	acia := &Acia6551{}
	acia.Reset()
	return acia
}

func (a *Acia6551) Reset() {
	a.RxChan = make(chan byte, 4096)
	a.TxChan = make(chan byte, 4096)

	a.tx = 0
	a.txEmpty = true

	a.rx = 0
	a.rxFull = false

	a.lastTxWrite = 0
	a.lastRxRead = 0
	a.overrun = false

	a.rxInterruptEnabled = false
	a.txInterruptEnabled = false

	a.InterruptChan = make(chan bool, 0)
}

func (a *Acia6551) Size() int {
	return 4
}

func (a *Acia6551) Read(address uint16) byte {
	switch address {
	case aciaData:
		return a.rxRead()
	case aciaStatus:
		return a.statusRegister()
	case aciaCommand:
		return a.commandRegister
	case aciaControl:
		return a.controlRegister
	default:
		panic(fmt.Errorf("ACIA 6551 cannot handle addressing 0x%04X", address))
	}
}

func (a *Acia6551) rxRead() byte {
	a.lastRxRead = unixTime()
	a.overrun = false
	a.rxFull = false
	return a.rx
}

func (a *Acia6551) RxWrite(data byte) {
	// Oh noes!
	if a.rxFull {
		a.overrun = true
	}

	a.rx = data
	a.rxFull = true

	if a.rxInterruptEnabled {
		// getbus.assertIrq()
	}
}

func (a *Acia6551) statusRegister() byte {
	now := unixTime()
	status := byte(0)

	if a.rxFull && (now >= (a.lastRxRead + a.baudRateDelay)) {
		status |= 0x08
	}

	if a.txEmpty && (now >= (a.lastTxWrite + a.baudRateDelay)) {
		status |= 0x10
	}

	if a.overrun {
		status |= 0x04
	}

	return status
}

func (a *Acia6551) Write(address uint16, value byte) {
	switch address {
	case aciaData:
		a.txWrite(value)
	case aciaStatus:
		a.Reset()
	case aciaCommand:
		a.setCommandRegister(value)
	case aciaControl:
		a.setControlRegister(value)
	default:
		panic(fmt.Errorf("ACIA 6551 cannot handle addressing 0x%04X", address))
	}
}

func (a *Acia6551) txWrite(value byte) {
	a.lastTxWrite = unixTime()
	a.tx = value
	a.txEmpty = false

	// Post for others
	a.debugTxOutput()
}

func (a *Acia6551) TxRead() byte {
	a.txEmpty = true
	return a.tx
}

func (a *Acia6551) HasTx() bool {
	return !a.txEmpty
}

func (a *Acia6551) HasRx() bool {
	return a.rxFull
}

func (a *Acia6551) debugTxOutput() {
	if a.HasTx() {
		a.TxChan <- a.TxRead()
	}
}

func (a *Acia6551) setCommandRegister(data byte) {
	fmt.Printf("Setting Acia6551 Command Register: %02X\n", data)
	a.commandRegister = data
	// TODO: Maybe implement IRQs.
}

func (a *Acia6551) setControlRegister(data byte) {
	a.controlRegister = data

	if data == 0x00 {
		a.Reset()
	} else {
		a.setBaudRate(baudRateSelectors[data&0x0f])
	}
}

func (a *Acia6551) setBaudRate(baudRate int) {
	fmt.Printf("Setting baudrate at %d\n", baudRate)

	a.baudRate = baudRate

	// Set baudRateDelay in nanoseconds. It's an approximation.
	if baudRate > 0 {
		a.baudRateDelay = int64((1.0 / float64(baudRate)) * 1000000000 * 8)
	} else {
		a.baudRateDelay = 0
	}
}

func unixTime() int64 {
	return time.Now().UnixNano()
}
