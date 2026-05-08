package tcp_services_impl

import (
	"encoding/json"
	"fmt"
	tcp_services "mangahub/internal/tcp"
	"mangahub/pkg/types"
	"net"
	"os"
	"time"
)

type TCPKeySyncServiceImpl struct {
	addr string
}

var _ tcp_services.TCPKeySyncServices = (*TCPKeySyncServiceImpl)(nil)

func NewTCPKeySyncService(addr string) tcp_services.TCPKeySyncServices {
	return &TCPKeySyncServiceImpl{addr: addr}
}

func (c *TCPKeySyncServiceImpl) SyncPublicKey(publicKeyPEM string) error {
	conn, err := net.DialTimeout("tcp", c.addr, 2*time.Second)
	if err != nil {
		return fmt.Errorf("failed to connect to TCP server: %v", err)
	}
	defer conn.Close()

	payload, _ := json.Marshal(map[string]string{
		"public_key": publicKeyPEM,
	})

	msg := types.TCPMessage{
		Action:  "pub_key:impl_sync_public_key",
		Payload: payload,
		Token:   os.Getenv("HANDSHAKE_KEY"),
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal TCP message: %v", err)
	}

	_, err = conn.Write(append(data, '\n'))
	if err != nil {
		return fmt.Errorf("failed to send TCP message: %v", err)
	}

	return nil
}
