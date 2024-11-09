package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/Yash-Khattar/HireWiz-Backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type InterviewController struct {
	DB *gorm.DB
}

func (i *InterviewController) InitInterview(c *gin.Context) {
	// Parse request body
	var input struct {
		SessionID string `json:"session_id" binding:"required"`
		UserID    uint   `json:"user_id" binding:"required"`
		JobID     uint   `json:"job_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fetch job description
	var job models.JobPost
	if err := i.DB.First(&job, input.JobID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}

	// Fetch user's resume URL from applications
	var application models.Application
	if err := i.DB.Debug().Where("job_post_id = ? AND user_id = ?", input.JobID, input.UserID).First(&application).Error; err != nil {
		// Log the actual error and input values
		fmt.Printf("Failed to find application - JobID: %d, UserID: %d, Error: %v\n", input.JobID, input.UserID, err)

		// Check if it's specifically a record not found error
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "No application found for this job and user combination",
				"details": fmt.Sprintf("JobID: %d, UserID: %d", input.JobID, input.UserID),
			})
			return
		}

		// Handle other database errors
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error while fetching application"})
		return
	}

	// Download resume from Cloudinary URL
	resp, err := http.Get(application.ResumeURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch resume"})
		return
	}
	defer resp.Body.Close()

	// Create multipart form data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add session_id field
	if err := writer.WriteField("session_id", input.SessionID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	// Add JD field
	if err := writer.WriteField("JD", job.Description); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	// Add resume file
	part, err := writer.CreateFormFile("resume", "resume.pdf")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}
	if _, err := io.Copy(part, resp.Body); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	writer.Close()

	// Make request to AI service
	req, err := http.NewRequest("POST", "https://ai-hr-4anm.onrender.com/init", body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	aiResp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to make AI service request"})
		return
	}
	defer aiResp.Body.Close()

	// Parse AI service response
	var aiResponse struct {
		Response string `json:"response"`
	}
	if err := json.NewDecoder(aiResp.Body).Decode(&aiResponse); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse AI service response"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Interview initialized successfully",
		"response": aiResponse.Response,
	})
}
