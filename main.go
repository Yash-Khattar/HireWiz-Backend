package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Yash-Khattar/HireWiz-Backend/database"
	"github.com/Yash-Khattar/HireWiz-Backend/handlers"
	"github.com/Yash-Khattar/HireWiz-Backend/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func keepAlive(baseURL string) {
	client := &http.Client{}
	for {
		resp, err := client.Get(baseURL)
		if err != nil {
			log.Printf("Keep-alive request failed: %v", err)
		} else {
			resp.Body.Close()
			log.Println("Keep-alive request successful")
		}
		time.Sleep(14 * time.Minute)
	}
}

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

	go keepAlive("https://hirewiz-backend.onrender.com:" + port)

	router := gin.New()
	router.Use(gin.Logger())

	// Add CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},  // Allow all origins
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Good morning !",
		})
	})

	db := database.DBinit()
	routes.AuthRouter(router, db)
	routes.JobRouter(router, db)
	routes.InterviewRouter(router, db)

	// Initialize Cloudinary
	if err := handlers.InitCloudinary(); err != nil {
		log.Fatal("Failed to initialize Cloudinary:", err)
	}

	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
