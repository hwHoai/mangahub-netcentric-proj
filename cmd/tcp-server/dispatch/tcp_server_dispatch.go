package dispatch

import (
	"mangahub/cmd/tcp-server/middleware"
	"mangahub/pkg/logger"
	"mangahub/pkg/types"
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

func (d *Dispatcher) Dispatch(conn net.Conn, msg types.TCPMessage) {
	action := msg.Action
	token := msg.Token

	// Run authentication middleware
	if err := middleware.AuthMiddleware(action, token); err != nil {
		logger.Error("Security Block", "error", err, "action", action)
		return
	}

	handler, exists := d.handlers[action]
	if !exists {
		logger.Warn("No handler registered for action", "action", action)
		return
	}

	handler(conn, msg.Payload)
}
