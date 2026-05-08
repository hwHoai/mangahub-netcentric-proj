package udp_services_impl

import (
	"encoding/json"
	"fmt"
	udp_services "mangahub/internal/udp"
	"mangahub/pkg/types"
	"net"
	"os"
)

type UDPKeySyncServiceImpl struct {
	addr string
}

var _ udp_services.UDPKeySyncServices = (*UDPKeySyncServiceImpl)(nil)

func NewUDPKeySyncService(addr string) udp_services.UDPKeySyncServices {
	return &UDPKeySyncServiceImpl{addr: addr}
}

func (c *UDPKeySyncServiceImpl) SyncPublicKey(publicKeyPEM string) error {
	serverAddr, err := net.ResolveUDPAddr("udp", c.addr)
	if err != nil {
		return fmt.Errorf("failed to resolve UDP address: %w", err)
	}

	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		return fmt.Errorf("failed to dial UDP server: %w", err)
	}
	defer conn.Close()

	payload, _ := json.Marshal(map[string]string{
		"public_key": publicKeyPEM,
	})

	msg := types.UDPMessage{
		Action:  "pub_key:impl_sync_public_key",
		Payload: payload,
		Token:   os.Getenv("HANDSHAKE_KEY"),
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal UDP message: %v", err)
	}

	_, err = conn.Write(append(data, '\n'))
	if err != nil {
		return fmt.Errorf("failed to send UDP message: %v", err)
	}

	return nil
}
