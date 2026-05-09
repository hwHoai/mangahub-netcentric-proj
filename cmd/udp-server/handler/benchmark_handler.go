package handler

import (
	"encoding/json"
	"mangahub/cmd/udp-server/dispatch"
	udp_pools "mangahub/cmd/udp-server/utils/pool"
	"mangahub/pkg/logger"
	"mangahub/pkg/types"
	"net"
	"time"
)

type BenchmarkHandler struct {
	pool udp_pools.UDPPool
}

func NewBenchmarkHandler(pool udp_pools.UDPPool) *BenchmarkHandler {
	return &BenchmarkHandler{pool: pool}
}

func (h *BenchmarkHandler) PingHandler(s *dispatch.UDPServer, addr *net.UDPAddr, msg types.UDPMessage) {
	// Simulate heavy work (e.g., Database Query, Image processing)
	time.Sleep(2000 * time.Millisecond)

	// Register client in pool
	h.pool.Register(addr.String(), addr)

	response := types.UDPMessage{
		Action:  "benchmark:res_register",
		Payload: json.RawMessage(`{"status": "ok", "msg": "REGISTER_SUCCESS"}`),
	}
	data, _ := json.Marshal(response)
	
	// Send registration confirmation
	s.Conn.WriteToUDP(append(data, '\n'), addr)
	
	logger.Info("UDP Subscriber registered (MOCK)", "addr", addr.String())
}

func (h *BenchmarkHandler) BroadcastHandler(s *dispatch.UDPServer, addr *net.UDPAddr, msg types.UDPMessage) {
	// Trigger a broadcast to all registered benchmark clients
	logger.Info("UDP Broadcast trigger received", "from", addr.String())

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

	h.pool.Broadcast(s.Conn, "all", payload)
}
