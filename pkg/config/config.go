package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL   string
	JWTSecret     string
	ServerAddress string
}

func LoadConfig() (*Config, error) {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		// Ignore if .env file is not found
		if !os.IsNotExist(err) {
			return nil, err
		}
	}

	// Read environment variables
	dbURL := os.Getenv("DATABASE_URL")
	jwtSecret := os.Getenv("JWT_SECRET")
	serverAddr := os.Getenv("SERVER_ADDRESS")

	if dbURL == "" {
		dbURL = "postgres://postgres:1234@localhost:5432/attendance_db?sslmode=disable"
	}
	if jwtSecret == "" {
		jwtSecret = "your-very-secret-key"
	}
	if serverAddr == "" {
		serverAddr = ":8080"
	}

	return &Config{
		DatabaseURL:   dbURL,
		JWTSecret:     jwtSecret,
		ServerAddress: serverAddr,
	}, nil
}
