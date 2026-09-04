package middleware

import (
	"crypto/subtle"
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIKeyAuth rejects any request without a valid X-API-Key header.
// validKey is the configured secret; loaded from env, never hardcoded.
func APIKeyAuth(validKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		provided := c.GetHeader("X-API-Key")

		// constant-time compare — avoids timing attacks that can leak the key
		if provided == "" || subtle.ConstantTimeCompare([]byte(provided), []byte(validKey)) != 1 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or missing API key",
				"code":  "UNAUTHORIZED",
			})
			return
		}
		c.Next() // key valid — proceed to the handler
	}
}
