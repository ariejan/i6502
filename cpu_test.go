package i6502

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewCpu(t *testing.T) {
	cpu, err := NewCpu()

	assert.NotNil(t, cpu)
	assert.Nil(t, err)
}
