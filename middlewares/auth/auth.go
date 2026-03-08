package auth

import (
	"os"

	"github.com/gin-gonic/gin"
)

func CheckSecretKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKeyHeader := c.GetHeader("x-api-key")
		secretKey := os.Getenv("SUPER_SECRET_KEY")

		if len(apiKeyHeader) > 0 && len(secretKey) > 0 && apiKeyHeader == secretKey {
			c.Next()
		} else {
			c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
		}
	}
}
