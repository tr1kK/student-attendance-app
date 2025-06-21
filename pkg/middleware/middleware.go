package middleware

import (
	"net/http"
	"student-attendance-app/pkg/auth"
	"student-attendance-app/pkg/config"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := auth.ValidateJWT(tokenString, cfg)

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		c.Set("userID", claims["id"])
		c.Set("userRole", claims["role"])
		c.Next()
	}
}

func RoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("userRole")
		if !exists || userRole.(string) != requiredRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "You do not have permission to access this resource"})
			return
		}
		c.Next()
	}
} 