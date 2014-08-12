package main

import (
	"github.com/ariejan/i6502/bus"
	"github.com/ariejan/i6502/cpu"
	"github.com/ariejan/i6502/devices"
	"github.com/ariejan/i6502/memory"
	"os"
)

type Machine struct {
	// Outgoing bytes, using the serial interface
	SerialTx chan byte

	// Incoming bytes, using the serial interface
	SerialRx chan byte

	// The cpu, bus etc.
	cpu *cpu.Cpu
	bus *bus.Bus
}

// Creates a new i6502 Machine instance
func CreateMachine() *Machine {
	// Channel for handling interrupts
	irqChan := make(chan bool, 0)
	nmiChan := make(chan bool, 0)

	ram := memory.CreateRam()

	rom, err := memory.LoadRomFromFile("rom/ehbasic.rom")
	if err != nil {
		panic(err)
	}

	acia6551 := devices.NewAcia6551(irqChan)
	via6522 := devices.NewVia6522(irqChan)

	bus, _ := bus.CreateBus()
	bus.Attach(ram, "32kB RAM", 0x0000)
	bus.Attach(rom, "16kB ROM", 0xC000)
	bus.Attach(via6522, "VIA 6522 Parallel", 0x8000)
	bus.Attach(acia6551, "ACIA 6551 Serial", 0x8800)

	cpu := &cpu.Cpu{Bus: bus, IrqChan: irqChan, NmiChan: nmiChan, ExitChan: make(chan int, 0)}

	machine := &Machine{SerialTx: make(chan byte, 256), SerialRx: make(chan byte, 256), cpu: cpu, bus: bus}

	// Run the CPU
	go func() {
		for {
			cpu.Step()
		}
	}()

	// Connect acia6551 Tx to SerialTx
	go func() {
		for {
			select {
			case data := <-acia6551.TxChan:
				machine.SerialTx <- data
			}
		}
	}()

	// Connect SerialRx to acia6551 Rx
	go func() {
		for {
			select {
			case data := <-machine.SerialRx:
				acia6551.RxWrite(data)
			}
		}
	}()

	go func() {
		for {
			select {
			case <-machine.cpu.ExitChan:
				ram.Dump("intcore")
				os.Exit(42)
			}
		}
	}()

	cpu.Reset()

	return machine
}

func (m *Machine) Reset() {
	m.cpu.Reset()
}
