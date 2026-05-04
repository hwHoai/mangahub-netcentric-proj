package pools

import "net"

type ConnectionPool interface {
	Register(userID string, conn net.Conn)
	Unregister(conn net.Conn)
	Broadcast(userID string, message []byte)
	BroadcastOne(userID string, message []byte, conn net.Conn)
}
