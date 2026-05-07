package dispatch

import (
	"encoding/json"
	"mangahub/cmd/udp-server/middleware"
	"mangahub/pkg/logger"
	"mangahub/pkg/types"
	"net"
)

type HandleFunc func(s *UDPServer, clientAddr *net.UDPAddr, payload types.UDPMessage)

type UDPServer struct {
	Conn     *net.UDPConn
	handlers map[string]HandleFunc
}

func NewUDPServer() *UDPServer {
	return &UDPServer{
		handlers:   make(map[string]HandleFunc),
	}
}

func (s *UDPServer) Start(port string) {
	// 1. Resolve UDP sender address
	addr, _ := net.ResolveUDPAddr("udp", port)

	// 2. Open UDP port
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		logger.Error("Start UDP Server error", "error", err)
		return
	}
	s.Conn = conn
	defer s.Conn.Close()
	logger.Info("UDP Server listening", "port", port)

	// 3. Handle incoming UDP messages
	buffer := make([]byte, 65535) // Max UDP packet size
	for {
		// Read UDP message
		n, clientAddr, err := s.Conn.ReadFromUDP(buffer)
		if err != nil {
			logger.Error("Error reading UDP message", "error", err)
			continue
		}

		message := make([]byte, n)
		copy(message, buffer[:n])
		logger.Debug("UDP received", "from", clientAddr.String(), "payload", string(message))
		
		// Process UDP message in a new goroutine
		go func (s *UDPServer, message []byte, clientAddr *net.UDPAddr) {
			// Parse UDP message
			var udpMsg types.UDPMessage
			err := json.Unmarshal(message, &udpMsg)
			if err != nil {
				logger.Error("Invalid UDP message format", "error", err)
				return
			}

			// check secret
			if err := middleware.AuthMiddleware(udpMsg.Action, udpMsg.Token); err != nil {
				logger.Error("Unauthorized UDP message", "action", udpMsg.Action, "error", err)
				return
			}

			// Dispatch to handler
			s.Dispatch(clientAddr, udpMsg)
		} (s, message, clientAddr)
	}
}

func (s *UDPServer) RegisterHandler(action string, handler HandleFunc) {
	s.handlers[action] = handler
}

func (s *UDPServer) Dispatch(clientAddr *net.UDPAddr, payload types.UDPMessage) {
	handler, exists := s.handlers[payload.Action]
	if !exists {
		logger.Warn("No handler registered for action", "action", payload.Action)
		return
	}

	handler(s, clientAddr, payload)
}