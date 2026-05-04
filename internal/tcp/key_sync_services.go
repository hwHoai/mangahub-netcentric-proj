package tcp_services

type TCPKeySyncServices interface {
	SyncPublicKey(publicKeyPEM string) error
}
