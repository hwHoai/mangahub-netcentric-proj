package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/google/uuid"
	"mangahub/pkg/types"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	target := flag.String("target", "localhost:8082", "Target TCP server address")
	conns := flag.Int("conns", 1000, "Number of concurrent connections to spawn")
	flag.Parse()

	fmt.Printf("🚀 Starting TCP Stress Test: %d connections to %s\n", *conns, *target)

	var wg sync.WaitGroup
	var activeConns int64
	var totalPongsReceived int64
	phase2Start := make(chan struct{})

	start := time.Now()

	mockChapter := fmt.Sprintf("manga_id:%s,chapter_id:%s", uuid.New().String(), uuid.New().String())

	fmt.Printf("Step 1: Registering %d subscribers...\n", *conns)
	for i := 0; i < *conns; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// Determine UserID (simulating 3 conns per user)
			userID := fmt.Sprintf("benchmark_user_%d", id/3)

			conn, err := net.DialTimeout("tcp", *target, 5*time.Second)
			if err != nil {
				return
			}
			defer conn.Close()

			// 1. REGISTER PHASE
			regMsg := types.TCPMessage{
				Action:  "benchmark:test_register",
				Payload: json.RawMessage(fmt.Sprintf(`{"user_id": "%s"}`, userID)),
			}
			data, _ := json.Marshal(regMsg)
			fmt.Fprintln(conn, string(data))

			// Wait for Registration Success
			scanner := bufio.NewScanner(conn)
			if scanner.Scan() {
				atomic.AddInt64(&activeConns, 1)
			} else {
				return
			}

			// Wait for the Burst Command
			<-phase2Start

			// BURST START
			
			// Send ONE heavy request
			payloadData, _ := json.Marshal(map[string]any{
				"user_id":      userID,
				"mock_chapter": mockChapter,
			})
			ping := types.TCPMessage{
				Action:  "benchmark:test_ping",
				Payload: json.RawMessage(payloadData),
			}
			pingData, _ := json.Marshal(ping)
			fmt.Fprintln(conn, string(pingData))
			
			// Wait for the response
			if scanner.Scan() {
				atomic.AddInt64(&totalPongsReceived, 1)
			}
			
			// Báo cáo thời gian hoàn thành của riêng connection này (latency)
			// Tuy nhiên ta sẽ lấy tổng kết ở main
		}(i)

		if i > 0 && i%500 == 0 {
			fmt.Printf("... %d connections attempted\n", i)
		}
		time.Sleep(333 * time.Microsecond)
	}

	// Wait for all registrations to be attempted
	for {
		if atomic.LoadInt64(&activeConns) >= int64(*conns) {
			break
		}
		if time.Since(start) > 30*time.Second {
			fmt.Printf("Timeout waiting for registrations. Active: %d/%d\n", atomic.LoadInt64(&activeConns), *conns)
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	regDuration := time.Since(start)
	fmt.Printf("✅ Registration Phase completed in %v. Active: %d/%d\n", regDuration, atomic.LoadInt64(&activeConns), *conns)

	fmt.Println("\n🚀 TRIGGERING SIMULTANEOUS BURST...")
	startBurst := time.Now()
	close(phase2Start) // Release the barrier

	wg.Wait()
	burstDuration := time.Since(startBurst)

	received := atomic.LoadInt64(&totalPongsReceived)
	fmt.Printf("✅ Burst Test finished.\n")
	fmt.Printf("----------------------------------\n")
	fmt.Printf("1. Registration Phase : %v\n", regDuration)
	fmt.Printf("2. Burst Load Phase   : %v\n", burstDuration)
	fmt.Printf("3. Delivery Report (Client-side):\n")
	fmt.Printf("   - Active Sessions  : %d\n", atomic.LoadInt64(&activeConns))
	fmt.Printf("   - Total Syncs Received: %d\n", received)
	fmt.Printf("----------------------------------\n")
}
