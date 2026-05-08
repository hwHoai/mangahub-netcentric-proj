package udp_services_impl

import (
	"encoding/json"
	"fmt"
	"net"

	udp_services "mangahub/internal/udp"
	"mangahub/pkg/types"
)
type UDPNotificationServicesImpl struct {
	serverAddr   *net.UDPAddr
	conn         *net.UDPConn
	handshakeKey string
}

var _ udp_services.UDPChapterNotificationServices = (*UDPNotificationServicesImpl)(nil)

// NewNotificationServices creates a new instance of UDPNotificationServices.
func NewNotificationServices(serverAddr *net.UDPAddr, conn *net.UDPConn, handshakeKey string) udp_services.UDPChapterNotificationServices {
	return &UDPNotificationServicesImpl{
		serverAddr:   serverAddr,
		conn:         conn,
		handshakeKey: handshakeKey,
	}
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
		Token:   c.handshakeKey,
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
		Token:   c.handshakeKey,
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
