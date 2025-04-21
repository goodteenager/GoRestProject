// internal/models/user.go

package models

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-" gorm:"index"`
	Name      string     `json:"name" binding:"required"`
	Email     string     `json:"email" binding:"required,email" gorm:"uniqueIndex"`
	Password  string     `json:"password" binding:"required"`
	Role      string     `json:"role" gorm:"default:user"`
}

// UserResponse is the safe representation of a User without sensitive data
type UserResponse struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
}

// ToResponse converts a User to a UserResponse (removing sensitive data)
func (u User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		Name:      u.Name,
		Email:     u.Email,
		Role:      u.Role,
	}
}

func (u *User) BeforeSave() error {
	// Установка роли по умолчанию, если не указана
	if u.Role == "" {
		u.Role = "user"
	}

	// Проверка, что роль допустима
	validRoles := map[string]bool{
		"user":      true,
		"admin":     true,
		"moderator": true,
	}

	if !validRoles[u.Role] {
		u.Role = "user" // Если недопустимая роль, устанавливаем "user"
	}

	return nil
}
