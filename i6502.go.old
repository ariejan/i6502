/*
i6502 is a software emulator of my i6502 home-built computer. It uses
the 65C02 microprocessor, 32kB RAM and 16kB ROM.
*/
package main

import (
	"fmt"
	"github.com/ariejan/i6502/bus"
	"github.com/ariejan/i6502/cpu"
	"github.com/ariejan/i6502/devices"
	"github.com/ariejan/i6502/memory"
	"os"
	"os/signal"
)

func main() {
	os.Exit(mainReturningStatus())
}

func mainReturningStatus() int {
	// 32kB RAM
	ram := memory.CreateRam()

	// 16kB ROM, filled from file
	rom, err := memory.LoadRomFromFile("rom/ehbasic.rom")
	if err != nil {
		panic(err)
	}

	acia6551 := devices.NewAcia6551()

	// 16-bit address bus
	bus, _ := bus.CreateBus()
	bus.Attach(ram, "32kB RAM", 0x0000)
	bus.Attach(rom, "16kB ROM", 0xC000)
	bus.Attach(acia6551, "ACIA 6551", 0x8800)

	fmt.Println(bus)

	exitChan := make(chan int, 0)

	cpu := &cpu.Cpu{Bus: bus, ExitChan: exitChan}
	cpu.Reset()

	go serialServer(cpu, acia6551)

	go func() {
		for {
			cpu.Step()
		}
	}()

	var (
		sig        os.Signal
		exitStatus int
	)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	select {
	case exitStatus = <-exitChan:
		// Okay, handle the rest of the code
	case sig = <-sigChan:
		fmt.Println("\nGot signal: ", sig)
		exitStatus = 1
	}

	fmt.Println(cpu)
	fmt.Println("Dumping RAM into `core` file")
	ram.Dump("core")

	return exitStatus
}
