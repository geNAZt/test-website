package websocket

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/session"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// allow all connections by default
		return true
	},
}

func Upgrade(w beego.Controller) *Connection {
	// Try to update the Request
	ws, err := upgrader.Upgrade(w.Ctx.ResponseWriter, w.Ctx.Request, nil)
	if err != nil {
		beego.BeeLogger.Warn("Could not upgrade Websocket Request %v", err)
		return nil
	}

	// Be 100% sure that the Session is there
	if w.CruSession == nil {
		w.StartSession()
	}

	// Create new connection
	c := &Connection{
		Send:    make(chan []byte, 256),
		ws:      ws,
		Session: w.CruSession,
	}

	// Tell the Hub we have a new Connection and start to pump messages
	Hub.register <- c
	go c.writePump()
	go c.readPump()

	// When the connection closes save the session
	close := make(chan struct{}, 1)
	go func() {
		select {
		case <-close:
			c.Session.SessionRelease(w.Ctx.ResponseWriter)
		}
	}()
	c.AppendChannel(close)
	c.Open = true

	return c
}

// connection is an middleman between the websocket connection and the hub.
type Connection struct {
	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	Send chan []byte

	// Custom appended channels
	customChannels []chan struct{}

	// Session from HTTP Request
	Session session.SessionStore

	// Boolean of the open state
	Open bool
}

// readPump pumps messages from the websocket connection to the hub.
func (c *Connection) readPump() {
	defer func() {
		Hub.unregister <- c
		c.ws.Close()
	}()

	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(readWait))
	c.ws.SetPongHandler(func(string) error {
		c.ws.SetReadDeadline(time.Now().Add(readWait))
		return nil
	})

	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			break
		}

		m := Message{}
		m.Connection = c
		m.Message = message

		// Tell the hub to handle it
		Hub.messages <- m
	}
}

// writePump pumps messages from the hub to the websocket connection.
func (c *Connection) writePump() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// write writes a message with the given opCode and payload.
func (c *Connection) write(opCode int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(opCode, payload)
}

func (c *Connection) AppendChannel(channel chan struct{}) {
	c.customChannels = append(c.customChannels, channel)
}

func (c *Connection) CloseCustomChannels() {
	c.Open = false

	for _, c := range c.customChannels {
		close(c)
	}
}
