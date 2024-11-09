package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Yash-Khattar/HireWiz-Backend/database"
	"github.com/Yash-Khattar/HireWiz-Backend/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/Yash-Khattar/HireWiz-Backend/handlers"
)

func main() {
	fmt.Println("Welcome to HireWiz!! 🚀🚀🚀")

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("err loading env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := gin.New()
	router.Use(gin.Logger())

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Good morning !",
		})
	})

	db := database.DBinit()
	routes.AuthRouter(router, db)
	routes.JobRouter(router, db)

	// Initialize Cloudinary
	if err := handlers.InitCloudinary(); err != nil {
		log.Fatal("Failed to initialize Cloudinary:", err)
	}

	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
