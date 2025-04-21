// internal/router/router.go

package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"go-rest-api/internal/handlers"
	"go-rest-api/internal/middleware"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	// Создаем обработчики
	userHandler := handlers.NewUserHandler(db)
	authHandler := handlers.NewAuthHandler(db)

	// Группа API
	api := router.Group("/api")
	{
		// Публичные маршруты аутентификации
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/register", authHandler.Register)
		}

		// Защищенные маршруты пользователей (с JWT аутентификацией)
		users := api.Group("/users")
		users.Use(middleware.JWTAuthMiddleware()) // Добавляем JWT аутентификацию
		{
			// Маршруты доступные всем аутентифицированным пользователям
			users.GET("/:id", userHandler.GetUser)
			users.PUT("/:id", userHandler.UpdateUserSelf) // Новый метод для обновления только своих данных

			// Защищенные маршруты для модераторов и администраторов
			modRoutes := users.Group("/")
			modRoutes.Use(middleware.RoleAuthMiddleware(middleware.RoleAdmin, middleware.RoleModerator))
			{
				modRoutes.GET("", userHandler.GetUsers)
			}

			// Защищенные маршруты для только администраторов
			adminRoutes := users.Group("/")
			adminRoutes.Use(middleware.RoleAuthMiddleware(middleware.RoleAdmin))
			{
				adminRoutes.POST("", userHandler.CreateUser)
				adminRoutes.PUT("/:id/admin", userHandler.UpdateUserAdmin) // Обновление с правами админа
				adminRoutes.DELETE("/:id", userHandler.DeleteUser)
			}
		}
	}

	return router
}
