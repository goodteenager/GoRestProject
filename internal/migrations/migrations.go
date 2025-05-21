package migrations

import (
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

type Migration struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"uniqueIndex"`
	CreatedAt time.Time
}

type MigrationFunc func(*gorm.DB) error

var migrationRegistry = make(map[string]MigrationFunc)

func Register(name string, fn MigrationFunc) {
	migrationRegistry[name] = fn
}

func RunMigrations(db *gorm.DB) error {
	if err := db.AutoMigrate(&Migration{}); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	for name, fn := range migrationRegistry {
		var count int64
		db.Model(&Migration{}).Where("name = ?", name).Count(&count)
		if count > 0 {
			continue
		}

		log.Printf("Running migration: %s", name)

		if err := db.Transaction(func(tx *gorm.DB) error {
			if err := fn(tx); err != nil {
				return err
			}

			return tx.Create(&Migration{Name: name}).Error
		}); err != nil {
			return fmt.Errorf("migration '%s' failed: %w", name, err)
		}

		log.Printf("Migration completed: %s", name)
	}
	return nil
}
