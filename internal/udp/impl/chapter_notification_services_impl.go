package udp_services_impl

import (
	"encoding/json"
	"fmt"
	"net"
	"os"

	udp_services "mangahub/internal/udp"
	"mangahub/pkg/types"
)
type UDPNotificationServicesImpl struct {
	serverAddr *net.UDPAddr
	conn       *net.UDPConn
}

var _ udp_services.UDPChapterNotificationServices = (*UDPNotificationServicesImpl)(nil)

// NewNotificationClient creates a new instance of NotificationClient.
func NewNotificationServicesImpl() (udp_services.UDPChapterNotificationServices, error) {
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

	return &UDPNotificationServicesImpl{
		serverAddr: serverAddr,
		conn:       conn,
	}, nil
}

func (c *UDPNotificationServicesImpl) SendNewChapterNotification(mangaID, chapterID, title string, chapterNumber float64) error {
	payload := types.NewChapterNotificationPayload{
		MangaID:       mangaID,
		ChapterID:     chapterID,
		ChapterTitle:  title,
		ChapterNumber: chapterNumber,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal inner payload: %w", err)
	}

	outerMsg := types.UDPMessage{
		Action:  "chapter:impl_broadcast_notification",
		Payload: payloadBytes,
		Token:   os.Getenv("HANDSHAKE_KEY"),
	}

	data, err := json.Marshal(outerMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal UDP message: %w", err)
	}

	data = append(data, '\n')

	_, err = c.conn.Write(data)
	if err != nil {
		return fmt.Errorf("failed to send UDP message: %w", err)
	}

	return nil
}

func (c *UDPNotificationServicesImpl) SendNewMessageNotification(roomID, senderName, content string) error {
	payload := types.NewMessageNotificationPayload{
		RoomID:     roomID,
		SenderName: senderName,
		Content:    content,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal inner payload: %w", err)
	}

	outerMsg := types.UDPMessage{
		Action:  "chat:impl_broadcast_message",
		Payload: payloadBytes,
		Token:   os.Getenv("HANDSHAKE_KEY"),
	}

	data, err := json.Marshal(outerMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal UDP message: %w", err)
	}

	data = append(data, '\n')

	_, err = c.conn.Write(data)
	if err != nil {
		return fmt.Errorf("failed to send UDP message: %w", err)
	}

	return nil
}

func (c *UDPNotificationServicesImpl) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
