package websocket

import (
	"math/rand"
	"time"
)

const (
	// Time allowed to write a message to the client.
	writeWait = 10 * time.Second

	// Time allowed to read the next message from the client.
	readWait = 60 * time.Second

	// Send pings to client with this period. Must be less than readWait.
	pingPeriod = (readWait * 9) / 10

	// Maximum message size allowed from client.
	maxMessageSize = 512
)

type Message struct {
	// Connection which has sent this message
	connection Connection

	// Message
	message []byte
}

type hub struct {
	// Registered connections.
	connections map[*Connection]bool

	// Inbound messages from the connections.
	broadcast chan []byte

	// Register requests from the connections.
	register chan *Connection

	// Unregister requests from connections.
	unregister chan *Connection

	// Message channel for incoming Messages
	messages chan *Message
}

var h = &hub{
	broadcast:   make(chan []byte, maxMessageSize),
	register:    make(chan *Connection, 1),
	unregister:  make(chan *Connection, 1),
	connections: make(map[*Connection]bool),
	messages:    make(chan *Message),
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
	go h.run()
}

func (h *hub) run() {
	for {
		select {
		case c := <-h.register:
			h.connections[c] = true
		case c := <-h.unregister:
			delete(h.connections, c)
			close(c.send)
		case m := <-h.broadcast:
			for c := range h.connections {
				select {
				case c.send <- m:
				default:
					close(c.send)
					delete(h.connections, c)
				}
			}
		}
	}
}
