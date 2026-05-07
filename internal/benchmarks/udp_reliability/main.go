package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"mangahub/pkg/types"
	"net"
	"time"
)

func main() {
	target := flag.String("target", "localhost:8083", "Target UDP server address")
	messages := flag.Int("n", 2000, "Number of notifications to send")
	flag.Parse()

	addr, err := net.ResolveUDPAddr("udp", *target)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	fmt.Printf("📡 [UDP-BENCH] Starting Reliability Test\n")
	fmt.Printf("📡 Target: %s | Load: %d packets\n", *target, *messages)

	success := 0
	start := time.Now()

	for i := 1; i <= *messages; i++ {
		msg := types.UDPMessage{
			Action:  "benchmark:test_ping",
			Payload: json.RawMessage(fmt.Sprintf(`{"id": %d, "ts": %d}`, i, time.Now().UnixNano())),
		}
		data, _ := json.Marshal(msg)
		
		_, err := conn.Write(data)
		if err == nil {
			success++
		}
		
		if i%500 == 0 {
			fmt.Printf(">>> Progress: %d/%d packets dispatched...\n", i, *messages)
		}
		time.Sleep(1 * time.Millisecond)
	}

	duration := time.Since(start)
	fmt.Printf("\n--- 📊 UDP Reliability Results ---\n")
	fmt.Printf("✅ Total Sent     : %d\n", *messages)
	fmt.Printf("📈 Success Rate   : %.2f%%\n", float64(success)/float64(*messages)*100)
	fmt.Printf("⏱️  Execution Time : %v\n", duration)
	fmt.Printf("🚀 Throughput     : %.0f packets/sec\n", float64(*messages)/duration.Seconds())
	fmt.Printf("----------------------------------\n")
}
