package dispatch

import (
	"encoding/json"
	"log"
	"mangahub/pkg/types"
	"net"
	"sync"
)

type HandleFunc func(s *UDPServer, payload types.UDPMessage)

type UDPServer struct {
    Conn       *net.UDPConn
    ClientsMap map[string]*net.UDPAddr
    Mutex      sync.RWMutex          
	handlers   map[string]HandleFunc

}

func NewUDPServer() *UDPServer {
	return &UDPServer{
		ClientsMap: make(map[string]*net.UDPAddr),
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
	buffer := make([]byte, 1024)
	for {
		// Read UDP message
		n, clientAddr, err := s.Conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("Error reading UDP message: %v", err)
			continue
		}

		// Register client address
		s.RegisterClient()

		// copy message to avoid data race
		message := make([]byte, n)
		copy(message, buffer[:n])
		log.Printf("UDP received from %s: %s", clientAddr.String(), message)
		
		// Process UDP message in a new goroutine
		go func (s *UDPServer, message []byte) {
			// Parse UDP message
			var udpMsg types.UDPMessage
			err := json.Unmarshal(message, &udpMsg)
			if err != nil {
				log.Printf("Invalid UDP message format: %v", err)
				return
			}

			// check secret
			if false { // TODO: check secret
				log.Printf("Unauthorized UDP message: invalid secret")
				return
			}

			// Dispatch to handler
			s.Dispatch(udpMsg)
		} (s, message)
	}
}

func (s *UDPServer) RegisterHandler(action string, handler HandleFunc) {
	s.handlers[action] = handler
}

func (s *UDPServer) Dispatch(payload types.UDPMessage) {
	handler, exists := s.handlers[payload.Action]
	if !exists {
		log.Printf("No handler registered for action: %s", payload.Action)
		return
	}

	handler(s, payload)
}

func (s *UDPServer) RegisterClient() {
	// TODO: implement client registration
	// E.g. When user allow notification, 
	// 		db store user_id and which notification they want to receive,
	// 		then when app start, it take user_id and udp address to register to server,
	// 		when server want to send notification, 
	// 		it look up user_id to get udp address and send notification to that address

	sampleClientAddr, _ := net.ResolveUDPAddr("udp", "localhost:12345")
	s.Mutex.Lock()
	s.ClientsMap["sample_user_id"] = sampleClientAddr
	s.Mutex.Unlock()
}