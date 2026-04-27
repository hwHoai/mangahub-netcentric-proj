package handler

import (
	"fmt"
	"mangahub/cmd/udp-server/dispatch"
	"mangahub/pkg/types"
)

func ExampleUDPHandler(s *dispatch.UDPServer, payload types.UDPMessage) {
	fmt.Printf("Request handled: %v\n", payload)
	response := []byte(`{"status": "Request handled successfully"}\n`)
	clientAddr, exists := s.ClientsMap[payload.UserID]
	if !exists {
		fmt.Printf("Client address not found for user ID: %s\n", payload.UserID)
		return
	}
	s.Conn.WriteToUDP(response, clientAddr)
}