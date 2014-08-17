package i6502

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func AciaSubject() *Acia6551 {
	tx := make(chan byte)
	rx := make(chan byte)
	acia, _ := NewAcia6551(rx, tx)

	return acia
}

func TestNewAcia6551(t *testing.T) {
	tx := make(chan byte)
	rx := make(chan byte)
	acia, err := NewAcia6551(rx, tx)

	assert.Nil(t, err)
	assert.Equal(t, 0x4, acia.Size())
}

func TestAciaReset(t *testing.T) {
	a := AciaSubject()
	a.Reset()

	assert.Equal(t, a.txData, 0)
	assert.True(t, a.txEmpty)

	assert.Equal(t, a.rxData, 0)
	assert.False(t, a.rxFull)

	assert.False(t, a.txIrqEnabled)
	assert.False(t, a.rxIrqEnabled)
}

func TestAciaCommand(t *testing.T) {
}
