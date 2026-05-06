package dispatch

import (
	"encoding/json"
	"log"
	"mangahub/cmd/udp-server/middleware"
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
		log.Fatalf("Start UDP Server error: %v", err)
	}
	s.Conn = conn
	defer s.Conn.Close()
	log.Printf("UDP Server listening on port %s", port)

	// 3. Handle incoming UDP messages
	buffer := make([]byte, 65535) // Max UDP packet size
	for {
		// Read UDP message
		n, clientAddr, err := s.Conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("Error reading UDP message: %v", err)
			continue
		}

		// Check if it's from API server (for now, assume any payload with 'broadcast' is trusted)
		// Or the client register packet.
		message := make([]byte, n)
		copy(message, buffer[:n])
		log.Printf("UDP received from %s: %s", clientAddr.String(), message)
		
		// Process UDP message in a new goroutine
		go func (s *UDPServer, message []byte, clientAddr *net.UDPAddr) {
			// Parse UDP message
			var udpMsg types.UDPMessage
			err := json.Unmarshal(message, &udpMsg)
			if err != nil {
				log.Printf("Invalid UDP message format: %v", err)
				return
			}

			// check secret
			if err := middleware.AuthMiddleware(udpMsg.Action, udpMsg.Token); err != nil {
				log.Printf("Unauthorized UDP message: %v", err)
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
		log.Printf("No handler registered for action: %s", payload.Action)
		return
	}

	handler(s, clientAddr, payload)
}
