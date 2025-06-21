package database

import (
	"fmt"
	"log"
	"student-attendance-app/pkg/models"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"
)

func InitDB(url string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Run migrations
	if err := db.AutoMigrate(&models.User{}, &models.Lesson{}, &models.Attendance{}, &models.GeneratedCode{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	// Seed the database
	seedDatabase(db)

	return db, nil
}

func seedDatabase(db *gorm.DB) {
	// Seed users
	seedUsers(db)
	// Seed lessons
	seedLessons(db)
}

func seedUsers(db *gorm.DB) {
	users := []models.User{
		{StudentID: "stu001", Name: "Иванов И.И.", Role: "student", Password: hashPassword("12345")},
		{StudentID: "stu002", Name: "Петрова А.А.", Role: "student", Password: hashPassword("23456")},
		{StudentID: "teach001", Name: "Преподаватель С.С.", Role: "teacher", Password: hashPassword("admin1")},
	}

	for _, user := range users {
		var existingUser models.User
		if db.First(&existingUser, "student_id = ?", user.StudentID).Error == gorm.ErrRecordNotFound {
			if err := db.Create(&user).Error; err != nil {
				log.Printf("failed to seed user %s: %v", user.StudentID, err)
			}
		}
	}
}

func seedLessons(db *gorm.DB) {
	lessons := []models.Lesson{
		{Name: "Математика", Day: "Понедельник", Time: "09:00-10:30", Teacher: "Преподаватель С.С.", Room: "101"},
		{Name: "Физика", Day: "Понедельник", Time: "10:45-12:15", Teacher: "Преподаватель С.С.", Room: "102"},
		{Name: "Химия", Day: "Вторник", Time: "09:00-10:30", Teacher: "Преподаватель С.С.", Room: "103"},
		{Name: "Биология", Day: "Вторник", Time: "10:45-12:15", Teacher: "Преподаватель С.С.", Room: "104"},
		{Name: "История", Day: "Среда", Time: "09:00-10:30", Teacher: "Преподаватель С.С.", Room: "105"},
	}

	for _, lesson := range lessons {
		var existingLesson models.Lesson
		if db.First(&existingLesson, "name = ? AND day = ? AND time = ?", lesson.Name, lesson.Day, lesson.Time).Error == gorm.ErrRecordNotFound {
			if err := db.Create(&lesson).Error; err != nil {
				log.Printf("failed to seed lesson %s: %v", lesson.Name, err)
			}
		}
	}
}

func hashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
} 