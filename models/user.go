package models

import "time"

// User model
type User struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	Email       string    `gorm:"not null;unique" json:"email"`
	Password    string    `json:"password" binding:"required" gorm:"not null"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Resume      string    `gorm:"type:varchar(255)" json:"resume,omitempty" validate:"omitempty,endswith=.pdf"` // URL/path to PDF resume
	Skills      string    `gorm:"type:text" json:"skills,omitempty"`
	Experience  string    `gorm:"type:text" json:"experience,omitempty"`
	PhoneNumber string    `json:"phone_number,omitempty"`
	LinkedIn    string    `json:"linkedin,omitempty"`
}
