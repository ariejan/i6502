package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

const (
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

type connection struct {
	// The i6502 machine
	machine *Machine

	// Websocket connection
	ws *websocket.Conn

	// Outgoing data channel
	send chan []byte
}

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

// Pump received websocket messages into the machine
func (c *connection) readPump() {
	defer func() {
		c.ws.Close()
	}()

	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			break
		}
		c.writeBytesToMachine(message)
	}
}

// Pump serial output from the machine into the socket
func (c *connection) writePump() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		c.ws.Close()
	}()

	for {
		select {
		case data, ok := <-c.machine.SerialTx:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, []byte{data}); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func (c *connection) write(messageType int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(10 * time.Second))
	return c.ws.WriteMessage(messageType, payload)
}

func (c *connection) writeBytesToMachine(data []byte) {
	for _, b := range data {
		log.Printf("%c", b)
		c.machine.SerialRx <- b
	}
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	c := &connection{machine: CreateMachine(), send: make(chan []byte, 256), ws: ws}

	go c.writePump()

	c.machine.Reset()
	c.readPump()
}
