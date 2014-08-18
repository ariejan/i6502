# i6502 - A 6502 Emulator

[![Build Status](https://travis-ci.org/ariejan/i6502.svg?branch=master)](https://travis-ci.org/ariejan/i6502)

This is an emulator of the i6502 hardward project I'm doing.

It's written in Golang and comes fully tested.

[Website](http://ariejan.github.io/i6502/) • [Documentation](http://godoc.org/github.com/ariejan/i6502)

## Background

The MOS 6502 Microprocessor has been around sinc 1975 and is used in many popular systems, like
the Apple II, Atari 2600, Commodore 64 and the Nintendo Entertainment System (NES).

It features an 8-bit accumulator and ALU, two 8-bit index registers and a 16-bit memory bus, allowing the processor to access up to 64kB of memory. 

I/O is mapped to memory, meaning that both RAM, ROM and I/O are addressed over the same 16-bit address bus.

Because of it's simple and elegant design and the plethora of information available about this microprocessor, the 6502 is very useful for learning and hobby projects.

The 65C02 is a updated version of the 6502. It includes some bug fixes and new instructions. The goal is for i6502 to fully support the 65C02 as well.

## What's included in the emulator?

 * 6502 Microprocessor, fully tested
 * 16-bit address bus, with attachable memory
 * RAM Memory

## What's not (yet) included?

 * Proper Golang packaging and documentation
 * 65C02 support
 * Roms
 * I/O (VIA 6522, ACIA 6551)
 * Batteries
 
## Getting started

Although this package contains everything you need, there is not single 'emulator' program yet. The CPU, address bus and RAM components are all available to you, but documentation is still lacking.

For now, you can checkout the project, and run the tests.

    go get github.com/ariejan/i6502
    cd $GOPATH/src/github.com/ariejan/i6502
    go get -t
    go test

## License

This project is licensed under the MIT, see [LICENSE](https://github.com/ariejan/i6502/blob/master/LICENSE) for full details.

## Contributors

 * Ariejan de Vroom (ariejan@ariejan.net)
 
 Special thanks to the awesome folk at [http://forum.6502.org](http://forum.6502.org) for their support and shared knowledge.
