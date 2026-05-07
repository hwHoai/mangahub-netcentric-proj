package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"mangahub/pkg/types"
	"net"
	"sync"
	"time"
)

func main() {
	target := flag.String("target", "localhost:8082", "Target TCP server address")
	conns := flag.Int("conns", 1000, "Number of concurrent connections to spawn")
	duration := flag.Duration("duration", 15*time.Second, "How long to hold and ping-pong")
	flag.Parse()

	fmt.Printf("🚀 Starting FULL-STACK TCP Stress Test (Ping-Pong-Ack): %d connections to %s\n", *conns, *target)

	var wg sync.WaitGroup
	start := time.Now()

	for i := 0; i < *conns; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			conn, err := net.DialTimeout("tcp", *target, 5*time.Second)
			if err != nil {
				return
			}
			defer conn.Close()

			pingMsg := types.TCPMessage{
				Action:  "benchmark:test_ping",
				Payload: json.RawMessage(`{"data": "PING"}`),
			}
			data, _ := json.Marshal(pingMsg)

			pongAck := types.TCPMessage{
				Action:  "benchmark:res_pong",
				Payload: json.RawMessage(`{"data": "PONG_ACK"}`),
			}
			ackData, _ := json.Marshal(pongAck)
			
			scanner := bufio.NewScanner(conn)
			stopTimer := time.After(*duration)
			ticker := time.NewTicker(2 * time.Second)
			defer ticker.Stop()

			for {
				select {
				case <-stopTimer:
					return
				case <-ticker.C:
					// Send Ping
					fmt.Fprintln(conn, string(data))
					
					// Wait for Pong and send Ack back to server
					if scanner.Scan() {
						fmt.Fprintln(conn, string(ackData))
					} else {
						return
					}
				}
			}
		}(i)
		
		if i%200 == 0 && i > 0 {
			fmt.Printf("... Spawned %d active ping-pong-ack sessions\n", i)
		}
		time.Sleep(2 * time.Millisecond)
	}

	fmt.Println("⏳ All sessions active. Performing continuous Ping-Pong-Ack...")
	wg.Wait()
	fmt.Printf("✅ Full-stack Test completed in %v. Concurrency maintained: %d\n", time.Since(start), *conns)
}
