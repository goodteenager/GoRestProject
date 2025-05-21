package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RoleAuthMiddleware(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Роль не найдена в токене"})
			c.Abort()
			return
		}

		userRole := role.(string)
		for _, r := range requiredRoles {
			if r == userRole {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "Недостаточно прав"})
		c.Abort()
	}
}
