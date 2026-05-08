package clients

import (
	"fmt"
	"mangahub/internal/udp"
	"mangahub/internal/udp/impl"
	"net"
	"os"
)

func NewUDPNotificationClient() (udp_services.UDPChapterNotificationServices, error) {
	host := os.Getenv("UDP_SERVER_HOST")
	if host == "" {
		host = "127.0.0.1"
	}
	port := os.Getenv("UDP_SERVER_PORT")
	if port == "" {
		port = "8083"
	}

	serverAddr, err := net.ResolveUDPAddr("udp", host+":"+port)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve UDP address: %w", err)
	}

	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to dial UDP server: %w", err)
	}

	handshakeKey := os.Getenv("HANDSHAKE_KEY")

	return udp_services_impl.NewNotificationServices(serverAddr, conn, handshakeKey), nil
}
