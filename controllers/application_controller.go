package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Yash-Khattar/HireWiz-Backend/handlers"
	"github.com/Yash-Khattar/HireWiz-Backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ApplicationController struct {
	DB *gorm.DB
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

	// Upload resume to Cloudinary
	resumeURL, err := handlers.UploadResume(file)
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
		"message": "Application submitted successfully",
		"application": application,
	})
} 