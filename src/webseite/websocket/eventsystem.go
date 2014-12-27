package websocket

type MessageAccepter func([]byte, *Connection) bool

type EventListener struct {
	incoming        chan *Message
	messageAccepter MessageAccepter
}

type EventSystem struct {
	eventListener []*EventListener
}

func (es *EventSystem) createListener(fn MessageAccepter) *EventListener {
	c := make(chan *Message, 5)

	listener := new(EventListener)
	listener.incoming = c
	listener.messageAccepter = fn

	return listener
}

func (es *EventSystem) Listen(fn MessageAccepter) <-chan *Message {
	var listener *EventListener
	listener = es.createListener(fn)
	es.eventListener = append(es.eventListener, listener)
	return listener.incoming
}

func (es *EventSystem) Emit(message *Message) {
	for _, c := range es.eventListener {
		if c.messageAccepter(message.message, message.connection) {
			c.incoming <- message
		}
	}
}

func (es *EventSystem) Close() {
	for _, c := range es.eventListener {
		close(c.incoming)
	}
}
