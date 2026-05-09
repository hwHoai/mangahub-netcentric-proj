package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"mangahub/pkg/types"
	"net"
	"sync/atomic"
	"time"
)

func main() {
	target := flag.String("target", "localhost:8083", "Target UDP server address")
	clientCount := flag.Int("clients", 2000, "Number of clients to simulate")
	flag.Parse()

	serverAddr, err := net.ResolveUDPAddr("udp", *target)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("🚀 [UDP-BROADCAST] Starting Broadcast Stress Test\n")
	fmt.Printf("🚀 Target Server: %s | Clients: %d\n", *target, *clientCount)

	var clients = make([]*net.UDPConn, *clientCount)

	// Step 1: Register clients
	fmt.Printf("Step 1: Registering %d dummy clients...\n", *clientCount)
	startReg := time.Now()
	var receivedCount int64 // Bộ đếm gói tin nhận được

	for i := 0; i < *clientCount; i++ {
		conn, err := net.ListenUDP("udp", nil)
		if err != nil {
			continue
		}
		clients[i] = conn

		// Mở goroutine lắng nghe cho từng client
		go func(c *net.UDPConn) {
			buf := make([]byte, 1024)
			// Đặt timeout 5 giây, nếu không nhận được gì thì tự đóng
			c.SetReadDeadline(time.Now().Add(5 * time.Second))
			_, _, err := c.ReadFromUDP(buf)
			if err == nil {
				atomic.AddInt64(&receivedCount, 1)
			}
		}(conn)

		pingMsg := types.UDPMessage{
			Action:  "benchmark:test_ping",
			Payload: json.RawMessage(fmt.Sprintf(`{"id": %d}`, i)),
		}
		data, _ := json.Marshal(pingMsg)
		conn.WriteToUDP(data, serverAddr)
		
		if i > 0 && i%1000 == 0 {
			fmt.Printf("... %d clients registered\n", i)
		}
		time.Sleep(100 * time.Microsecond)
	}
	fmt.Printf("✅ Registration completed in %v\n", time.Since(startReg))

	// Step 2: Trigger Broadcast
	fmt.Printf("\nStep 2: Triggering Server-side Broadcast...\n")
	controller, _ := net.DialUDP("udp", nil, serverAddr)
	triggerMsg := types.UDPMessage{
		Action:  "benchmark:trigger_broadcast",
		Payload: json.RawMessage(`{"cmd": "START_BROADCAST"}`),
	}
	triggerData, _ := json.Marshal(triggerMsg)
	
	startBroadcast := time.Now()
	controller.Write(triggerData)

	fmt.Println("⏳ Waiting 4 seconds for delivery results...")
	time.Sleep(4 * time.Second)
	
	received := atomic.LoadInt64(&receivedCount)
	fmt.Printf("✅ Broadcast Stress Test finished.\n")
	fmt.Printf("----------------------------------\n")
	fmt.Printf("1. Registration Phase : %v\n", time.Since(startReg)-time.Since(startBroadcast))
	fmt.Printf("2. Delivery Report    :\n")
	fmt.Printf("   - Total Clients    : %d\n", *clientCount)
	fmt.Printf("   - Total Received   : %d\n", received)
	fmt.Printf("   - Delivery Rate    : %.2f%%\n", float64(received)/float64(*clientCount)*100)
	fmt.Printf("----------------------------------\n")
	fmt.Printf("💡 Note: Check server logs for internal execution time.\n")

	// Cleanup
	for _, c := range clients {
		if c != nil {
			c.Close()
		}
	}
}
