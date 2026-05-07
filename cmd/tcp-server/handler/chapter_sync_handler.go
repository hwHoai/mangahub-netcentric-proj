package handler

import (
	"encoding/json"
	"fmt"
	pool_impl "mangahub/cmd/tcp-server/utils/pool/impl"
	"net"
)

type ChapterSyncHandler struct {
	pool *pool_impl.ChapterSyncPool
}

func NewChapterSyncHandler(pool *pool_impl.ChapterSyncPool) *ChapterSyncHandler {
	return &ChapterSyncHandler{pool: pool}
}

func (h *ChapterSyncHandler) RegisterConnectionHandler(conn net.Conn, payload any) {
	var data struct {
		UserID string `json:"user_id"`
	}
	
	raw, ok := payload.(json.RawMessage)
	if !ok {
		fmt.Println("Invalid payload type for register")
		return
	}

	if err := json.Unmarshal(raw, &data); err != nil {
		fmt.Printf("Error unmarshaling register payload: %v\n", err)
		return
	}

	h.pool.Register(data.UserID, conn)

	newProgressMsg, _ := json.Marshal(map[string]any{
		"action": "chapter_sync:on_new_read_progress",
		"payload": map[string]string{
			"chapter_id": h.pool.GetLastChapter(data.UserID),
		},
	})

	h.pool.BroadcastOne(data.UserID, newProgressMsg, conn)
	fmt.Printf("User %s registered connection\n", data.UserID)
	conn.Write([]byte("Registered successfully\n"))
}

func (h *ChapterSyncHandler) BroadcastReadHandler(conn net.Conn, payload any) {
	var data struct {
		UserID    string `json:"user_id"`
		ChapterID string `json:"chapter_id"`
	}

	raw, ok := payload.(json.RawMessage)
	if !ok {
		fmt.Println("Invalid payload type for broadcast_read")
		return
	}

	if err := json.Unmarshal(raw, &data); err != nil {
		fmt.Printf("Error unmarshaling broadcast_read payload: %v\n", err)
		return
	}

	fmt.Printf("Broadcasting sync_reading for user %s, chapter %s\n", data.UserID, data.ChapterID)
	
	h.pool.UpdateLastChapter(data.UserID, data.ChapterID)

	syncMsg, _ := json.Marshal(map[string]any{
		"action": "chapter_sync:on_new_read_progress",
		"payload": map[string]string{
			"chapter_id": data.ChapterID,
		},
	})
	
	h.pool.Broadcast(data.UserID, syncMsg)
}
