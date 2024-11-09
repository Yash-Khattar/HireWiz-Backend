package middleware

import (
	"strings"

	"github.com/Yash-Khattar/HireWiz-Backend/handlers"
	"github.com/gin-gonic/gin"
)

func CompanyAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(401, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		token := parts[1]
		companyID, err := handlers.ValidateJWT(token)
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Set company ID in context for use in handlers
		c.Set("company_id", companyID)
		c.Next()
	}
}