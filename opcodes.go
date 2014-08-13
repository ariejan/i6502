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

	0x69: OpType{0x69, adc, immediate, 2, 2},
	0x65: OpType{0x65, adc, zeropage, 2, 2},
	0x75: OpType{0x75, adc, zeropageX, 2, 4},
	0x6D: OpType{0x6D, adc, absolute, 3, 4},
	0x7D: OpType{0x7D, adc, absoluteX, 3, 4},
	0x79: OpType{0x79, adc, absoluteY, 3, 4},
	0x61: OpType{0x61, adc, indirectX, 2, 6},
	0x71: OpType{0x71, adc, indirectY, 2, 5},
}
