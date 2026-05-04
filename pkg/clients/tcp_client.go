package clients

import (
	"mangahub/internal/tcp"
	"mangahub/internal/tcp/impl"
	"os"
)

func NewTCPChapterSyncClient() tcp_services.TCPChapterSyncServices {
	serverHost := os.Getenv("SERVER_HOST")
	if serverHost == "" {
		serverHost = "localhost"
	}
	port := os.Getenv("TCP_SERVER_PORT")
	if port == "" {
		port = "8082"
	}
	addr := serverHost + ":" + port

	return tcp_services_impl.NewTCPChapterSyncService(addr)
}
