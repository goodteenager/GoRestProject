package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"go-rest-api/internal/constants"

	"go-rest-api/internal/models"
)

type UserHandler struct {
	DB *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{DB: db}
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	var users []models.User
	if err := h.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}

	var response []models.UserResponse
	for _, user := range users {
		response = append(response, user.ToResponse())
	}

	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

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

func (h *UserHandler) CreateUser(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingUser models.User
	if err := h.DB.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User with that email already exists"})
		return
	}

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

	updateData := map[string]interface{}{}

	if input.Name != "" {
		updateData["name"] = input.Name
	}

	if input.Email != "" {
		var existingUser models.User
		if result := h.DB.Where("email = ? AND id != ?", input.Email, id).First(&existingUser); result.RowsAffected > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already in use"})
			return
		}
		updateData["email"] = input.Email
	}

	if input.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		updateData["password"] = string(hashedPassword)
	}

	if input.Role != "" {
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

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	if err := h.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	currentUserID, _ := c.Get("user_id")

	// Проверяем, пытается ли пользователь удалить самого себя
	if user.ID == currentUserID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You cannot delete yourself"})
		return
	}

	if err := h.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func (h *UserHandler) UpdateUserSelf(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID"})
		return
	}

	userID, _ := c.Get("user_id")

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

	updateData := map[string]interface{}{}

	if input.Name != "" {
		updateData["name"] = input.Name
	}

	if input.Email != "" {
		var existingUser models.User
		if result := h.DB.Where("email = ? AND id != ?", input.Email, id).First(&existingUser); result.RowsAffected > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "Email уже используется"})
			return
		}
		updateData["email"] = input.Email
	}

	if input.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось хешировать пароль"})
			return
		}
		updateData["password"] = string(hashedPassword)
	}

	if input.Role != "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Вы не можете изменить свою роль"})
		return
	}

	h.DB.Model(&user).Updates(updateData)

	c.JSON(http.StatusOK, user.ToResponse())
}

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

	updateData := map[string]interface{}{}

	if input.Name != "" {
		updateData["name"] = input.Name
	}

	if input.Email != "" {
		var existingUser models.User
		if result := h.DB.Where("email = ? AND id != ?", input.Email, id).First(&existingUser); result.RowsAffected > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "Email уже используется"})
			return
		}
		updateData["email"] = input.Email
	}

	if input.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось хешировать пароль"})
			return
		}
		updateData["password"] = string(hashedPassword)
	}

	if input.Role != "" {
		if input.Role != constants.RoleUser &&
			input.Role != constants.RoleAdmin &&
			input.Role != constants.RoleModerator {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Недопустимая роль"})
			return
		}
		updateData["role"] = input.Role
	}

	h.DB.Model(&user).Updates(updateData)

	c.JSON(http.StatusOK, user.ToResponse())
}

// Добавим новый метод для модератора
func (h *UserHandler) UpdateUserModerator(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	currentUserRole, _ := c.Get("role")

	var user models.User
	if err := h.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var input struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updateData := make(map[string]interface{})

	// Модератор может обновлять только имя и пароль
	if input.Name != "" {
		updateData["name"] = input.Name
	}

	if input.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		updateData["password"] = string(hashedPassword)
	}

	// Админ может обновлять любые поля кроме своей записи
	if currentUserRole == "admin" {
		var adminInput models.User
		if err := c.ShouldBindJSON(&adminInput); err == nil {
			if adminInput.Email != "" {
				updateData["email"] = adminInput.Email
			}
			if adminInput.Role != "" {
				updateData["role"] = adminInput.Role
			}
		}
	}

	if err := h.DB.Model(&user).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, user.ToResponse())
}
func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var user models.User
	if err := h.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user.ToResponse())
}
