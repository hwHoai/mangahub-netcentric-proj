package handler
	
import (
	"fmt"
	"net"
)

func ExampleTCPHandler(conn net.Conn, payload any) {
	fmt.Printf("Request handled: %v\n", payload)
	conn.Write([]byte(`{"status": "Request handled successfully"}\n`))
}