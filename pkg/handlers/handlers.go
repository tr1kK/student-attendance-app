package handlers

import (
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"student-attendance-app/pkg/auth"
	"student-attendance-app/pkg/config"
	"student-attendance-app/pkg/models"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type LoginRequest struct {
	Identifier string `json:"identifier" binding:"required" example:"student1"`
	Password   string `json:"password" binding:"required" example:"password123"`
}

type RegisterRequest struct {
	Identifier string `json:"identifier" binding:"required" example:"newstudent"`
	Password   string `json:"password" binding:"required" example:"securepassword"`
	Name       string `json:"name" binding:"required" example:"New Student"`
	Email      string `json:"email" binding:"required" example:"new@example.com"`
	GroupID    uint   `json:"group_id" binding:"required" example:"1"`
}

type SubmitAttendanceRequest struct {
	LessonID uint   `json:"lesson_id" binding:"required" example:"1"`
	Code     string `json:"code" binding:"required" example:"12345"`
}

// Auth Handlers

// Login godoc
// @Summary Вход пользователя
// @Description Аутентифицирует пользователя и возвращает JWT токен.
// @Tags auth
// @Accept  json
// @Produce  json
// @Param   credentials body LoginRequest true "Учетные данные для входа"
// @Success 200 {object} map[string]interface{} "Успешный вход"
// @Failure 400 {object} map[string]interface{} "Неверный запрос"
// @Failure 401 {object} map[string]interface{} "Неверные учетные данные"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /auth/login [post]
func Login(c *gin.Context, db *gorm.DB, cfg *config.Config) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := db.First(&user, "identifier = ?", req.Identifier).Error; err != nil {
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

	c.JSON(http.StatusOK, gin.H{"token": token, "user": user})
}

// Register godoc
// @Summary Регистрация нового студента
// @Description Создает нового пользователя-студента.
// @Tags auth
// @Accept  json
// @Produce  json
// @Param   user body RegisterRequest true "Данные для регистрации пользователя"
// @Success 200 {object} map[string]interface{} "Пользователь успешно зарегистрирован"
// @Failure 400 {object} map[string]interface{} "Неверный запрос"
// @Failure 409 {object} map[string]interface{} "Конфликт (идентификатор или email уже существует)"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /auth/register [post]
func Register(c *gin.Context, db *gorm.DB, cfg *config.Config) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 14)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Ensure group exists
	var group models.Group
	if err := db.First(&group, req.GroupID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	user := models.User{
		Identifier: req.Identifier,
		Password:   string(hashedPassword),
		Name:       req.Name,
		Email:      req.Email,
		Role:       "student",
		GroupID:    &req.GroupID,
	}

	if err := db.Create(&user).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			c.JSON(http.StatusConflict, gin.H{"error": "Identifier or email already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

// GetGroups godoc
// @Summary Получить все группы
// @Description Возвращает список всех студенческих групп.
// @Tags public
// @Produce  json
// @Success 200 {array} models.Group "Список групп"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /groups [get]
func GetGroups(c *gin.Context, db *gorm.DB) {
	var groups []models.Group
	if err := db.Find(&groups).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve groups"})
		return
	}
	c.JSON(http.StatusOK, groups)
}

// Lesson Handlers

// GetLessons godoc
// @Summary Получить занятия
// @Description Возвращает список занятий. Для студентов - занятия их группы. Для преподавателей/администраторов - все занятия.
// @Tags lessons
// @Produce  json
// @Security BearerAuth
// @Success 200 {array} models.Lesson "Список занятий"
// @Failure 404 {object} map[string]interface{} "Пользователь не найден"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /api/lessons [get]
func GetLessons(c *gin.Context, db *gorm.DB) {
	userRole, _ := c.Get("userRole")
	userID, _ := c.Get("userID")

	var lessons []models.Lesson

	if userRole == "student" {
		var currentUser models.User
		if err := db.First(&currentUser, userID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		if currentUser.GroupID == nil {
			c.JSON(http.StatusOK, []models.Lesson{}) // Student has no group, return empty
			return
		}

		// Fetch lessons for the student's group
		if err := db.Preload("Group").Where("group_id = ?", currentUser.GroupID).Find(&lessons).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve lessons for group"})
			return
		}
	} else {
		// For teachers and admins, fetch all lessons
		if err := db.Preload("Group").Order("day, time").Find(&lessons).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve all lessons"})
			return
		}
	}

	c.JSON(http.StatusOK, lessons)
}

// Attendance Handlers

// SubmitAttendance godoc
// @Summary Отметить посещаемость
// @Description Студент отправляет код посещаемости для определенного занятия.
// @Tags student
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param   attendance body SubmitAttendanceRequest true "Данные для отметки посещаемости"
// @Success 200 {object} map[string]interface{} "Посещаемость успешно отмечена"
// @Failure 400 {object} map[string]interface{} "Неверный или просроченный код"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /api/student/attendance [post]
func SubmitAttendance(c *gin.Context, db *gorm.DB) {
	var req SubmitAttendanceRequest
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

// GetStudentAttendance godoc
// @Summary Получить записи о посещаемости студента
// @Description Получает все записи о посещаемости для залогиненного студента.
// @Tags student
// @Produce  json
// @Security BearerAuth
// @Success 200 {array} models.Attendance "Список записей о посещаемости"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /api/student/attendance [get]
func GetStudentAttendance(c *gin.Context, db *gorm.DB) {
	userID, _ := c.Get("userID")
	var attendance []models.Attendance
	if err := db.Preload("Lesson").Where("student_id = ?", uint(userID.(float64))).Find(&attendance).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve attendance records"})
		return
	}
	c.JSON(http.StatusOK, attendance)
}

// GetLessonAttendance godoc
// @Summary Получить посещаемость занятия
// @Description Получает все записи о посещаемости для определенного занятия.
// @Tags teacher
// @Produce  json
// @Security BearerAuth
// @Param lessonId path int true "ID Занятия"
// @Success 200 {array} models.Attendance "Список записей о посещаемости"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /api/teacher/attendance/{lessonId} [get]
func GetLessonAttendance(c *gin.Context, db *gorm.DB) {
	lessonID := c.Param("lessonId")
	var attendance []models.Attendance
	if err := db.Preload("Student").Where("lesson_id = ?", lessonID).Find(&attendance).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve attendance records"})
		return
	}
	c.JSON(http.StatusOK, attendance)
}

