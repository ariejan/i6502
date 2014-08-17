/*
The i6502 package contains all the components needed to construct
a working MOS 6502 emulated computer using different common parts,
like the MOS 6502, WDC 65C02, VIA 6522 and ACIA 6551.

The CPU is the core of the system. It features 8-bit registers and
ALU, and can address 16-bit of memory. It features a 16-bit program
counter (PC) that indicates where from memory the next instruction will
be read.

Besides the Cpu, there is also an AddressBus, which maps the 16-bit
address space to different attached components that implement the Memory
interface. Ram is one such component.

Creating a new emulator is easy and straightforward. All that's required
is a Cpu, and AddressBus and attached components.

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

Running a program from RAM is possible by loading it into
memory at the specified address. Note that this also sets the
Program Counter to the beginning of the loaded program.

Keep in mind that 0x00xx is reserved for Zeropage instructions and
0x01xx is reserved for the stack.

    import "io/ioutil"

    program, err := ioutil.ReadFile(path)

    // This will load the program (if it fits within memory)
    // at 0x0200 and set cpu.PC to 0x0200 as well.
    cpu.LoadProgram(program, 0x0200)

Running a program is as easy as calling `cpu.Step()`, which will
read and execute a single instruction.

    go for {
        cpu.Step()
    }()

*/
package i6502
