package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"go-rest-api/internal/models"
)

// UserHandler содержит обработчики для работы с пользователями
type UserHandler struct {
	DB *gorm.DB
}

// NewUserHandler создает новый экземпляр UserHandler
func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{DB: db}
}

// GetUsers возвращает список всех пользователей
func (h *UserHandler) GetUsers(c *gin.Context) {
	var users []models.User
	if err := h.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}

	// Преобразование списка пользователей в ответ без паролей
	var response []models.UserResponse
	for _, user := range users {
		response = append(response, user.ToResponse())
	}

	c.JSON(http.StatusOK, response)
}

// GetUser возвращает пользователя по ID
func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if err := h.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user.ToResponse())
}

// CreateUser создает нового пользователя
func (h *UserHandler) CreateUser(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверяем, существует ли пользователь с таким email
	var existingUser models.User
	if err := h.DB.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User with that email already exists"})
		return
	}

	if err := h.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, user.ToResponse())
}

// UpdateUser обновляет данные пользователя
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var user models.User
	if err := h.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Обновляем только те поля, которые предоставлены
	updateData := map[string]interface{}{}

	if input.Name != "" {
		updateData["name"] = input.Name
	}

	if input.Email != "" {
		// Проверка, не занят ли email другим пользователем
		var existingUser models.User
		if result := h.DB.Where("email = ? AND id != ?", input.Email, id).First(&existingUser); result.RowsAffected > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already in use"})
			return
		}
		updateData["email"] = input.Email
	}

	if input.Password != "" {
		updateData["password"] = input.Password
	}

	h.DB.Model(&user).Updates(updateData)

	c.JSON(http.StatusOK, user.ToResponse())
}

// DeleteUser удаляет пользователя
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if err := h.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	h.DB.Delete(&user)
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
