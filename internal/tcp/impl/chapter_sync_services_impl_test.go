package tcp_services_impl

import (
	"bufio"
	"encoding/json"
	"net"
	"os"
	"testing"
	"time"

	"mangahub/pkg/types"

	"github.com/stretchr/testify/assert"
)

func TestSyncReading_Success(t *testing.T) {
	// 1. Setup mock TCP server
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	assert.NoError(t, err)
	defer listener.Close()

	addr := listener.Addr().String()

	// Wait for connection in a goroutine
	done := make(chan bool)
	go func() {
		conn, err := listener.Accept()
		assert.NoError(t, err)
		defer conn.Close()

		// Read the message
		scanner := bufio.NewScanner(conn)
		if scanner.Scan() {
			var msg types.TCPMessage
			err := json.Unmarshal(scanner.Bytes(), &msg)
			assert.NoError(t, err)
			assert.Equal(t, "chapter_sync:impl_broadcast_read", msg.Action)

			// Parse payload
			var payload map[string]string
			json.Unmarshal(msg.Payload, &payload)
			assert.Equal(t, "user-1", payload["user_id"])
			assert.Equal(t, "chapter-1", payload["chapter_id"])
		}
		done <- true
	}()

	// 2. Set environment variable for Token
	os.Setenv("HANDSHAKE_KEY", "test-key")

	// 3. Test the client
	service := NewTCPChapterSyncService(addr)
	err = service.SyncReading("user-1", "chapter-1")

	assert.NoError(t, err)

	// Wait for server to process
	select {
	case <-done:
	case <-time.After(1 * time.Second):
		t.Fatal("Test timed out waiting for TCP server")
	}
}

func TestSyncReading_ConnectionFailure(t *testing.T) {
	// Use a non-existent port to force connection failure
	service := NewTCPChapterSyncService("127.0.0.1:12345")
	err := service.SyncReading("user-1", "chapter-1")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to connect to TCP server")
}
