package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func DBinit() *gorm.DB{
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("err loading env file")
	}
	dns := os.Getenv("DNS")
	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Unable to connect to db : %v", err)
	}

	fmt.Println("DB Connected successfully")
	
	if err:= DBMigrator(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	return db


}