package dispatch

import (
	"log"
	"net"
)

type HandleFunc func(conn net.Conn, payload any)

type Dispatcher struct {
	handlers map[string]HandleFunc
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		handlers: make(map[string]HandleFunc),
	}
}

func (d *Dispatcher) RegisterHandler(action string, handler HandleFunc) {
	d.handlers[action] = handler
}

func (d *Dispatcher) Dispatch(conn net.Conn, action string, payload any) {
	handler, exists := d.handlers[action]
	if !exists {
		log.Printf("No handler registered for action: %s", action)
		return
	}

	handler(conn, payload)
}

