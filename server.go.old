package main

import (
	"bufio"
	"fmt"
	"github.com/ariejan/i6502/cpu"
	"github.com/ariejan/i6502/devices"
	"net"
)

func handleConnection(c net.Conn, cpu *cpu.Cpu, acia *devices.Acia6551) {

	// Force telnet into character mode
	c.Write([]byte("\377\375\042\377\373\001"))
	c.Write([]byte("-- i6502 Serial Terminal --\n"))

	// Transfer output to the client
	go func() {
		for {
			select {
			case data := <-acia.TxChan:
				c.Write([]byte{data})
			}
		}
	}()

	go func() {
		reader := bufio.NewReader(c)

		for {
			b, err := reader.ReadByte()
			if err != nil {
				panic(err)
			}

			// Push to CPU
			acia.RxChan <- b
		}
	}()

	fmt.Println("Client connected. Resetting CPU.")
	cpu.Reset()
}

func serialServer(cpu *cpu.Cpu, acia *devices.Acia6551) {
	listen, err := net.Listen("tcp", ":6000")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listen.Accept()
		defer conn.Close()

		if err != nil {
			continue
		}
		go handleConnection(conn, cpu, acia)
	}
}
