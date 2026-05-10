package handler

import (
	"crypto/sha256"
	"encoding/json"
	"mangahub/cmd/udp-server/dispatch"
	udp_pools "mangahub/cmd/udp-server/utils/pool"
	"mangahub/pkg/logger"
	"mangahub/pkg/types"
	"net"
	"sort"
)

type BenchmarkHandler struct {
	pool udp_pools.UDPPool
}

func NewBenchmarkHandler(pool udp_pools.UDPPool) *BenchmarkHandler {
	return &BenchmarkHandler{pool: pool}
}

func (h *BenchmarkHandler) PingHandler(s *dispatch.UDPServer, addr *net.UDPAddr, msg types.UDPMessage) {
	
	_, port, _ := net.SplitHostPort(addr.String())
	h.pool.Register(port, addr)

	response := types.UDPMessage{
		Action:  "benchmark:res_register",
		Payload: json.RawMessage(`{"status": "ok"}`),
	}
	data, _ := json.Marshal(response)
	s.Conn.WriteToUDP(data, addr)
}

func (h *BenchmarkHandler) AckHandler(s *dispatch.UDPServer, addr *net.UDPAddr, msg types.UDPMessage) {
	var data struct {
		NotificationID string `json:"notification_id"`
	}
	json.Unmarshal(msg.Payload, &data)

	// Sử dụng Port làm ID duy nhất trên localhost để tránh mismatch IPv4/IPv6 string
	_, port, _ := net.SplitHostPort(addr.String())
	h.pool.ProcessAck(data.NotificationID, port)
}

func (h *BenchmarkHandler) BroadcastHandler(s *dispatch.UDPServer, addr *net.UDPAddr, msg types.UDPMessage) {
	// --- EXTREME DATA PREPARATION (Before Broadcast) ---
	tokenSecret := []byte("manga_hub_vanguard_broadcast_prep_2026")
	for i := 0; i < 20000; i++ {
		hash := sha256.Sum256(tokenSecret)
		tokenSecret = hash[:]
	}

	priorityList := make([]int, 2000)
	for i := 0; i < 2000; i++ {
		priorityList[i] = 2000 - i
	}
	sort.Ints(priorityList)

	_ = make([]byte, 512*1024)
	// ----------------------------------------------------

	logger.Info("🚀 UDP Broadcast Triggered (Heavy Prep Done)", "from", addr.String())

	// Mock Chapter Notification
	payload := map[string]interface{}{
		"action": "chapter:on_new_chapter_notification",
		"payload": types.NewChapterNotificationPayload{
			ID:            "notif_benchmark_001",
			MangaID:       "manga_123",
			ChapterID:     "chap_99",
			ChapterTitle:  "The Final Refactor",
			ChapterNumber: 99.0,
		},
	}

	h.pool.Broadcast(s.Conn, "notif_001", payload)
	
}
