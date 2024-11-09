package handlers

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(companyID uint) (string , error){
	claims:= jwt.MapClaims{
		"id" : companyID,
		"exp" : time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

// ValidateJWT validates the JWT token and returns the company ID
func ValidateJWT(tokenString string) (uint, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return 0, err
	}

	// Validate token and extract claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check if token is expired
		exp, ok := claims["exp"].(float64)
		if !ok {
			return 0, errors.New("invalid exp claim")
		}
		if float64(time.Now().Unix()) > exp {
			return 0, errors.New("token expired") 
		}

		// Extract and validate company ID
		if id, ok := claims["id"].(float64); ok {
			return uint(id), nil
		}
		return 0, errors.New("invalid id claim")
	}

	return 0, errors.New("invalid token")
}