package handlers

import(
	"net/http"
	"os"
	"time"
	"todoList/db"
	"todoList/models"

	"github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"

)

func RegisterUser(content *gin.Context) {
	var input models.RegisterInput

	if err := content.ShouldBindJSON(&input); err != nil {
		content.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		content.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
	}
	if err := db.DB.Create(&user).Error; err != nil {
		content.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	content.JSON(http.StatusCreated, gin.H{"data": gin.H{"id": user.ID, "name": user.Name, "email": user.Email}})
}

func LoginUser(content *gin.Context) {
	var input models.LoginInput

	if err := content.ShouldBindJSON(&input); err != nil {
		content.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := db.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		content.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		content.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		content.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	content.JSON(http.StatusOK, gin.H{"token": tokenString})
}