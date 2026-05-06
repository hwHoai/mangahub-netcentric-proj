package handler

import (
	"encoding/json"
	"fmt"
	"mangahub/cmd/udp-server/dispatch"
	"mangahub/cmd/udp-server/utils"
	"mangahub/pkg/types"
	jwt_impl "mangahub/pkg/utils/jwt/impl"
	"net"
)

type KeySyncHandler struct{}

func NewKeySyncHandler() *KeySyncHandler {
	return &KeySyncHandler{}
}

func (h *KeySyncHandler) SyncPublicKeyHandler(s *dispatch.UDPServer, clientAddr *net.UDPAddr, payload types.UDPMessage) {
	var data struct {
		PublicKey string `json:"public_key"`
	}

	if err := json.Unmarshal(payload.Payload, &data); err != nil {
		fmt.Printf("Error unmarshaling sync_public_key payload: %v\n", err)
		return
	}

	fmt.Println("Received Public Key from API Server. Ready to verify tokens.")
	
	// Parse and store in global variable
	jwtUtil := jwt_impl.NewJWTUtil(nil)
	pubKey, err := jwtUtil.ParsePublicKeyPEM(data.PublicKey)
	if err != nil {
		fmt.Printf("Error parsing public key: %v\n", err)
		return
	}

	utils.SetPublicKey(pubKey)
	fmt.Println("Public Key successfully stored and updated for UDP Server.")
}
