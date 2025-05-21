package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"go-rest-api/internal/migrations"
)

func InitDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := migrations.RunMigrations(db); err != nil {
		return nil, err
	}

	log.Println("Database connected and migrated successfully")
	return db, nil
}
