package udp_services_impl

import (
	"encoding/json"
	"net"
	"os"
	"strings"
	"testing"
	"time"

	"mangahub/pkg/types"

	"github.com/stretchr/testify/assert"
)

func TestUDPNotificationServicesImpl_SendNotifications(t *testing.T) {
	// 1. Start a mock UDP server
	serverAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	assert.NoError(t, err)

	conn, err := net.ListenUDP("udp", serverAddr)
	assert.NoError(t, err)
	defer conn.Close()

	// Extract the random port assigned
	addrStr := conn.LocalAddr().String()
	parts := strings.Split(addrStr, ":")
	port := parts[len(parts)-1]

	os.Setenv("UDP_SERVER_HOST", "127.0.0.1")
	os.Setenv("UDP_SERVER_PORT", port)
	os.Setenv("HANDSHAKE_KEY", "test-secret")

	// Setup receiver channel
	done := make(chan types.UDPMessage, 2)

	go func() {
		buf := make([]byte, 1024)
		for i := 0; i < 2; i++ {
			n, _, err := conn.ReadFromUDP(buf)
			if err != nil {
				return
			}
			var msg types.UDPMessage
			json.Unmarshal(buf[:n], &msg)
			done <- msg
		}
	}()

	// 2. Initialize the client
	remoteAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:"+port)
	assert.NoError(t, err)
	clientConn, err := net.DialUDP("udp", nil, remoteAddr)
	assert.NoError(t, err)

	client := NewNotificationServices(remoteAddr, clientConn, "test-secret")
	defer client.Close()

	// 3. Test Chapter Notification
	err = client.SendNewChapterNotification("manga-1", "chapter-1", "New Title", 1.0)
	assert.NoError(t, err)

	// 4. Test Message Notification
	err = client.SendNewMessageNotification("room-1", "user-1", "Hello world")
	assert.NoError(t, err)

	// 5. Verify the received messages
	var received []types.UDPMessage
	for i := 0; i < 2; i++ {
		select {
		case msg := <-done:
			received = append(received, msg)
		case <-time.After(2 * time.Second):
			t.Fatal("Timeout waiting for UDP messages")
		}
	}

	assert.Len(t, received, 2)
	assert.Equal(t, "chapter:impl_broadcast_notification", received[0].Action)
	assert.Equal(t, "chat:impl_broadcast_message", received[1].Action)
}
