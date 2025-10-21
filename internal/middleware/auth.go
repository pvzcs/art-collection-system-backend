package middleware

import (
	"art-collection-system/internal/service"
	"art-collection-system/internal/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT token from Authorization header
// and stores user_id and user_role in Gin Context
func AuthMiddleware(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.Error(c, 401, "未授权")
			c.Abort()
			return
		}

		// Remove "Bearer " prefix
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			// No "Bearer " prefix found
			utils.Error(c, 401, "令牌格式无效")
			c.Abort()
			return
		}

		// Validate token
		user, err := authService.ValidateToken(token)
		if err != nil {
			utils.Error(c, 401, "令牌无效")
			c.Abort()
			return
		}

		// Store user_id and user_role in context
		c.Set("user_id", user.ID)
		c.Set("user_role", user.Role)

		c.Next()
	}
}
