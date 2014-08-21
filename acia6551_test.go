package i6502

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func AciaSubject() *Acia6551 {
	acia, _ := NewAcia6551()
	return acia
}

func TestNewAcia6551(t *testing.T) {
	acia, err := NewAcia6551()

	assert.Nil(t, err)
	assert.Equal(t, 0x4, acia.Size())
}

func TestAciaAsMemory(t *testing.T) {
	assert.Implements(t, (*Memory)(nil), new(Acia6551))
}

func TestAciaReset(t *testing.T) {
	a := AciaSubject()

	a.Reset()

	assert.Equal(t, a.tx, 0)
	assert.True(t, a.txEmpty)

	assert.Equal(t, a.rx, 0)
	assert.False(t, a.rxFull)

	assert.False(t, a.txIrqEnabled)
	assert.False(t, a.rxIrqEnabled)

	assert.False(t, a.overrun)
	assert.Equal(t, 0, a.controlData)
}

func TestAciaReaderWithTxEmpty(t *testing.T) {
	a := AciaSubject()

	// Nothing to read
	assert.True(t, a.txEmpty)

	value := make([]byte, 1)
	bytesRead, _ := a.Read(value)

	assert.Equal(t, 0, bytesRead)
}

func TestAciaWriteByteAndReader(t *testing.T) {
	a := AciaSubject()

	// CPU writes data
	a.WriteByte(aciaData, 0x42)

	// System reads from Tx
	value := make([]byte, 1)
	bytesRead, _ := a.Read(value)

	if assert.Equal(t, 1, bytesRead) {
		assert.Equal(t, 0x42, value[0])
	}
}

func TestAciaWriterAndReadByte(t *testing.T) {
	a := AciaSubject()

	// System writes a single byte
	bytesWritten, _ := a.Write([]byte{0x42})

	if assert.Equal(t, 1, bytesWritten) {
		assert.Equal(t, 0x42, a.ReadByte(aciaData))
	}

	// System writes multiple bytes
	bytesWritten, _ = a.Write([]byte{0x42, 0x32, 0xAB})

	if assert.Equal(t, 3, bytesWritten) {
		assert.Equal(t, 0xAB, a.ReadByte(aciaData))
	}
}

func TestAciaCommandRegister(t *testing.T) {
	a := AciaSubject()
	assert.False(t, a.rxIrqEnabled)
	assert.False(t, a.txIrqEnabled)

	a.WriteByte(aciaCommand, 0x02) // b0000 0010 RX Irq enabled
	assert.True(t, a.rxIrqEnabled)
	assert.False(t, a.txIrqEnabled)

	a.WriteByte(aciaCommand, 0x04) // b0000 0100 TX Irq enabled
	assert.False(t, a.rxIrqEnabled)
	assert.True(t, a.txIrqEnabled)

	a.WriteByte(aciaCommand, 0x06) // b0000 0110 RX + TX Irq enabled
	assert.True(t, a.rxIrqEnabled)
	assert.True(t, a.txIrqEnabled)

	assert.Equal(t, 0x06, a.ReadByte(aciaCommand))
}

func TestAciaControlRegister(t *testing.T) {
	a := AciaSubject()

	a.WriteByte(aciaControl, 0xB8)
	assert.Equal(t, 0xB8, a.ReadByte(aciaControl))
}

func TestAciaStatusRegister(t *testing.T) {
	a := AciaSubject()

	a.rxFull = false
	a.txEmpty = false
	a.overrun = false
	assert.Equal(t, 0x00, a.ReadByte(aciaStatus))

	a.rxFull = true
	a.txEmpty = false
	a.overrun = false
	assert.Equal(t, 0x08, a.ReadByte(aciaStatus))

	a.rxFull = false
	a.txEmpty = true
	a.overrun = false
	assert.Equal(t, 0x10, a.ReadByte(aciaStatus))

	a.rxFull = false
	a.txEmpty = false
	a.overrun = true
	assert.Equal(t, 0x04, a.ReadByte(aciaStatus))
}

func TestAciaIntegration(t *testing.T) {
	// Create a system
	// * 32kB RAM at 0x0000-7FFFF
	// * ACIA at 0x8800-8803
	ram, _ := NewRam(0x8000)
	acia, _ := NewAcia6551()
	bus, _ := NewAddressBus()
	bus.Attach(ram, 0x0000)
	bus.Attach(acia, 0x8800)
	cpu, _ := NewCpu(bus)

	program := []byte{
		0xA9, 0x00, // LDA #$00
		0x8D, 0x01, 0x88, // STA AciaStatus (Reset)
		0xA9, 0x42, // LDA #$42
		0x8D, 0x00, 0x88, // STA AciaData (Write)
		0xAD, 0x00, 0x88, // LDA AciaData (Read)
	}

	cpu.LoadProgram(program, 0x0200)
	cpu.Steps(2)

	acia.Write([]byte{0xAB})

	cpu.Steps(3)

	value := make([]byte, 1)
	bytesRead, _ := acia.Read(value)

	if assert.Equal(t, 1, bytesRead) {
		assert.Equal(t, 0x42, value[0])
	}

	assert.Equal(t, 0xAB, cpu.A)
}
