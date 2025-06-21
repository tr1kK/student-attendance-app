package auth

import (
	"fmt"
	"student-attendance-app/pkg/config"
	"student-attendance-app/pkg/models"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenerateJWT(user models.User, cfg *config.Config) (string, error) {
	claims := jwt.MapClaims{
		"id":   user.ID,
		"role": user.Role,
		"exp":  time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWTSecret))
}

func ValidateJWT(tokenString string, cfg *config.Config) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cfg.JWTSecret), nil
	})
} 