// Admin Handlers

// AdminGetUsers godoc
// @Summary Получить всех пользователей (Админ)
// @Description Получает список всех пользователей.
// @Tags admin
// @Produce  json
// @Security BearerAuth
// @Success 200 {array} models.User "Список пользователей"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /api/admin/users [get]
func AdminGetUsers(c *gin.Context, db *gorm.DB) {
	var users []models.User
	if err := db.Preload("Group").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

// AdminCreateUser godoc
// @Summary Создать пользователя (Админ)
// @Description Создает нового пользователя с указанными данными.
// @Tags admin
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param user body models.User true "Объект пользователя"
// @Success 200 {object} models.User "Созданный пользователь"
// @Failure 400 {object} map[string]interface{} "Неверный запрос"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /api/admin/users [post]
func AdminCreateUser(c *gin.Context, db *gorm.DB) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = string(hashedPassword)

	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// AdminUpdateUser godoc
// @Summary Обновить пользователя (Админ)
// @Description Обновляет данные существующего пользователя.
// @Tags admin
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param id path int true "ID Пользователя"
// @Param user body models.User true "Объект пользователя"
// @Success 200 {object} models.User "Обновленный пользователь"
// @Failure 400 {object} map[string]interface{} "Неверный запрос"
// @Failure 404 {object} map[string]interface{} "Пользователь не найден"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /api/admin/users/{id} [put]
func AdminUpdateUser(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var user models.User
	if err := db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// AdminDeleteUser godoc
// @Summary Удалить пользователя (Админ)
// @Description Удаляет пользователя по его ID.
// @Tags admin
// @Produce  json
// @Security BearerAuth
// @Param id path int true "ID Пользователя"
// @Success 200 {object} map[string]interface{} "Пользователь успешно удален"
// @Failure 404 {object} map[string]interface{} "Пользователь не найден"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /api/admin/users/{id} [delete]
func AdminDeleteUser(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	if err := db.Delete(&models.User{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// AdminGetGroups godoc
// @Summary Получить все группы (Админ)
// @Description Получает список всех групп.
// @Tags admin
// @Produce  json
// @Security BearerAuth
// @Success 200 {array} models.Group "Список групп"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /api/admin/groups [get]
func AdminGetGroups(c *gin.Context, db *gorm.DB) {
	var groups []models.Group
	if err := db.Find(&groups).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve groups"})
		return
	}
	c.JSON(http.StatusOK, groups)
}

// Code Handlers

// GenerateCode godoc
// @Summary Сгенерировать код посещаемости
// @Description Генерирует новый 5-значный код для занятия, который истекает через 5 минут.
// @Tags teacher
// @Produce  json
// @Security BearerAuth
// @Param lessonId path int true "ID Занятия"
// @Success 200 {object} models.GeneratedCode "Сгенерированный код"
// @Failure 400 {object} map[string]interface{} "Неверный запрос"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /api/teacher/lessons/{lessonId}/code [post]
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

// DeactivateCode godoc
// @Summary Деактивировать код посещаемости
// @Description Деактивирует активный код для занятия.
// @Tags teacher
// @Produce  json
// @Security BearerAuth
// @Param lessonId path int true "ID Занятия"
// @Success 200 {object} map[string]interface{} "Код успешно деактивирован"
// @Failure 400 {object} map[string]interface{} "Неверный запрос"
// @Failure 404 {object} map[string]interface{} "Активный код не найден"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /api/teacher/lessons/{lessonId}/code [delete]
func DeactivateCode(c *gin.Context, db *gorm.DB) {
	var req struct {
		LessonID uint `json:"lesson_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	result := db.Model(&models.GeneratedCode{}).
		Where("lesson_id = ? AND is_active = ?", req.LessonID, true).
		Update("is_active", false)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to deactivate code"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No active code found for this lesson"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Code deactivated successfully"})
}
