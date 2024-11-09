package routes

import (
	controller "github.com/Yash-Khattar/HireWiz-Backend/controllers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InterviewRouter(incomingInterviewRoutes *gin.Engine, db *gorm.DB) {
	interviewController := controller.InterviewController{DB: db}
	incomingInterviewRoutes.POST("initiate", interviewController.InitInterview)
}
