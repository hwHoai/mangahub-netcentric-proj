package controllers

import (
	"fmt"
	"mangahub/cmd/websocket-server/utils"
	"mangahub/pkg/logger"
	jwt_impl "mangahub/pkg/utils/jwt/impl"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type KeySyncController struct{}

func NewKeySyncController() *KeySyncController {
	return &KeySyncController{}
}

func (h *KeySyncController) SyncPublicKeyHandler(c *gin.Context) {
	// Simple handshake check
	handshakeKey := c.GetHeader("X-Handshake-Key")
	if handshakeKey != os.Getenv("HANDSHAKE_KEY") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized handshake"})
		return
	}

	var data struct {
		PublicKey string `json:"public_key"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	jwtUtil := jwt_impl.NewJWTUtil(nil)
	pubKey, err := jwtUtil.ParsePublicKeyPEM(data.PublicKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error parsing public key: %v", err)})
		return
	}

	utils.SetPublicKey(pubKey)
	logger.Info("Public key synced", "public_key", pubKey)
	c.JSON(http.StatusOK, gin.H{"message": "Public key synced"})
}
