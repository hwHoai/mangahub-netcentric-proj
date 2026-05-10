package middleware

import (
	"errors"
	"mangahub/cmd/tcp-server/utils"
	"mangahub/pkg/logger"
	jwt_impl "mangahub/pkg/utils/jwt/impl"
	"os"
	"strings"
)

func AuthMiddleware(action, token string) error {
	// 1. Determine category based on prefix or content
	category := "unknown"
	if strings.HasPrefix(action, "req_") || strings.Contains(action, ":req_") {
		category = "req"
	} else if strings.HasPrefix(action, "res_") || strings.Contains(action, ":res_") {
		category = "res"
	} else if strings.HasPrefix(action, "impl") || strings.Contains(action, ":impl") {
		category = "impl"
	} else if strings.HasPrefix(action, "test_") || strings.Contains(action, ":test_") {
		category = "test"
	}

	// 2. Switch case for each category
	switch category {
	case "req":
		// JWT verify
		publicKey := utils.GetPublicKey()
		if publicKey == nil {
			return errors.New("unauthorized: public key not synced yet")
		}

		jwtUtil := jwt_impl.NewJWTUtil(nil)
		_, err := jwtUtil.VerifyJWTToken(token, publicKey)
		if err != nil {
			return errors.New("unauthorized: invalid or expired token")
		}
		return nil

	case "res":
		// Allowed directly
		return nil

	case "impl":
		// Handshake key check
		handshakeKey := os.Getenv("HANDSHAKE_KEY")
		if handshakeKey == "" {
			logger.Warn("HANDSHAKE_KEY environment variable is not set")
		}
		if token != handshakeKey {
			return errors.New("forbidden: invalid handshake key")
		}
		return nil
	
	case "test":
		// Allowed for benchmarking
		return nil

	default:
		logger.Warn("Action has no recognized security category", "action", action)
		return nil
	}
}
