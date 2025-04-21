package middleware

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go-rest-api/internal/models"
)

// JWTAuthMiddleware handles JWT authentication
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Check if the header has the Bearer prefix
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header must be in the format: Bearer {token}"})
			c.Abort()
			return
		}

		// Parse the JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: " + err.Error()})
			c.Abort()
			return
		}

		// Validate the token
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Check if the token is expired
			if exp, ok := claims["exp"].(float64); ok {
				if time.Now().Unix() > int64(exp) {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
					c.Abort()
					return
				}
			}

			// Set claims to context for later use in handlers
			userID, ok := claims["user_id"].(float64)
			if !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
				c.Abort()
				return
			}

			// Set user ID and role in context
			c.Set("user_id", uint(userID))
			if role, ok := claims["role"].(string); ok {
				c.Set("role", role)
			}

			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
	}
}

// RoleAuthMiddleware checks if user has the required role
func RoleAuthMiddleware(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get role from context (set by JWTAuthMiddleware)
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Role not found in token"})
			c.Abort()
			return
		}

		// Check if user has required role
		userRole := role.(string)
		for _, r := range requiredRoles {
			if r == userRole {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		c.Abort()
	}
}

// GenerateToken creates a new JWT token for a user
func GenerateToken(user models.User) (string, error) {
	// Get JWT secret from environment variable
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default_jwt_secret_key" // Default secret (should be overridden in production)
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	})

	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
