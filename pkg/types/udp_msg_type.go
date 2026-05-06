package types

import "encoding/json"

type UDPMessage struct {
	Action  string          `json:"action"`
	Token   string          `json:"token"`
	Payload json.RawMessage `json:"payload"`
}

type NewChapterNotificationPayload struct {
	MangaID       string  `json:"manga_id"`
	ChapterID     string  `json:"chapter_id"`
	ChapterTitle  string  `json:"chapter_title"`
	ChapterNumber float64 `json:"chapter_number"`
}

type NewMessageNotificationPayload struct {
	RoomID     string `json:"room_id"`
	SenderName string `json:"sender_name"`
	Content    string `json:"content"`
}