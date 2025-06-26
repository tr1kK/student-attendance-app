package database

import (
	"fmt"
	"log"
	"student-attendance-app/pkg/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(url string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Run migrations
	if err := db.AutoMigrate(
		&models.Group{},
		&models.User{},
		&models.Lesson{},
		&models.Attendance{},
		&models.GeneratedCode{},
	); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	seedDatabase(db)

	return db, nil
}

func seedDatabase(db *gorm.DB) {
	// Seed Groups if they don't exist
	groups := []models.Group{
		{Name: "Group A"},
		{Name: "Group B"},
		{Name: "Group C"},
	}
	for _, group := range groups {
		db.FirstOrCreate(&group, models.Group{Name: group.Name})
	}
	log.Println("Groups seeded.")

	// Seed Users
	seedUsers(db)

	// Seed Lessons
	seedLessons(db)
}

func seedUsers(db *gorm.DB) {
	users := []struct {
		Identifier string
		Password   string // Plaintext password
		Name       string
		Email      string
		Role       string
		GroupName  string // Group name to find ID
	}{
		// Password: 12345
		{Identifier: "student001", Password: "12345", Name: "Test Student", Email: "student@test.com", Role: "student", GroupName: "Group A"},
		// Password: admin1
		{Identifier: "teacher001", Password: "admin1", Name: "Test Teacher", Email: "teacher@test.com", Role: "teacher", GroupName: "Group B"},
		// Password: rootpass
		{Identifier: "admin001", Password: "rootpass", Name: "Test Admin", Email: "admin@test.com", Role: "admin", GroupName: "Group C"},
	}

	for _, u := range users {
		// Check if user already exists
		var existingUser models.User
		if db.First(&existingUser, "identifier = ?", u.Identifier).Error == gorm.ErrRecordNotFound {
			// Find group to get its ID
			var group models.Group
			if err := db.First(&group, "name = ?", u.GroupName).Error; err != nil {
				log.Printf("Could not find group %s for user %s: %v", u.GroupName, u.Identifier, err)
				continue
			}

			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
			if err != nil {
				log.Printf("Failed to hash password for user %s: %v", u.Identifier, err)
				continue
			}

			user := models.User{
				Identifier: u.Identifier,
				Password:   string(hashedPassword),
				Name:       u.Name,
				Email:      u.Email,
				Role:       u.Role,
				GroupID:    &group.ID,
			}
			if err := db.Create(&user).Error; err != nil {
				log.Printf("Failed to create user %s: %v", u.Identifier, err)
			} else {
				log.Printf("User %s created.", u.Identifier)
			}
		}
	}
}

func seedLessons(db *gorm.DB) {
	// Get group IDs
	var groupA, groupB, groupC models.Group
	db.First(&groupA, "name = ?", "Group A")
	db.First(&groupB, "name = ?", "Group B")
	db.First(&groupC, "name = ?", "Group C")

	lessons := []models.Lesson{
		// Group A
		{Name: "Алгебра", Day: "Понедельник", Time: "09:00-10:30", Teacher: "Анна Владимировна", Room: "101", GroupID: &groupA.ID},
		{Name: "История", Day: "Понедельник", Time: "12:30-14:00", Teacher: "Иван Петрович", Room: "203", GroupID: &groupA.ID},
		{Name: "Геометрия", Day: "Вторник", Time: "10:45-12:15", Teacher: "Анна Владимировна", Room: "101", GroupID: &groupA.ID},
		// Group B
		{Name: "Физика", Day: "Среда", Time: "09:00-10:30", Teacher: "Петр Сидорович", Room: "203", GroupID: &groupB.ID},
		{Name: "Химия", Day: "Четверг", Time: "12:30-14:00", Teacher: "Мария Ивановна", Room: "305", GroupID: &groupB.ID},
		// Group C
		{Name: "Информатика", Day: "Пятница", Time: "10:45-12:15", Teacher: "Сергей Николаевич", Room: "404", GroupID: &groupC.ID},
	}

	for _, lesson := range lessons {
		db.FirstOrCreate(&lesson, models.Lesson{Name: lesson.Name, Day: lesson.Day, Time: lesson.Time})
	}
	log.Println("Lessons seeded.")
}

func hashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}
