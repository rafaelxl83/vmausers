package controllers

import (
	"net/http"
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

	if err := user.ValidatePasswordRestrictions(""); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	if err := middlewares.CreateUser(&user); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	context.JSON(http.StatusCreated, gin.H{"userId": user.ID, "email": user.Email, "username": user.FirstName})
}
