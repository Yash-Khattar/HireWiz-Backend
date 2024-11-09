package main

import (
	"fmt"
	"log"
	"os"
	"net/http"
	"time"

	"github.com/Yash-Khattar/HireWiz-Backend/database"
	"github.com/Yash-Khattar/HireWiz-Backend/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/Yash-Khattar/HireWiz-Backend/handlers"
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

	go keepAlive("http://localhost:" + port)

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
