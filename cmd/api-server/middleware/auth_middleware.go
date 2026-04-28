package middleware

import (
	"crypto/rsa"
	"fmt"
	"mangahub/internal/auth"
	auth_service_impl "mangahub/internal/auth/impl"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	publicKey  *rsa.PublicKey
	jwtService auth.JWTService
}

func NewAuthMiddleware(publicKey *rsa.PublicKey) *AuthMiddleware {
	return &AuthMiddleware{
		publicKey:  publicKey,
		jwtService: auth_service_impl.NewJWTService(),
	}
}

func (am *AuthMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{
				"error": "missing authorization header",
			})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>" format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(401, gin.H{
				"error": "invalid authorization header format",
			})
			c.Abort()
			return
		}

		accessToken := parts[1]

		// Verify JWT token using the in-memory public key.
		jwtClaims, err := am.jwtService.VerifyJWTToken(accessToken, am.publicKey)
		if err != nil {
			c.JSON(401, gin.H{
				"error": fmt.Sprintf("invalid or expired token: %v", err),
			})
			c.Abort()
			return
		}

		// Token is valid, set user info in context and call next handler
		c.Set("user_id", jwtClaims.Subject)
		c.Set("claims", jwtClaims)
		c.Next()
	}
}
