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

type TCPChapterSyncServiceImpl struct {
	addr string
}

var _ tcp_services.TCPChapterSyncServices = (*TCPChapterSyncServiceImpl)(nil)

func NewTCPChapterSyncService(addr string) tcp_services.TCPChapterSyncServices {
	return &TCPChapterSyncServiceImpl{addr: addr}
}

func (c *TCPChapterSyncServiceImpl) SyncReading(userID string, chapterID string) error {
	fmt.Println("SyncReading", userID, chapterID)
	conn, err := net.DialTimeout("tcp", c.addr, 2*time.Second)
	if err != nil {
		return fmt.Errorf("failed to connect to TCP server: %v", err)
	}
	defer conn.Close()

	payload, _ := json.Marshal(map[string]string{
		"user_id":    userID,
		"chapter_id": chapterID,
	})

	msg := types.TCPMessage{
		Action:  "chapter_sync:impl_broadcast_read",
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
