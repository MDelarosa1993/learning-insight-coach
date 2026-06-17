package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"learning-insight-coach/config"
)

func APIKeyAuth(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		expected := cfg.APIKey

		if expected == "" {
			c.Next()
			return
		}

		provided := c.GetHeader("X-API-Key")

		if provided != expected {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "missing or invalid API key",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}