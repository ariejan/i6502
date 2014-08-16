package i6502

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStackPushPopPeek(t *testing.T) {
	assert := assert.New(t)
	cpu, _, _ := NewRamMachine()

	assert.Equal(0xFF, cpu.SP)

	cpu.stackPush(0x42)
	cpu.stackPush(0xA0)

	assert.Equal(0xFD, cpu.SP)
	assert.Equal(0x42, cpu.bus.Read(0x1FF))
	assert.Equal(0xA0, cpu.bus.Read(0x1FE))

	peekValue := cpu.stackPeek()
	assert.Equal(0xFD, cpu.SP)
	assert.Equal(0xA0, peekValue)

	popValue := cpu.stackPop()
	assert.Equal(0xFE, cpu.SP)
	assert.Equal(0xA0, popValue)
}
