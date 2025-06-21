package handlers

import (
	"math/rand"
	"net/http"
	"strconv"
	"student-attendance-app/pkg/auth"
	"student-attendance-app/pkg/config"
	"student-attendance-app/pkg/models"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Auth Handlers
func Login(c *gin.Context, db *gorm.DB, cfg *config.Config) {
	var req struct {
		StudentID string `json:"student_id" binding:"required"`
		Password  string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := db.First(&user, "student_id = ?", req.StudentID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := auth.GenerateJWT(user, cfg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// Lesson Handlers
func GetLessons(c *gin.Context, db *gorm.DB) {
	var lessons []models.Lesson
	if err := db.Find(&lessons).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve lessons"})
		return
	}
	c.JSON(http.StatusOK, lessons)
}

// Attendance Handlers
func SubmitAttendance(c *gin.Context, db *gorm.DB) {
	var req struct {
		LessonID uint   `json:"lesson_id" binding:"required"`
		Code     string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userID")

	// Validate code
	var generatedCode models.GeneratedCode
	err := db.Where("lesson_id = ? AND code = ? AND is_active = ? AND expires_at > ?", req.LessonID, req.Code, true, time.Now()).
		Order("created_at desc").
		First(&generatedCode).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired code"})
		return
	}

	// Save attendance
	attendance := models.Attendance{
		LessonID:    req.LessonID,
		StudentID:   uint(userID.(float64)),
		SubmittedAt: time.Now(),
	}
	if err := db.Create(&attendance).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save attendance"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Attendance marked successfully"})
}

// Code Handlers
func GenerateCode(c *gin.Context, db *gorm.DB) {
	var req struct {
		LessonID uint `json:"lesson_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Deactivate previous codes for this lesson
	db.Model(&models.GeneratedCode{}).Where("lesson_id = ?", req.LessonID).Update("is_active", false)

	// Generate new code
	code := strconv.Itoa(10000 + rand.Intn(90000))
	newCode := models.GeneratedCode{
		LessonID:  req.LessonID,
		Code:      code,
		ExpiresAt: time.Now().Add(15 * time.Minute),
		IsActive:  true,
	}

	if err := db.Create(&newCode).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate code"})
		return
	}

	c.JSON(http.StatusOK, newCode)
} 