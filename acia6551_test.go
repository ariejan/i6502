package i6502

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewAcia6551(t *testing.T) {
	tx := make(chan byte)
	rx := make(chan byte)
	acia, err := NewAcia6551(rx, tx)

	assert.Nil(t, err)
	assert.Equal(t, 0x4, acia.Size())
}
