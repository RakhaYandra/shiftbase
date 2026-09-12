package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORS minimal untuk shiftbase-web (dev :5173). Tanpa lib tambahan.
func CORS(frontendURL string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if origin := c.GetHeader("Origin"); origin == frontendURL {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Headers", "Authorization,Content-Type")
			c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		}
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
