package cpu

import (
	"fmt"
	"github.com/ariejan/i6502/bus"
)

// Instruction is a self-contained, variable length instruction,
// it includes the the operation type and a possible 8 or 16 bit operand.
type Instruction struct {
	OpType

	// 2 byte instructions have a single byte operand
	Op8 uint8

	// 3 byte instruction have a double byte operand
	Op16 uint16
}

func (i *Instruction) String() string {
	var output string

	switch i.Bytes {
	case 1:
		output = fmt.Sprintf("0x%02X - %s\n", i.Opcode, instructionNames[i.id])
	case 2:
		output = fmt.Sprintf("0x%02X - %s %02X\n", i.Opcode, instructionNames[i.id], i.Op8)
	case 3:
		output = fmt.Sprintf("0x%02X - %s %04X\n", i.Opcode, instructionNames[i.id], i.Op16)
	}

	return output
}

func ReadInstruction(pc uint16, bus *bus.Bus) Instruction {
	// Read the opcode
	opcode := bus.Read(pc)

	// Do we know this opcode in our optypes table?
	optype, ok := optypes[opcode]
	if !ok {
		panic(fmt.Sprintf("Unknown opcode $%02X at $%04X", opcode, pc))
	}

	instruction := Instruction{OpType: optype}
	switch instruction.Bytes {
	case 1: // No operand.
	case 2:
		instruction.Op8 = bus.Read(pc + 1)
	case 3:
		instruction.Op16 = bus.Read16(pc + 1)
	default:
		panic(fmt.Sprintf("Unknown instruction length (%d) for $%02X at $04X", instruction.Bytes, opcode, pc))
	}

	return instruction
}
