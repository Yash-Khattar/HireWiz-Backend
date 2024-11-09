package controllers

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"

	"github.com/Yash-Khattar/HireWiz-Backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ApplicationController struct {
	DB *gorm.DB
}

func (a *ApplicationController) uploadToGCS(file *multipart.FileHeader) (string, error) {
	ctx := context.Background()
	fmt.Println("Uploading to GCS")

	// Set up credentials from the JSON file
	credentialsFile := "ai-hr-441207-c32baa40bcf7.json"
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(credentialsFile))
	if err != nil {
		return "", fmt.Errorf("failed to create client: %v", err)
	}
	fmt.Println("Client created")
	defer client.Close()
	fmt.Println("Client closed")
	bucketName := "hire-wiz" // Your bucket name
	bucket := client.Bucket(bucketName)
	fmt.Println("Bucket created")
	// Generate a unique filename
	filename := fmt.Sprintf("resumes/%d-%s", time.Now().Unix(), file.Filename)
	fmt.Println("Filename created")
	obj := bucket.Object(filename)
	fmt.Println("Object created")
	writer := obj.NewWriter(ctx)
	fmt.Println("Writer created")

	// Open the uploaded file
	src, err := file.Open()
	fmt.Println("Source opened")
	if err != nil {
		return "", err
	}
	defer src.Close()
	fmt.Println("Source closed")

	// Copy the file content to GCS
	if _, err := io.Copy(writer, src); err != nil {
		return "", err
	}
	fmt.Println("Copy completed")
	if err := writer.Close(); err != nil {
		return "", err
	}
	fmt.Println("Writer closed")
	// Generate the public URL
	publicURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, filename)
	fmt.Println("Public URL generated" + publicURL)
	return publicURL, nil
}

func (a *ApplicationController) ApplyForJob(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Type assert to uint
	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	if userIDUint == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get job ID from URL
	jobID := c.Param("id")

	// Check if job exists
	var job models.JobPost
	if err := a.DB.First(&job, jobID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}

	// Check if already applied
	var existingApplication models.Application
	if err := a.DB.Where("job_post_id = ? AND user_id = ?", jobID, userIDUint).First(&existingApplication).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Already applied for this job"})
		return
	}

	// Get resume file
	file, err := c.FormFile("resume")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Resume file is required"})
		return
	}

	// Upload resume to GCS instead of Cloudinary
	resumeURL, err := a.uploadToGCS(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload resume"})
		return
	}

	int_jobID, err := strconv.Atoi(jobID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID"})
		return
	}
	// Create application
	application := models.Application{
		JobPostID: uint(int_jobID),
		UserID:    userIDUint,
		AppliedAt: time.Now(),
		Status:    "pending",
		ResumeURL: resumeURL,
	}

	if err := a.DB.Create(&application).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create application"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":     "Application submitted successfully",
		"application": application,
	})
}
