package database

import (
	"github.com/Yash-Khattar/HireWiz-Backend/models"
	"gorm.io/gorm"
)

func DBMigrator(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Company{},
		&models.User{},
		&models.JobPost{},
		&models.Application{},
	)
}