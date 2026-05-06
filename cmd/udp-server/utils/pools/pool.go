package udp_pools

import (
	"net"
)

type UDPPool interface {
	Register(userID string, addr *net.UDPAddr)
	Broadcast(conn *net.UDPConn, mangaID string, payload map[string]interface{})
}
