package i6502

// OpCode table
const (
	_ = iota
	adc
	and
	asl
	bcc
	bcs
	beq
	bit
	bmi
	bne
	bpl
	brk
	bvc
	bvs
	clc
	cld
	cli
	clv
	cmp
	cpx
	cpy
	dec
	dex
	dey
	eor
	inc
	inx
	iny
	jmp
	jsr
	lda
	ldx
	ldy
	lsr
	nop
	ora
	pha
	php
	pla
	plp
	rol
	ror
	rti
	rts
	sbc
	sec
	sed
	sei
	sta
	stx
	sty
	tax
	tay
	tsx
	txa
	txs
	tya
)

var instructionNames = [...]string{
	"",
	"ADC",
	"AND",
	"ASL",
	"BCC",
	"BCS",
	"BEQ",
	"BIT",
	"BMI",
	"BNE",
	"BPL",
	"BRK",
	"BVC",
	"BVS",
	"CLC",
	"CLD",
	"CLI",
	"CLV",
	"CMP",
	"CPX",
	"CPY",
	"DEC",
	"DEX",
	"DEY",
	"EOR",
	"INC",
	"INX",
	"INY",
	"JMP",
	"JSR",
	"LDA",
	"LDX",
	"LDY",
	"LSR",
	"NOP",
	"ORA",
	"PHA",
	"PHP",
	"PLA",
	"PLP",
	"ROL",
	"ROR",
	"RTI",
	"RTS",
	"SBC",
	"SEC",
	"SED",
	"SEI",
	"STA",
	"STX",
	"STY",
	"TAX",
	"TAY",
	"TSX",
	"TXA",
	"TXS",
	"TYA",
}

// OpType is the operation type, it includes the instruction and
// addressing mode. It also includes some extra information for the
// emulator, like number of cycles
type OpType struct {
	// The actual Opcode byte read from memory
	Opcode byte

	// Opcode ID
	opcodeId uint8

	// Addressing Mode ID
	addressingId uint8

	// Size of this instruction, either 1, 2 or 3 bytes
	Size uint8

	// Number of clock cycles this instruction needs
	Cycles uint8
}

