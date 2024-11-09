package routes

import (
	controller "github.com/Yash-Khattar/HireWiz-Backend/controllers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AuthRouter(incomingAuthRoutes *gin.Engine, db *gorm.DB) {
	authController := controller.AuthController{DB: db}
	incomingAuthRoutes.POST("company/register", authController.CompanyRegister)
	incomingAuthRoutes.POST("company/login", authController.CompanyLogin)
	incomingAuthRoutes.POST("user/register", authController.UserRegister)
	incomingAuthRoutes.POST("user/login", authController.UserLogin)
}
