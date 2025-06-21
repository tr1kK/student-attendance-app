package models

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	StudentID string    `gorm:"unique;not null" json:"student_id"`
	Password  string    `gorm:"not null" json:"-"` // Omit from JSON responses
	Name      string    `gorm:"not null" json:"name"`
	Role      string    `gorm:"not null" json:"role"` // 'student' or 'teacher'
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Lesson struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Day       string    `gorm:"not null" json:"day"`
	Time      string    `gorm:"not null" json:"time"`
	Teacher   string    `gorm:"not null" json:"teacher"`
	Room      string    `gorm:"not null" json:"room"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Attendance struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	LessonID    uint      `gorm:"not null" json:"lesson_id"`
	StudentID   uint      `gorm:"not null" json:"student_id"`
	SubmittedAt time.Time `gorm:"not null" json:"submitted_at"`
	Lesson      Lesson    `gorm:"foreignKey:LessonID" json:"lesson"`
	Student     User      `gorm:"foreignKey:StudentID" json:"student"`
}

type GeneratedCode struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	LessonID  uint      `gorm:"not null" json:"lesson_id"`
	Code      string    `gorm:"not null" json:"code"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	IsActive  bool      `gorm:"not null;default:true" json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	Lesson    Lesson    `gorm:"foreignKey:LessonID" json:"lesson"`
} 