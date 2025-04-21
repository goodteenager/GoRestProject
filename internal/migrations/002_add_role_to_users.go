// internal/migrations/002_add_role_to_users.go

package migrations

import (
	"gorm.io/gorm"
)

func init() {
	Register("002_add_role_to_users", addRoleToUsers)
}

func addRoleToUsers(db *gorm.DB) error {
	return db.Exec(`
		ALTER TABLE users ADD COLUMN IF NOT EXISTS role VARCHAR(50) DEFAULT 'user' NOT NULL;
	`).Error
}
