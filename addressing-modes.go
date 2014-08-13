package i6502

// Addressing modes
const (
	_ = iota
	absolute
	absoluteX
	absoluteY
	accumulator
	immediate
	implied
	indirect
	indirectX
	indirectY
	relative
	zeropage
	zeropageX
	zeropageY
)

var addressingNames = [...]string{
	"",
	"absolute",
	"absoluteX",
	"absoluteY",
	"accumulator",
	"immediate",
	"implied",
	"(indirect)",
	"(indirect,X)",
	"(indirect),Y",
	"relative",
	"zeropage",
	"zeropageX",
	"zeropageY",
}
