package controllers

import (
	"net/http"
	"strings"
	"vmausers/auth"
	"vmausers/middlewares"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
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
		log.Errorf("GenerateToken: Bad Request [%v]", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	// check if email exists and password is correct
	user, err := middlewares.GetUserByEmail(request.Email)
	if err != nil {
		log.Errorf("GenerateToken: Bad Request [%s][%v]", user.Email, err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	credentialError := user.Password.CheckPassword(request.Password)
	if credentialError != nil {
		log.Warnf("GenerateToken: Not Acceptable [%s][%v]", user.Email, err)
		context.JSON(http.StatusNotAcceptable, gin.H{"error": "invalid credentials"})
		context.Abort()
		return
	}

	tokenString, err := auth.GenerateJWT(user.Email, user.Email[:strings.IndexByte(user.Email, '@')])
	if err != nil {
		log.Warnf("GenerateToken: Internal ServerError [%s][%v]", user.Email, err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	log.Debugf("GenerateToken: New token [%s][%v]", user.Email, tokenString)
	//http://jwt.io/
	context.JSON(http.StatusOK, gin.H{"token": tokenString})
}
