package websocket_impl

import (
	"bytes"
	"encoding/json"
	"fmt"
	ws "mangahub/internal/websocket"
	"net/http"
	"os"
)

type WSKeySyncServiceImpl struct {
	addr string
}

func NewWSKeySyncService(addr string) ws.WSKeySyncService {
	return &WSKeySyncServiceImpl{addr: addr}
}

func (s *WSKeySyncServiceImpl) SyncPublicKey(publicKeyPEM string) error {
	url := fmt.Sprintf("http://%s/impl/sync-public-key", s.addr)

	payload, _ := json.Marshal(map[string]string{
		"public_key": publicKeyPEM,
	})

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Handshake-Key", os.Getenv("HANDSHAKE_KEY"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to sync public key, status: %d", resp.StatusCode)
	}

	return nil
}
