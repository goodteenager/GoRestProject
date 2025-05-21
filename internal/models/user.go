package models

import (
	"errors"
	"go-rest-api/internal/constants"
	"time"
)

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

type UserResponse struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
}

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
	validRoles := map[string]bool{
		constants.RoleUser:      true,
		constants.RoleAdmin:     true,
		constants.RoleModerator: true,
	}

	if !validRoles[u.Role] {
		return errors.New("invalid user role")
	}

	return nil
}
