package i6502

import (
	"fmt"
)

type Instruction struct {
	OpType // Embed OpType

	Op8  byte   // 8-bit operand for 2-byte instructions
	Op16 uint16 // 16-bit operand for 3-byte instructions

	Address uint16 // Address location where this instruction got read, for debugging purposes
}

// Return a string containing debug information about the instruction and operands.
func (i Instruction) String() (output string) {
	switch i.Size {
	case 1:
		output = fmt.Sprintf("~~~ 0x%04X: 0x%02X - %s [%s] {%d}\n", i.Address, i.Opcode, instructionNames[i.opcodeId], addressingNames[i.addressingId], i.Cycles)
	case 2:
		output = fmt.Sprintf("~~~ 0x%04X: 0x%02X - %s 0x%02X [%s] {%d}\n", i.Address, i.Opcode, instructionNames[i.opcodeId], i.Op8, addressingNames[i.addressingId], i.Cycles)
	case 3:
		output = fmt.Sprintf("~~~ 0x%04X: 0x%02X - %s 0x%04X [%s] {%d}\n", i.Address, i.Opcode, instructionNames[i.opcodeId], i.Op16, addressingNames[i.addressingId], i.Cycles)
	}

	return
}
