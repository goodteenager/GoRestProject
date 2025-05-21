package router

import (
	"github.com/gin-gonic/gin"
	"go-rest-api/internal/constants"
	"go-rest-api/internal/handlers"
	"go-rest-api/internal/middleware"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	router.Use(middleware.CORSMiddleware())

	router.Static("/static", "./static")
	router.LoadHTMLGlob("static/*.html")

	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	userHandler := handlers.NewUserHandler(db)
	authHandler := handlers.NewAuthHandler(db)

	api := router.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/register", authHandler.Register)
		}

		users := api.Group("/users")
		users.Use(middleware.JWTAuthMiddleware())
		{
			users.GET("/:id", userHandler.GetUser)
			users.PUT("/:id", userHandler.UpdateUserSelf)

			modAdminRoutes := users.Group("")
			modAdminRoutes.Use(middleware.RoleAuthMiddleware(constants.RoleAdmin, constants.RoleModerator))
			{
				modAdminRoutes.GET("", userHandler.GetUsers)
				modAdminRoutes.PUT("/:id/mod", userHandler.UpdateUserModerator)
			}

			adminRoutes := users.Group("")
			adminRoutes.Use(middleware.RoleAuthMiddleware(constants.RoleAdmin))
			{
				adminRoutes.POST("", userHandler.CreateUser)
				adminRoutes.DELETE("/:id", userHandler.DeleteUser)
				adminRoutes.PUT("/:id/admin", userHandler.UpdateUserAdmin)
			}
		}
	}

	return router
}
