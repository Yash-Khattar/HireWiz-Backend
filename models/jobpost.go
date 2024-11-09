package models

import "time"

// JobPost model
type JobPost struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"not null" json:"title" binding:"required"`
	Description string    `gorm:"type:text" json:"description" binding:"required"`
	Location    string    `json:"location"`
	CompanyID   uint      `json:"company_id"`
	Company     Company   `json:"company" gorm:"foreignKey:CompanyID"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	// Many-to-many relationship with Users through Applications
	Applicants []User `gorm:"many2many:applications;" json:"applicants,omitempty"`
}

// Application model to store additional application data
type Application struct {
	JobPostID uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"primaryKey"`
	AppliedAt time.Time `json:"applied_at"`
	Status    string    `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending, accepted, rejected
	ResumeURL string    `json:"resume_url,omitempty"`
}
