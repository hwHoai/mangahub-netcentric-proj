package handler

import (
	"encoding/json"
	"mangahub/cmd/tcp-server/utils"
	"mangahub/pkg/logger"
	jwt_impl "mangahub/pkg/utils/jwt/impl"
	"net"
)

type KeySyncHandler struct{}

func NewKeySyncHandler() *KeySyncHandler {
	return &KeySyncHandler{}
}

func (h *KeySyncHandler) SyncPublicKeyHandler(conn net.Conn, payload any) {
	var data struct {
		PublicKey string `json:"public_key"`
	}

	raw, ok := payload.(json.RawMessage)
	if !ok {
		logger.Error("Invalid payload type for sync_public_key")
		return
	}

	if err := json.Unmarshal(raw, &data); err != nil {
		logger.Error("Error unmarshaling sync_public_key payload", "error", err)
		return
	}

	logger.Info("Received Public Key from API Server. Ready to verify tokens.")
	
	// Parse and store in global variable
	jwtUtil := jwt_impl.NewJWTUtil(nil)
	pubKey, err := jwtUtil.ParsePublicKeyPEM(data.PublicKey)
	if err != nil {
		logger.Error("Error parsing public key", "error", err)
		return
	}

	utils.SetPublicKey(pubKey)
	logger.Info("Public Key successfully stored and updated.")
}
