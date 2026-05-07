package middleware

import (
	"fmt"
	"mangahub/cmd/websocket-server/utils"
	jwt_impl "mangahub/pkg/utils/jwt/impl"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Query("token")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}

		pubKey := utils.GetPublicKey()
		if pubKey == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Public key not synced yet"})
			c.Abort()
			return
		}

		jwtUtil := jwt_impl.NewJWTUtil(nil)
		claims, err := jwtUtil.VerifyJWTToken(tokenString, pubKey)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Invalid token: %v", err)})
			c.Abort()
			return
		}

		// Set user info in context
		c.Set("user_id", claims.Subject)
		c.Next()
	}
}
