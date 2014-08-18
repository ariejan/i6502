package i6502

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func AciaSubject() (*Acia6551, chan byte, chan byte) {
	tx := make(chan byte)
	rx := make(chan byte)
	acia, _ := NewAcia6551(rx, tx)

	return acia, rx, tx
}

func TestNewAcia6551(t *testing.T) {
	tx := make(chan byte)
	rx := make(chan byte)
	acia, err := NewAcia6551(rx, tx)

	assert.Nil(t, err)
	assert.Equal(t, 0x4, acia.Size())
}

func TestAciaReset(t *testing.T) {
	a, _, _ := AciaSubject()

	a.Reset()

	assert.Equal(t, a.txData, 0)
	assert.True(t, a.txEmpty)

	assert.Equal(t, a.rxData, 0)
	assert.False(t, a.rxFull)

	assert.False(t, a.txIrqEnabled)
	assert.False(t, a.rxIrqEnabled)

	assert.False(t, a.overrun)
	assert.Equal(t, 0, a.controlData)
}

func TestAciaReadData(t *testing.T) {
	a, _, _ := AciaSubject()

	a.Rx <- 0x42

	assert.True(t, a.rxFull)

	fmt.Printf("Reading...\n")
	value := a.Read(aciaData)
	assert.Equal(t, 0x42, value)
	assert.False(t, a.rxFull)
}
