package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"go-rest-api/internal/handlers"
)

// SetupRouter настраивает маршруты API
func SetupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	// Создаем обработчик пользователей
	userHandler := handlers.NewUserHandler(db)

	// Группа API
	api := router.Group("/api")
	{
		// Маршруты пользователей
		users := api.Group("/users")
		{
			users.GET("", userHandler.GetUsers)
			users.GET("/:id", userHandler.GetUser)
			users.POST("", userHandler.CreateUser)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
		}
	}

	return router
}
