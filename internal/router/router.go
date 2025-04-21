// internal/router/router.go

package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"go-rest-api/internal/handlers"
	"go-rest-api/internal/middleware"
)

// SetupRouter настраивает маршруты API
func SetupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	// Создаем обработчики
	userHandler := handlers.NewUserHandler(db)
	authHandler := handlers.NewAuthHandler(db)

	// Группа API
	api := router.Group("/api")
	{
		// Публичные маршруты аутентификации (без middleware)
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/register", authHandler.Register)
		}

		// Защищенные маршруты пользователей (с middleware)
		users := api.Group("/users")
		users.Use(middleware.JWTAuthMiddleware()) // Добавляем JWT аутентификацию
		{
			// Маршруты доступные всем аутентифицированным пользователям
			users.GET("/:id", userHandler.GetUser)

			// Маршруты доступные только администраторам
			adminRoutes := users.Group("/")
			adminRoutes.Use(middleware.RoleAuthMiddleware("admin"))
			{
				adminRoutes.GET("", userHandler.GetUsers)
				adminRoutes.POST("", userHandler.CreateUser)
				adminRoutes.PUT("/:id", userHandler.UpdateUser)
				adminRoutes.DELETE("/:id", userHandler.DeleteUser)
			}
		}
	}

	return router
}