var opTypes = map[uint8]OpType{
	0xEA: OpType{0xEA, nop, implied, 1, 2},

	// Set instruction
	0x38: OpType{0x38, sec, implied, 1, 2},
	0xF8: OpType{0xF8, sed, implied, 1, 2},
	0x78: OpType{0x78, sei, implied, 1, 2},

	// Clear instructions
	0x18: OpType{0x18, clc, implied, 1, 2},
	0xD8: OpType{0xD8, cld, implied, 1, 2},
	0x58: OpType{0x58, cli, implied, 1, 2},
	0xB8: OpType{0xB8, clv, implied, 1, 2},

	// ADC
	0x69: OpType{0x69, adc, immediate, 2, 2},
	0x65: OpType{0x65, adc, zeropage, 2, 3},
	0x75: OpType{0x75, adc, zeropageX, 2, 4},
	0x6D: OpType{0x6D, adc, absolute, 3, 4},
	0x7D: OpType{0x7D, adc, absoluteX, 3, 4},
	0x79: OpType{0x79, adc, absoluteY, 3, 4},
	0x61: OpType{0x61, adc, indirectX, 2, 6},
	0x71: OpType{0x71, adc, indirectY, 2, 5},

	// SBC
	0xE9: OpType{0xE9, sbc, immediate, 2, 2},
	0xE5: OpType{0xE5, sbc, zeropage, 2, 3},
	0xF5: OpType{0xF5, sbc, zeropageX, 2, 4},
	0xED: OpType{0xED, sbc, absolute, 3, 4},
	0xFD: OpType{0xFD, sbc, absoluteX, 3, 4},
	0xF9: OpType{0xF9, sbc, absoluteY, 3, 4},
	0xE1: OpType{0xE1, sbc, indirectX, 2, 6},
	0xF1: OpType{0xF1, sbc, indirectY, 2, 5},

	// Increments
	0xE8: OpType{0xE8, inx, implied, 1, 2},
	0xC8: OpType{0xC8, iny, implied, 1, 2},

	0xE6: OpType{0xE6, inc, zeropage, 2, 5},
	0xF6: OpType{0xF6, inc, zeropageX, 2, 6},
	0xEE: OpType{0xEE, inc, absolute, 3, 6},
	0xFE: OpType{0xFE, inc, absoluteX, 3, 7},

	// Decrements
	0xCA: OpType{0xCA, dex, implied, 1, 2},
	0x88: OpType{0x88, dey, implied, 1, 2},

	0xC6: OpType{0xC6, dec, zeropage, 2, 5},
	0xD6: OpType{0xD6, dec, zeropageX, 2, 6},
	0xCE: OpType{0xCE, dec, absolute, 3, 6},
	0xDE: OpType{0xDE, dec, absoluteX, 3, 7},

	// LDA
	0xA9: OpType{0xA9, lda, immediate, 2, 2},
	0xA5: OpType{0xA5, lda, zeropage, 2, 3},
	0xB5: OpType{0xB5, lda, zeropageX, 2, 4},
	0xAD: OpType{0xAD, lda, absolute, 3, 4},
	0xBD: OpType{0xBD, lda, absoluteX, 3, 4},
	0xB9: OpType{0xB9, lda, absoluteY, 3, 4},
	0xA1: OpType{0xA1, lda, indirectX, 2, 6},
	0xB1: OpType{0xB1, lda, indirectY, 2, 5},

	// LDX
	0xA2: OpType{0xA2, ldx, immediate, 2, 2},
	0xA6: OpType{0xA6, ldx, zeropage, 2, 3},
	0xB6: OpType{0xB6, ldx, zeropageY, 2, 4},
	0xAE: OpType{0xAE, ldx, absolute, 3, 4},
	0xBE: OpType{0xBE, ldx, absoluteY, 3, 4},

	// LDY
	0xA0: OpType{0xA0, ldy, immediate, 2, 2},
	0xA4: OpType{0xA4, ldy, zeropage, 2, 3},
	0xB4: OpType{0xB4, ldy, zeropageY, 2, 4},
	0xAC: OpType{0xAC, ldy, absolute, 3, 4},
	0xBC: OpType{0xBC, ldy, absoluteY, 3, 4},

	// ORA
	0x09: OpType{0x09, ora, immediate, 2, 2},
	0x05: OpType{0x05, ora, zeropage, 2, 3},
	0x15: OpType{0x15, ora, zeropageX, 2, 4},
	0x0D: OpType{0x0D, ora, absolute, 3, 4},
	0x1D: OpType{0x1D, ora, absoluteX, 3, 4},
	0x19: OpType{0x19, ora, absoluteY, 3, 4},
	0x01: OpType{0x01, ora, indirectX, 2, 6},
	0x11: OpType{0x11, ora, indirectY, 2, 5},

	// AND
	0x29: OpType{0x29, and, immediate, 2, 2},
	0x25: OpType{0x25, and, zeropage, 2, 3},
	0x35: OpType{0x35, and, zeropageX, 2, 4},
	0x2D: OpType{0x2D, and, absolute, 3, 4},
	0x3D: OpType{0x3D, and, absoluteX, 3, 4},
	0x39: OpType{0x39, and, absoluteY, 3, 4},
	0x21: OpType{0x21, and, indirectX, 2, 6},
	0x31: OpType{0x31, and, indirectY, 2, 5},

	// EOR
	0x49: OpType{0x49, eor, immediate, 2, 2},
	0x45: OpType{0x45, eor, zeropage, 2, 3},
	0x55: OpType{0x55, eor, zeropageX, 2, 4},
	0x4D: OpType{0x4D, eor, absolute, 3, 4},
	0x5D: OpType{0x5D, eor, absoluteX, 3, 4},
	0x59: OpType{0x59, eor, absoluteY, 3, 4},
	0x41: OpType{0x41, eor, indirectX, 2, 6},
	0x51: OpType{0x51, eor, indirectY, 2, 5},

	// STA
	0x85: OpType{0x85, sta, zeropage, 2, 3},
	0x95: OpType{0x95, sta, zeropageX, 2, 4},
	0x8D: OpType{0x8D, sta, absolute, 3, 4},
	0x9D: OpType{0x9D, sta, absoluteX, 3, 5},
	0x99: OpType{0x99, sta, absoluteY, 3, 5},
	0x81: OpType{0x81, sta, indirectX, 2, 6},
	0x91: OpType{0x91, sta, indirectY, 2, 6},

	// STX
	0x86: OpType{0x86, stx, zeropage, 2, 3},
	0x96: OpType{0x96, stx, zeropageY, 2, 4},
	0x8E: OpType{0x8E, stx, absolute, 3, 4},

	// STY
	0x84: OpType{0x84, sty, zeropage, 2, 3},
	0x94: OpType{0x94, sty, zeropageX, 2, 4},
	0x8C: OpType{0x8C, sty, absolute, 3, 4},

	// TAX
	0xAA: OpType{0xAA, tax, implied, 1, 2},

	// TAY
	0xA8: OpType{0xA8, tay, implied, 1, 2},

	// TXA
	0x8A: OpType{0x8A, txa, implied, 1, 2},

	// TXA
	0x98: OpType{0x98, tya, implied, 1, 2},

	// TSX
	0xBA: OpType{0xBA, tsx, implied, 1, 2},

	// TXS
	0x9A: OpType{0x9A, txs, implied, 1, 2},

	// ASL
	0x0A: OpType{0x0A, asl, accumulator, 1, 2},
	0x06: OpType{0x06, asl, zeropage, 2, 5},
	0x16: OpType{0x16, asl, zeropageX, 2, 6},
	0x0E: OpType{0x0E, asl, absolute, 3, 6},
	0x1E: OpType{0x1E, asl, absoluteX, 3, 7},

	// LSR
	0x4A: OpType{0x4A, lsr, accumulator, 1, 2},
	0x46: OpType{0x46, lsr, zeropage, 2, 5},
	0x56: OpType{0x56, lsr, zeropageX, 2, 6},
	0x4E: OpType{0x4E, lsr, absolute, 3, 6},
	0x5E: OpType{0x5E, lsr, absoluteX, 3, 7},
}
