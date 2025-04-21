// internal/database/db.go

package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"go-rest-api/internal/migrations"
)

// InitDB инициализирует подключение к базе данных и выполняет миграции
func InitDB(dsn string) (*gorm.DB, error) {
	// Подключение к базе данных
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Запуск миграций
	if err := migrations.RunMigrations(db); err != nil {
		return nil, err
	}

	log.Println("Database connected and migrated successfully")
	return db, nil
}
