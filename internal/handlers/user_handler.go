// internal/handlers/user_handler.go

package handlers

import (
	"go-rest-api/internal/middleware"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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

// GetUsers возвращает список всех пользователей (только для админов)
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
// (пользователь может получить только свои данные, админ - любые)
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// Get the authenticated user's ID from context
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	// Check if user is authorized to access this resource
	if userID.(uint) != uint(id) && role.(string) != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only access your own user data"})
		return
	}

	var user models.User
	if err := h.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user.ToResponse())
}

// CreateUser создает нового пользователя (только для админов)
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

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = string(hashedPassword)

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
		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		updateData["password"] = string(hashedPassword)
	}

	if input.Role != "" {
		// Only admins can change roles
		role, _ := c.Get("role")
		if role.(string) != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can change user roles"})
			return
		}
		updateData["role"] = input.Role
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

// internal/handlers/user_handler.go (дополнения)

// UpdateUserSelf позволяет пользователю обновить только свои данные
func (h *UserHandler) UpdateUserSelf(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID"})
		return
	}

	// Получаем ID аутентифицированного пользователя из контекста
	userID, _ := c.Get("user_id")

	// Проверяем, имеет ли пользователь право обновлять эти данные
	if userID.(uint) != uint(id) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Вы можете обновлять только свои данные"})
		return
	}

	var user models.User
	if err := h.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return
	}

	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Обновляем только разрешенные поля (имя, email, пароль)
	updateData := map[string]interface{}{}

	if input.Name != "" {
		updateData["name"] = input.Name
	}

	if input.Email != "" {
		// Проверка, не занят ли email другим пользователем
		var existingUser models.User
		if result := h.DB.Where("email = ? AND id != ?", input.Email, id).First(&existingUser); result.RowsAffected > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "Email уже используется"})
			return
		}
		updateData["email"] = input.Email
	}

	if input.Password != "" {
		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось хешировать пароль"})
			return
		}
		updateData["password"] = string(hashedPassword)
	}

	// Пользователь не может изменить свою роль
	if input.Role != "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Вы не можете изменить свою роль"})
		return
	}

	h.DB.Model(&user).Updates(updateData)

	c.JSON(http.StatusOK, user.ToResponse())
}

// UpdateUserAdmin обновляет данные пользователя с правами администратора
func (h *UserHandler) UpdateUserAdmin(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID"})
		return
	}

	var user models.User
	if err := h.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return
	}

	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Обновляем данные пользователя
	updateData := map[string]interface{}{}

	if input.Name != "" {
		updateData["name"] = input.Name
	}

	if input.Email != "" {
		// Проверка, не занят ли email другим пользователем
		var existingUser models.User
		if result := h.DB.Where("email = ? AND id != ?", input.Email, id).First(&existingUser); result.RowsAffected > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "Email уже используется"})
			return
		}
		updateData["email"] = input.Email
	}

	if input.Password != "" {
		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось хешировать пароль"})
			return
		}
		updateData["password"] = string(hashedPassword)
	}

	// Администратор может изменить роль пользователя
	if input.Role != "" {
		// Проверка на допустимые роли
		if input.Role != middleware.RoleUser &&
			input.Role != middleware.RoleAdmin &&
			input.Role != middleware.RoleModerator {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Недопустимая роль"})
			return
		}
		updateData["role"] = input.Role
	}

	h.DB.Model(&user).Updates(updateData)

	c.JSON(http.StatusOK, user.ToResponse())
}
