package models

import "time"

type Group struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"unique;not null" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type User struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Identifier string    `gorm:"unique;not null" json:"identifier"`
	Password   string    `gorm:"not null" json:"-"` // Omit from JSON responses
	Name       string    `gorm:"not null" json:"name"`
	Email      string    `gorm:"unique;not null" json:"email"`
	Role       string    `gorm:"not null" json:"role"` // 'student', 'teacher', or 'admin'
	GroupID    *uint     `json:"group_id"`
	Group      Group     `gorm:"foreignKey:GroupID;references:ID" json:"group"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Lesson struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"uniqueIndex:idx_lesson_name_day_time" json:"name"`
	Day       string    `gorm:"uniqueIndex:idx_lesson_name_day_time" json:"day"`
	Time      string    `gorm:"uniqueIndex:idx_lesson_name_day_time" json:"time"`
	Teacher   string    `json:"teacher"`
	Room      string    `json:"room"`
	GroupID   *uint     `json:"group_id"`
	Group     Group     `json:"group"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Attendance struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	LessonID    uint      `gorm:"not null" json:"lesson_id"`
	StudentID   uint      `gorm:"not null" json:"student_id"`
	SubmittedAt time.Time `gorm:"not null" json:"submitted_at"`
	Lesson      Lesson    `gorm:"foreignKey:LessonID;references:ID" json:"lesson"`
	Student     User      `gorm:"foreignKey:StudentID;references:ID" json:"student"`
}

type GeneratedCode struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	LessonID  uint      `gorm:"not null" json:"lesson_id"`
	Code      string    `gorm:"not null" json:"code"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	IsActive  bool      `gorm:"not null;default:true" json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	Lesson    Lesson    `gorm:"foreignKey:LessonID;references:ID" json:"lesson"`
}
