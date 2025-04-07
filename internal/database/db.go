package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"go-rest-api/internal/models"
)

// InitDB инициализирует подключение к базе данных и выполняет миграции
func InitDB(dsn string) (*gorm.DB, error) {
	// Подключение к базе данных
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Автоматическая миграция моделей
	if err := db.AutoMigrate(&models.User{}); err != nil {
		return nil, err
	}

	log.Println("Database connected and migrated successfully")
	return db, nil
}

// GetDB возвращает экземпляр базы данных
func GetDB(db *gorm.DB) *gorm.DB {
	return db
}
