package controllers

import (
	"net/http"
	"strings"
	"vmausers/auth"
	"vmausers/middlewares"

	"github.com/gin-gonic/gin"
)

type TokenRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// @BasePath /api/v1

// GenerateToken godoc
// @Summary Endpoint to generate a new JWT token
// @Schemes
// @Description Generates a new JWT token
// @Tags token
// @Accept json
// @Produce json
// @Param user formData models.User true "User Data"
// @Success 200 {string} user
// @Router /user/register [post]
func GenerateToken(context *gin.Context) {
	var request TokenRequest

	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	// check if email exists and password is correct
	user, err := middlewares.GetUserByEmail(request.Email)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	credentialError := user.Password.CheckPassword(request.Password)
	if credentialError != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		context.Abort()
		return
	}

	tokenString, err := auth.GenerateJWT(user.Email, user.Email[:strings.IndexByte(user.Email, '@')])
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	//http://jwt.io/
	context.JSON(http.StatusOK, gin.H{"token": tokenString})
}
