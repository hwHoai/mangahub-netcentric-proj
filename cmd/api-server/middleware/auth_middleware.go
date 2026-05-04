package middleware

import (
	"crypto/rsa"
	"fmt"
	"mangahub/pkg/utils/jwt"
	jwt_impl "mangahub/pkg/utils/jwt/impl"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	publicKey  *rsa.PublicKey
	jwtUtil    jwt.JWTUtil
}

func NewAuthMiddleware(publicKey *rsa.PublicKey) *AuthMiddleware {
	return &AuthMiddleware{
		publicKey: publicKey,
		jwtUtil:   jwt_impl.NewJWTUtil(nil),
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
		jwtClaims, err := am.jwtUtil.VerifyJWTToken(accessToken, am.publicKey)
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
