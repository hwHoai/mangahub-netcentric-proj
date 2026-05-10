package udp_pools

import (
	"net"
)

type UDPPool interface {
	Register(userID string, addr *net.UDPAddr)
	Unregister(userID string, reason string)
	Broadcast(conn *net.UDPConn, id string, payload map[string]interface{})
	ProcessAck(notificationID string, userID string)
}
