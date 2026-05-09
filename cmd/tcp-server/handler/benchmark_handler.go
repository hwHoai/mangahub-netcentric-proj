package handler

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	pool_impl "mangahub/cmd/tcp-server/utils/pool/impl"
	"mangahub/pkg/logger"
	"mangahub/pkg/types"
	"net"
	"sort"
)

type BenchmarkHandler struct {
	pool *pool_impl.BenchmarkPool
}

func NewBenchmarkHandler(pool *pool_impl.BenchmarkPool) *BenchmarkHandler {
	return &BenchmarkHandler{pool: pool}
}

func (h *BenchmarkHandler) RegisterHandler(conn net.Conn, payload any) {
	var data struct {
		UserID string `json:"user_id"`
	}

	raw, ok := payload.(json.RawMessage)
	if !ok {
		logger.Error("Invalid payload type for benchmark:register")
		return
	}

	if err := json.Unmarshal(raw, &data); err != nil {
		logger.Error("Error unmarshaling benchmark:register payload", "error", err)
		return
	}

	h.pool.Register(data.UserID, conn)

	response := types.TCPMessage{
		Action: "benchmark:res_register",
		Payload: json.RawMessage(`{"status": "ok"}`),
	}
	dataBytes, _ := json.Marshal(response)
	conn.Write(append(dataBytes, '\n'))
	logger.Info("Benchmark client registered", "userID", data.UserID, "addr", conn.RemoteAddr().String())
}

func (h *BenchmarkHandler) PingHandler(conn net.Conn, payload any) {
	// --- EXTREME REAL-WORLD HEAVY LOAD SIMULATION ---
	
	// 1. Extreme Security (Simulating ultra-secure hashing / Argon2 equivalent)
	tokenSecret := []byte("manga_hub_vanguard_secure_key_2026_extreme")
	for i := 0; i < 20000; i++ {
		hash := sha256.Sum256(tokenSecret)
		tokenSecret = hash[:]
	}

	// 2. Complex Data Sorting (Simulating Large-scale Subscriber Ranking)
	priorityList := make([]int, 2000)
	for i := 0; i < 2000; i++ {
		priorityList[i] = 2000 - i
	}
	sort.Ints(priorityList)

	// 3. Massive Memory Pressure (8GB Churn for 16k requests)
	buffer := make([]byte, 512*1024)
	for i := 0; i < 100; i++ {
		buffer[i] = byte(i)
	}

	// --- END OF EXTREME SIMULATION ---

	var data struct {
		UserID      string `json:"user_id"`
		MockChapter string `json:"mock_chapter"`
	}

	raw, ok := payload.(json.RawMessage)
	if !ok {
		logger.Error("Invalid payload type for benchmark:test_ping")
		return
	}

	if err := json.Unmarshal(raw, &data); err != nil {
		logger.Error("Error unmarshaling benchmark:test_ping payload", "error", err, "raw", string(raw))
		return
	}

	// Broadcast pong with mock chapter to all connections of the user
	response := types.TCPMessage{
		Action: "benchmark:res_pong",
		Payload: json.RawMessage(fmt.Sprintf(`{"user_id": "%s", "mock_chapter": "%s"}`, data.UserID, data.MockChapter)),
	}
	dataBytes, _ := json.Marshal(response)
	h.pool.Broadcast(data.UserID, dataBytes)

	// logger.Info("Benchmark broadcasted pong", "userID", data.UserID, "chapter", data.MockChapter)
}

func (h *BenchmarkHandler) PongHandler(conn net.Conn, payload any) {
	// Log received pong from client if any
	logger.Info("Received PONG from client (Benchmark)", "addr", conn.RemoteAddr().String())
}
