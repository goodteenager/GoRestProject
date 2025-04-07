package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"go-rest-api/internal/config"
	"go-rest-api/internal/database"
	"go-rest-api/internal/router"
)

func main() {
	// Загрузка переменных окружения
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Инициализация конфигурации
	cfg := config.NewConfig()

	// Инициализация базы данных
	db, err := database.InitDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Настройка и запуск маршрутизатора
	r := router.SetupRouter(db)

	// Запуск сервера
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
