package main

import (
	"log"
	"student-attendance-app/pkg/config"
	"student-attendance-app/pkg/database"
	"student-attendance-app/pkg/router"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load environment variables
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	// Initialize database
	db, err := database.InitDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("could not initialize database: %v", err)
	}

	// Set up router
	r := gin.Default()
	router.SetupRouter(r, db)

	// Start server
	if err := r.Run(cfg.ServerAddress); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
} 