package controllers

import (
	"net/http"
	"strings"
	"vmausers/middlewares"
	"vmausers/models"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
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

// RegisterUser godoc
// @Summary Endpoint to register a new user
// @Schemes
// @Description Add a new user
// @Tags user
// @Accept json
// @Produce json
// @Param user formData models.User true "User Data"
// @Success 200 {object} models.User
// @Failure 204 {string} string "No Content"
// @Failure 400 {string} string "Bad request"
// @Failure 406 {string} string "Not Acceptable"
// @Failure 500 {string} string "Server Error"
// @Router /user/register [post]
func RegisterUser(context *gin.Context) {
	var user models.User

	if err := context.ShouldBindJSON(&user); err != nil {
		log.Errorf("RegisterUser: Bad Request [%v]", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	if !models.IsEmailValid(user.Email) {
		log.Errorf("RegisterUser: Not Acceptable [%v][Invalid email format]", user)
		context.JSON(http.StatusNotAcceptable, gin.H{"error": "Invalid email"})
		context.Abort()
		return
	}

	if found, err := middlewares.GetUserByEmail(user.Email); err == nil || found != nil {
		errMessage := "User already registered"
		if err != nil {
			errMessage = err.Error()
		}

		log.Errorf("RegisterUser: Bad Request [%v][%v]", user, err)
		context.JSON(http.StatusBadRequest, gin.H{"error": errMessage})
		context.Abort()
		return
	}

	if err := user.Password.ValidatePasswordRestrictions(""); err != nil {
		log.Errorf("RegisterUser: Not Acceptable [%v][%v]", user, err)
		context.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	user.Password = *models.NewPassword(user.Password.EncryptedPass)
	if err := middlewares.CreateUser(&user); err != nil {
		log.Errorf("RegisterUser: Internal Server Error [%v][%v]", user, err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	log.Infof("RegisterUser: [%v]", user)
	context.JSON(http.StatusCreated, gin.H{"userId": user.ID, "email": user.Email, "username": user.Email[:strings.IndexByte(user.Email, '@')]})
}

func GetUserById(context *gin.Context) {
	user, err := middlewares.GetUserById(context.Param("id"))
	if err != nil {
		log.Errorf("GetUserById: Bad Request [%v]", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": err})
		context.Abort()
		return
	}

	log.Infof("GetUserById: [%v]", user)
	context.JSON(http.StatusOK, gin.H{"user": user})
}

// GetUserByEmail godoc
// @Summary Endpoint to load an user by it's email
// @Schemes
// @Description Get an user
// @Tags user
// @Accept json
// @Produce json
// @Param email path string true "User Email"
// @Success 200 {object} models.User
// @Failure 400 {string} string "Bad request"
// @Router /secured/user/{email} [get]
func GetUserByEmail(context *gin.Context) {
	user, err := middlewares.GetUserByEmail(context.Param("email"))
	if err != nil {
		log.Errorf("GetUserByEmail: Bad Request [%v]", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": err})
		context.Abort()
		return
	}

	log.Infof("GetUserByEmail: [%v]", user)
	context.JSON(http.StatusOK, gin.H{"user": user})
}

// GetManyUsers godoc
// @Summary Endpoint to load a list of users limited to 100
// @Schemes
// @Description Get a list of users
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {array} models.User
// @Failure 400 {string} string "Bad request"
// @Router /secured/user [get]
func GetManyUsers(context *gin.Context) {
	usersList, err := middlewares.GetManyUsers()
	if err != nil {
		log.Errorf("GetManyUsers: Bad Request [%v]", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": err})
		context.Abort()
		return
	}

	log.Infof("GetManyUsers: [%v]", len(*usersList))
	context.JSON(http.StatusOK, gin.H{"users": usersList})
}

// UpdateUser godoc
// @Summary Endpoint to update common user information
// @Schemes
// @Description Update user information
// @Tags user
// @Accept json
// @Produce json
// @Param user formData models.User true "User Data"
// @Success 200 {object} models.User
// @Failure 204 {string} string "No Content"
// @Router /secured/user/update [put]
func UpdateUser(context *gin.Context) {
	var user models.User

	if err := context.ShouldBindJSON(&user); err != nil {
		log.Errorf("UpdateUser: No Content [%v]", err)
		context.JSON(http.StatusNoContent, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	currUser, err := middlewares.GetUserByEmail(user.Email)
	currUser.UpdateValues(user)

	if err = middlewares.UpdateUser(currUser); err != nil {
		log.Errorf("UpdateUser: No Content [%v]", err)
		context.JSON(http.StatusNoContent, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	log.Infof("UpdateUser: [%v][%v]", currUser, user)
	context.JSON(http.StatusOK, gin.H{"user": currUser})
}

// UpdateUserPassword godoc
// @Summary Endpoint to update the user password
// @Schemes
// @Description Update the user password
// @Tags user
// @Accept json
// @Produce json
// @Param email query string true "User email"
// @Param password query string true "User password"
// @Success 200
// @Failure 204 {string} string "No Content"
// @Failure 400 {string} string "Bad Request"
// @Router /secured/user/update/password [put]
func UpdateUserPassword(context *gin.Context) {
	user, err := middlewares.GetUserByEmail(context.Query("email"))
	if err != nil {
		log.Errorf("UpdateUserPassword: No Content [%v]", err)
		context.JSON(http.StatusNoContent, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	if err := user.Password.ValidatePasswordRestrictions(context.Param("password")); err != nil {
		log.Errorf("UpdateUserPassword: Bad Request [%v]", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	newPassword := models.NewPassword(context.Query("password"))
	if err := middlewares.UpdateUserPassword(user, *newPassword); err != nil {
		log.Warnf("UpdateUserPassword: No Content [%v]", err)
		context.JSON(http.StatusNoContent, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	log.Infof("UpdateUserPassword: [%v][%v]", user, newPassword)
	context.JSON(http.StatusOK, gin.H{})
}

// UpdateUserEmail godoc
// @Summary Endpoint to register a new user
// @Schemes
// @Description Add a new user
// @Tags user
// @Accept json
// @Produce json
// @Param email query string true "User email"
// @Param newemail query string true "User newemail"
// @Success 200 {object} models.User
// @Failure 204 {string} string "No Content"
// @Failure 400 {string} string "Bad request"
// @Router /secured/user/update/email [put]
func UpdateUserEmail(context *gin.Context) {
	user, err := middlewares.GetUserByEmail(context.Query("email"))
	if err != nil {
		log.Warnf("UpdateUserEmail: No Content [%v]", err)
		context.JSON(http.StatusNoContent, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	newEmail := context.Query("newemail")
	if !models.IsEmailValid(newEmail) {
		log.Warnf("UpdateUserEmail: Bad Request [%v][%v]", newEmail, err)
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email"})
		context.Abort()
		return
	}

	err = middlewares.UpdateUserEmail(user, newEmail)
	if err != nil {
		log.Errorf("UpdateUserEmail: No Content [%v][%v]", newEmail, err)
		context.JSON(http.StatusNoContent, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	log.Infof("UpdateUserEmail: [%v][%v]", user, newEmail)
	context.JSON(http.StatusOK, gin.H{"user": user})
}

// DeleteUserById godoc
// @Summary Endpoint to exclude an user
// @Schemes
// @Description Delete an user
// @Tags user
// @Accept json
// @Produce json
// @Param id path string true "User Id"
// @Success 200
// @Failure 204 {string} string "No Content"
// @Router /secured/user/{id} [delete]
func DeleteUserById(context *gin.Context) {
	err := middlewares.DeleteUserById(context.Param("id"))
	if err != nil {
		log.Errorf("DeleteUserById: No Content [%v]", err)
		context.JSON(http.StatusNoContent, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	log.Infof("DeleteUserById: [%v]", context.Param("id"))
	context.JSON(http.StatusOK, gin.H{})
}

// RegisterUser godoc
// @Summary Endpoint to register a new user
// @Schemes
// @Description Add a new user
// @Tags user
// @Accept json
// @Produce json
// @Param email path string true "User Data"
// @Success 200
// @Failure 204 {string} string "No Content"
// @Router /user/register [delete]
func DeleteUserByEmail(context *gin.Context) {
	err := middlewares.DeleteUserByEmail(context.Param("email"))
	if err != nil {
		log.Errorf("DeleteUserByEmail: No Content [%v]", err)
		context.JSON(http.StatusNoContent, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	log.Infof("DeleteUserByEmail: [%v]", context.Param("email"))
	context.JSON(http.StatusOK, gin.H{})
}
