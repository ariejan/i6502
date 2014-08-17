/*
The i6502 package contains all the components needed to construct
a working MOS 6502 emulated computer using different common parts,
like the MOS 6502 or WDC 65C02, VIA 6522 (parallel I/O) and
ACIA 6551 (serial I/O).

The CPU is the core of the system. It features an 8-bit accumulator (A)
and two general purpose 8-bit index registers (X, Y). There is a
16-bit program counter (PC). The 8-bit stack pointer (SP) points to
the 0x0100-0x1FF address space moves downward. The status register (P)
contains bits indicating Zero, Negative, Break, Decimal, IrqDisable,
Carry and Overflow conditions. The 6502 uses a 16-bit address bus to
access 8-bit data values.

The AddressBus can be used to attach different components to different
parts of the 16-bit address space, accessible by the 6502. Common
layouts are

 * 64kB RAM at 0x0000-FFFF

Or

 * 32kB RAM  at 0x0000-7FFF
 * VIA 6522  at 0x8000-800F
 * ACIA 6551 at 0x8800-8803
 * 16kB ROM  at 0xC000-FFFF

Creating a new emulated machine entails three steps:

 1. Create the different memory components (Ram, Rom, IO)
 2. Create the AddressBus and attach memory
 3. Create the Cpu with the AddressBus

Example: create an emulator using the full 64kB address space for RAM

    import "github.com/ariejan/i6502"

    // Create Ram, 64kB in size
    ram, err := i6502.NewRam(0x10000)

    // Create the AddressBus
    bus, err := i6502.NewAddressBus()

    // And attach the Ram at 0x0000
    bus.Attach(ram, 0x0000)

    // Create the Cpu, with the AddressBus
    cpu, err := i6502.NewCpu(bus)

The hardware pins `IRQ` and `RESB` are implemented and mapped to
the functions `Interrupt()` and `Reset()`.

Running a program from memory can be done by loading the binary
data into memory using `LoadProgram`. Keep in mind that the first
two memory pages (0x0000-01FF) are reserved for zeropage and stack
memory.

Example of loading a binary program from disk into memory:

    import "io/ioutil"

    program, err := ioutil.ReadFile(path)

    // This will load the program (if it fits within memory)
    // at 0x0200 and set cpu.PC to 0x0200 as well.
    cpu.LoadProgram(program, 0x0200)

With all memory connected and a program loaded, all that's left
is executing instructions on the Cpu. A single call to `Step()` will
read and execute a single (1, 2 or 3 byte) instruction from memory.

To create a Cpu and have it running, simple create a go-routine.

    go for {
        cpu.Step()
    }()

*/
package i6502
