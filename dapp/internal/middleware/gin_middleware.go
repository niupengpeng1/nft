package middleware

import (
	"net/http"
	"time"

	"dapp/pkg/logger"
	"github.com/gin-gonic/gin"
)

// CORSMiddleware returns a Gin middleware function for CORS support
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")
		
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	}
}

// LoggingMiddleware returns a Gin middleware function for request logging
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		
		c.Next()
		
		latency := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		
		logger.Info(
			"[GIN] %s %s %d %v %s %s",
			method,
			path,
			statusCode,
			latency,
			clientIP,
			query,
		)
	}
}

// RecoveryMiddleware returns a Gin middleware function for panic recovery
func RecoveryMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		logger.Error("Panic recovered: %v", recovered)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Internal server error",
		})
	})
}

// RateLimitMiddleware is a placeholder for rate limiting
// TODO: Implement actual rate limiting logic
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implement rate limiting here
		// Example: using github.com/juju/ratelimit or redis
		c.Next()
	}
}

// AuthMiddleware is a placeholder for authentication
// TODO: Implement JWT or API key authentication
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from header
		// token := c.GetHeader("Authorization")
		
		// Validate token
		// if invalid {
		//     c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		//     c.Abort()
		//     return
		// }
		
		c.Next()
	}
}
