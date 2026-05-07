package udp_services

// UDPKeySyncServices is the interface for syncing the public key to the UDP server.
type UDPKeySyncServices interface {
	SyncPublicKey(publicKeyPEM string) error
}
