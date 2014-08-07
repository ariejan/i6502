# i6502 - A 6502/65C02 Emulator

The i6502 is a emulator/soft-prototype of a hardward device I'm building.

The goal of this project is to learn more about the following:

 * Go
 * CPU/Microprocessor Design
 * Computer Architecture
 * Assembler / Low-Level C
 * Operating Systems
 * Electronics (the hardware building part)

A test ROM file is included, but it does little more than loading a value into
the accumulator and storing it in memory.

## What's included in the emulator?

 * 6502 (not fully 65C02 yet) CPU
 * 16-bit address bus
 * 32kB RAM and 16kB ROM modules, addressable via the address bus
 * ROM loadable from file

## What's not (yet) included?

 * 65C02 support
 * I/O (6522, 6551)
 * Batteries
 * Tests ;-)

## License

This project is licensed under the MIT, see [LICENSE](https://github.com/ariejan/i6502/blob/master/LICENSE) for full details.

## Contributors

 * Ariejan de Vroom (ariejan@ariejan.net)

