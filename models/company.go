package models

import "time"

// Company model
type Company struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"not null;unique" json:"name"`
	Description string    `json:"description"`
	Email       string    `gorm:"not null;unique" json:"email" binding:"required,email"`
	Password    string    `json:"password" binding:"required" gorm:"not null"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Website     string    `json:"website"`
	JobPosts    []JobPost `gorm:"foreignKey:CompanyID" json:"job_posts,omitempty"`
}
