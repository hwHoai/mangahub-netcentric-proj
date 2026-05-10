package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"mangahub/pkg/types"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	target := flag.String("target", "localhost:8083", "Target UDP server address")
	clientCount := flag.Int("clients", 1000, "Number of concurrent clients to spawn")
	flag.Parse()

	serverAddr, err := net.ResolveUDPAddr("udp", *target)
	if err != nil {
		fmt.Printf("❌ Failed to resolve address: %v\n", err)
		return
	}

	fmt.Printf("🚀 [UDP-BROADCAST] Starting Broadcast Stress Test\n")
	fmt.Printf("🚀 Target Server: %s | Clients: %d\n", *target, *clientCount)

	var totalPongsReceived int64
	var totalBroadcastReceived int64
	var wg sync.WaitGroup

	fmt.Printf("Step 1: Registering %d subscribers...\n", *clientCount)
	for i := 0; i < *clientCount; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			conn, err := net.DialUDP("udp", nil, serverAddr)
			if err != nil {
				return
			}
			defer conn.Close()

			// Dùng channel để báo hiệu khi nhận được broadcast
			finishChan := make(chan bool)

			// Goroutine lắng nghe
			go func(c *net.UDPConn) {
				buf := make([]byte, 2048)
				for {
					c.SetReadDeadline(time.Now().Add(25 * time.Second))
					n, err := c.Read(buf) // Dùng Read cho socket đã Dial
					if err != nil {
						finishChan <- false
						return
					}

					var msg types.UDPMessage
					if err := json.Unmarshal(buf[:n], &msg); err != nil {
						continue
					}

					if msg.Action == "benchmark:res_register" {
						atomic.AddInt64(&totalPongsReceived, 1)
						continue 
					}

					if msg.Action == "benchmark:test_broadcast" {
						atomic.AddInt64(&totalBroadcastReceived, 1)
						
						ackMsg := types.UDPMessage{
							Action: "benchmark:test_ack",
							Payload: json.RawMessage(fmt.Sprintf(`{"notification_id": "notif_001", "user_id": "%s"}`, c.LocalAddr().String())),
						}
						ackData, _ := json.Marshal(ackMsg)
						c.Write(ackData) // Dùng Write vì socket đã Dial
						finishChan <- true
						return
					}
				}
			}(conn)

			// Gửi gói tin đăng ký
			regMsg := types.UDPMessage{
				Action:  "benchmark:test_register",
				Payload: json.RawMessage(fmt.Sprintf(`{"user_id": "user_%d"}`, id/3)),
			}
			regData, _ := json.Marshal(regMsg)
			conn.Write(regData)

			// Đợi cho đến khi nhận được broadcast hoặc timeout
			<-finishChan
		}(i)

		if i > 0 && i%1000 == 0 {
			fmt.Printf("... %d registered\n", i)
		}
		time.Sleep(100 * time.Microsecond) // Avoid OS buffer overflow during registration
	}

	time.Sleep(3 * time.Second) // Wait for all registrations to settle
	fmt.Printf("\n✅ Step 1: All clients registered. Total Responses (Ack): %d\n", atomic.LoadInt64(&totalPongsReceived))

	// Reset counter for broadcast
	atomic.StoreInt64(&totalPongsReceived, 0)

	fmt.Printf("\nStep 2: Triggering Server-side Broadcast (Heavy Prep)...\n")
	controller, _ := net.DialUDP("udp", nil, serverAddr)
	triggerMsg := types.UDPMessage{
		Action:  "benchmark:test_trigger_broadcast",
		Payload: json.RawMessage(`{"cmd": "START_BROADCAST"}`),
	}
	data, _ := json.Marshal(triggerMsg)
	
	startBroadcast := time.Now()
	controller.Write(append(data, '\n'))
	dispatchDuration := time.Since(startBroadcast)

	fmt.Println("⏳ Waiting 15 seconds for broadcast delivery & ACKs...")
	time.Sleep(15 * time.Second)
	
	received := atomic.LoadInt64(&totalBroadcastReceived)
	fmt.Printf("\n✅ UDP Broadcast Stress Test finished.\n")
	fmt.Printf("----------------------------------\n")
	fmt.Printf("1. Prep & Dispatch Duration : %v\n", dispatchDuration)
	fmt.Printf("2. Delivery Report (Client-side):\n")
	fmt.Printf("   - Total Target Clients   : %d\n", *clientCount)
	fmt.Printf("   - Total Received         : %d\n", received)
	fmt.Printf("   - Delivery Rate          : %.2f%%\n", float64(received)/float64(*clientCount)*100)
	fmt.Printf("----------------------------------\n")
}
