package migrations

import (
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

// Migration represents a single database migration
type Migration struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"uniqueIndex"`
	CreatedAt time.Time
}

// MigrationFunc defines the signature of a migration function
type MigrationFunc func(*gorm.DB) error

// migrationRegistry stores all registered migrations
var migrationRegistry = make(map[string]MigrationFunc)

// Register adds a new migration to the registry
func Register(name string, fn MigrationFunc) {
	migrationRegistry[name] = fn
}

// RunMigrations runs all pending migrations
func RunMigrations(db *gorm.DB) error {
	// Create migrations table if it doesn't exist
	if err := db.AutoMigrate(&Migration{}); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Run each migration in the registry
	for name, fn := range migrationRegistry {
		// Check if migration has already been applied
		var count int64
		db.Model(&Migration{}).Where("name = ?", name).Count(&count)
		if count > 0 {
			// Migration already applied, skip
			continue
		}

		log.Printf("Running migration: %s", name)

		// Run the migration inside a transaction
		if err := db.Transaction(func(tx *gorm.DB) error {
			// Run the migration function
			if err := fn(tx); err != nil {
				return err
			}

			// Record the migration
			return tx.Create(&Migration{Name: name}).Error
		}); err != nil {
			return fmt.Errorf("migration '%s' failed: %w", name, err)
		}

		log.Printf("Migration completed: %s", name)
	}

	return nil
}
