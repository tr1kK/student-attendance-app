package router

import (
	"student-attendance-app/pkg/config"
	"student-attendance-app/pkg/handlers"
	"student-attendance-app/pkg/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(r *gin.Engine, db *gorm.DB) {
	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	// Public routes
	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/login", func(c *gin.Context) {
			handlers.Login(c, db, cfg)
		})
	}

	// Authenticated routes
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware(cfg))
	{
		// Lesson routes (accessible to all authenticated users)
		api.GET("/lessons", func(c *gin.Context) {
			handlers.GetLessons(c, db)
		})

		// Student routes
		studentRoutes := api.Group("/student")
		studentRoutes.Use(middleware.RoleMiddleware("student"))
		{
			studentRoutes.POST("/attendance", func(c *gin.Context) {
				handlers.SubmitAttendance(c, db)
			})
		}

		// Teacher routes
		teacherRoutes := api.Group("/teacher")
		teacherRoutes.Use(middleware.RoleMiddleware("teacher"))
		{
			teacherRoutes.POST("/codes", func(c *gin.Context) {
				handlers.GenerateCode(c, db)
			})
		}
	}
} 