package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggerMiddleware records all API requests (method, path, status code, duration)
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		startTime := time.Now()

		// Get request information
		method := c.Request.Method
		path := c.Request.URL.Path
		clientIP := c.ClientIP()

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(startTime)

		// Get response status code
		statusCode := c.Writer.Status()

		// Log the request
		// TODO: In production, use a proper logging library like Zap
		fmt.Printf("[%s] %s %s %d %v %s\n",
			time.Now().Format("2006-01-02 15:04:05"),
			method,
			path,
			statusCode,
			duration,
			clientIP,
		)
	}
}
