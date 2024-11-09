package routes

import (
	"github.com/Yash-Khattar/HireWiz-Backend/controllers"
	"github.com/Yash-Khattar/HireWiz-Backend/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func JobRouter(router *gin.Engine, db *gorm.DB) {
	jobController := controllers.JobController{DB: db}

	// Public routes (no auth required)
	router.GET("/jobs/public", jobController.GetAllJobs)

	// Protected routes (require company auth)
	jobs := router.Group("/jobs")
	jobs.Use(middleware.CompanyAuthMiddleware())
	{
		jobs.POST("/create", jobController.CreateJob)
		jobs.GET("/getjobs", jobController.GetJobs)
		jobs.GET("/getbyid/:id", jobController.GetJob)
		jobs.PUT("/update/:id", jobController.UpdateJob)
		jobs.DELETE("/delete/:id", jobController.DeleteJob)
	}
}