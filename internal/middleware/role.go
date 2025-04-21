// internal/middleware/role.go

package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Константы для ролей пользователей
const (
	RoleUser      = "user"
	RoleAdmin     = "admin"
	RoleModerator = "moderator" // Добавим еще одну роль для примера
)

// RoleAuthMiddleware проверяет, имеет ли пользователь необходимую роль
func RoleAuthMiddleware(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем роль из контекста (устанавливается в JWTAuthMiddleware)
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Роль не найдена в токене"})
			c.Abort()
			return
		}

		// Проверяем, есть ли у пользователя требуемая роль
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

//// CheckPermission проверяет права доступа для определенного ресурса
//func CheckPermission(c *gin.Context, resourceOwnerID uint) bool {
//	// Получаем ID пользователя и роль из контекста
//	userID, _ := c.Get("user_id")
//	role, _ := c.Get("role")
//
//	// Администраторы могут получить доступ к любому ресурсу
//	if role.(string) == RoleAdmin {
//		return true
//	}
//
//	// Обычные пользователи могут получить доступ только к своим ресурсам
//	return userID.(uint) == resourceOwnerID
//}

//// OwnerOrAdminMiddleware проверяет, является ли пользователь владельцем ресурса или админом
//func OwnerOrAdminMiddleware() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		// Код будет добавлен в обработчики для проверки
//		c.Next()
//	}
//}
