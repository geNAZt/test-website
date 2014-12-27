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
	Connection *Connection

	// Message
	Message []byte
}

type hub struct {
	// Registered connections.
	connections map[*Connection]bool

	// Inbound messages from the connections.
	Broadcast chan []byte

	// Register requests from the connections.
	register chan *Connection

	// Unregister requests from connections.
	unregister chan *Connection

	// Message channel for incoming Messages
	messages chan Message

	// EventSystem for incoming Messages
	eventSystem *EventSystem
}

var Hub = &hub{
	Broadcast:   make(chan []byte, maxMessageSize),
	register:    make(chan *Connection, 1),
	unregister:  make(chan *Connection, 1),
	connections: make(map[*Connection]bool),
	messages:    make(chan Message),
	eventSystem: &EventSystem{},
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
	go Hub.run()
}

func (h *hub) run() {
	for {
		select {
		case c := <-h.register:
			h.connections[c] = true
		case c := <-h.unregister:
			c.CloseCustomChannels()
			close(c.Send)
			delete(h.connections, c)
		case m := <-h.Broadcast:
			for c := range h.connections {
				select {
				case c.Send <- m:
				default:
					c.CloseCustomChannels()
					close(c.Send)
					delete(h.connections, c)
				}
			}
		case m := <-h.messages:
			h.eventSystem.Emit(m)
		}
	}
}

func (h *hub) Listen(fn MessageAccepter) <-chan Message {
	return h.eventSystem.Listen(fn)
}
