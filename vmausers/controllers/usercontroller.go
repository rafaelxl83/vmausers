package controllers

import (
	"net/http"
	"strings"
	"vmausers/middlewares"
	"vmausers/models"

	"github.com/gin-gonic/gin"
)

func ValidatePasswordRestrictions(providedPassword string) error {
	pass := models.Password{}
	err := pass.ValidatePasswordRestrictions(providedPassword)
	if err != nil {
		return err
	}

	return nil
}

func CheckPassword(user models.User, providedPassword string) error {
	err := user.Password.CheckPassword(providedPassword)
	if err != nil {
		return err
	}

	return nil
}

func RegisterUser(context *gin.Context) {
	var user models.User

	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	if found, err := middlewares.GetUserByEmail(user.Email); err != nil || found != nil {
		errMessage := "User already registered"
		if err != nil {
			errMessage = err.Error()
		}

		context.JSON(http.StatusBadRequest, gin.H{"error": errMessage})
		context.Abort()
		return
	}

	if err := user.Password.ValidatePasswordRestrictions(""); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	user.Password = *models.NewPassword(user.Password.EncryptedPass)
	if err := middlewares.CreateUser(&user); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	context.JSON(http.StatusCreated, gin.H{"userId": user.ID, "email": user.Email, "username": user.Email[:strings.IndexByte(user.Email, '@')]})
}
