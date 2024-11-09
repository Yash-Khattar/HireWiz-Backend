package controllers

import (
	"net/http"
	"strings"
	"time"

	"github.com/Yash-Khattar/HireWiz-Backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type JobController struct {
	DB *gorm.DB
}

// CreateJob handles job creation
func (j *JobController) CreateJob(c *gin.Context) {
	// Get company ID from context (set by auth middleware)
	companyID := c.GetUint("company_id")
	if companyID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Create input struct for validation
	var input struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description" binding:"required"`
		Location    string `json:"location" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create job post
	job := models.JobPost{
		Title:       input.Title,
		Description: input.Description,
		Location:    input.Location,
		CompanyID:   companyID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := j.DB.Create(&job).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create job"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Job created successfully",
		"job":     job,
	})
}

// GetJobs returns all jobs for a company
func (j *JobController) GetJobs(c *gin.Context) {
	companyID := c.GetUint("company_id")
	if companyID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var jobs []models.JobPost
	if err := j.DB.Where("company_id = ?", companyID).Find(&jobs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch jobs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"jobs": jobs})
}

// GetJob returns a specific job
func (j *JobController) GetJob(c *gin.Context) {
	companyID := c.GetUint("company_id")
	jobID := c.Param("id")

	var job models.JobPost
	if err := j.DB.Where("id = ? AND company_id = ?", jobID, companyID).First(&job).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch job"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"job": job})
}

// UpdateJob handles job updates
func (j *JobController) UpdateJob(c *gin.Context) {
	companyID := c.GetUint("company_id")
	jobID := c.Param("id")

	// Check if job exists and belongs to company
	var job models.JobPost
	if err := j.DB.Where("id = ? AND company_id = ?", jobID, companyID).First(&job).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch job"})
		return
	}

	// Bind update data
	var input models.JobPost
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields
	input.UpdatedAt = time.Now()
	if err := j.DB.Model(&job).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update job"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Job updated successfully",
		"job":     job,
	})
}

// DeleteJob handles job deletion
func (j *JobController) DeleteJob(c *gin.Context) {
	companyID := c.GetUint("company_id")
	jobID := c.Param("id")

	// Check if job exists and belongs to company
	var job models.JobPost
	if err := j.DB.Where("id = ? AND company_id = ?", jobID, companyID).First(&job).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch job"})
		return
	}

	// Delete job
	if err := j.DB.Delete(&job).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete job"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Job deleted successfully"})
}

// GetAllJobsPublic returns all jobs for public viewing (users)
func (j *JobController) GetAllJobs(c *gin.Context) {
	var jobs []models.JobPost

	// Create the base query
	query := j.DB.Preload("Company", func(db *gorm.DB) *gorm.DB {
		// Select only the public company fields
		return db.Select("id, name, description, website")
	})

	// Handle search query if provided
	if search := c.Query("search"); search != "" {
		query = query.Where("LOWER(title) LIKE ? OR LOWER(description) LIKE ?",
			"%"+strings.ToLower(search)+"%",
			"%"+strings.ToLower(search)+"%")
	}

	// Handle location filter if provided
	if location := c.Query("location"); location != "" {
		query = query.Where("LOWER(location) = ?", strings.ToLower(location))
	}

	// Execute the query
	if err := query.Find(&jobs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch jobs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"jobs":  jobs,
		"count": len(jobs),
	})
}
