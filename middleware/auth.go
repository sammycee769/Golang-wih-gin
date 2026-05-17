package middleware

import (
	"net/http"
	"os"
	"strings"
	"todoList/db"
	"todoList/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gin-gonic/gin"
)

func RequireAuth(content *gin.Context) {
	authHeader := content.GetHeader("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		content.AbortWithStatusJSON(http.StatusUnauthorized,gin.H{"error":"missing token"})
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		content.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	claims, _ := token.Claims.(jwt.MapClaims)
	userID := uint(claims["sub"].(float64))

	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		content.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	content.Set("user", user)
	content.Next()
}