package types

import "encoding/json"

type UDPMessage struct {
	Action  string          `json:"action"`
	Token   string          `json:"token"`
	Payload json.RawMessage `json:"payload"`
}

type NewChapterNotificationPayload struct {
	ID            string  `json:"id,omitempty"` // Unique ID for ACK
	MangaID       string  `json:"manga_id"`
	ChapterID     string  `json:"chapter_id"`
	ChapterTitle  string  `json:"chapter_title"`
	ChapterNumber float64 `json:"chapter_number"`
}

type NewMessageNotificationPayload struct {
	ID         string `json:"id,omitempty"` // Unique ID for ACK
	RoomID     string `json:"room_id"`
	SenderName string `json:"sender_name"`
	Content    string `json:"content"`
}

type NotificationAckPayload struct {
	NotificationID string `json:"notification_id"`
	UserID         string `json:"user_id"`
}