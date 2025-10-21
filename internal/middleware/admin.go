package middleware

import (
	"art-collection-system/internal/utils"

	"github.com/gin-gonic/gin"
)

// AdminMiddleware validates that the user has admin role
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user_role from context (set by AuthMiddleware)
		role, exists := c.Get("user_role")
		if !exists {
			utils.Error(c, 403, "需要管理员权限")
			c.Abort()
			return
		}

		// Verify role is "admin"
		if role != "admin" {
			utils.Error(c, 403, "需要管理员权限")
			c.Abort()
			return
		}

		c.Next()
	}
}
