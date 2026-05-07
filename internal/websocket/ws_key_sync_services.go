package websocket_services

type WSKeySyncService interface {
	SyncPublicKey(publicKeyPEM string) error
}