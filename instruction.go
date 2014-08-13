package i6502

import (
	"fmt"
)

type Instruction struct {
	// Embed OpType
	OpType

	// 8-bit operand for 2-byte instructions
	Op8 byte

	// 16-bit operand for 3-byte instructions
	Op16 uint16

	// Address location where this instruction got read
	Address uint16
}

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

func (c *Cpu) readNextInstruction() Instruction {
	// Read the opcode
	opcode := c.bus.Read(c.PC)

	optype, ok := opTypes[opcode]
	if !ok {
		panic(fmt.Sprintf("Unknown or unimplemented opcode 0x%02X", opcode))
	}

	instruction := Instruction{OpType: optype}
	switch instruction.Size {
	case 1: // Zero operand instruction
	case 2: // 8-bit operand
		instruction.Op8 = c.bus.Read(c.PC + 1)
	case 3: // 16-bit operand
		instruction.Op16 = c.bus.Read16(c.PC + 1)
	}

	return instruction
}
