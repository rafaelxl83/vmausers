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

// @BasePath /api/v1

// RegisterUser godoc
// @Summary Endpoint to register a new user
// @Schemes
// @Description Add a new user
// @Tags user
// @Accept json
// @Produce json
// @Param user formData models.User true "User Data"
// @Success 200 {object} user
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /user/register [post]
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

// GetUserById godoc
// @Summary Endpoint to load an user by it's ID
// @Schemes
// @Description Get an user
// @Tags user
// @Accept json
// @Produce json
// @Param id path string true "User Id"
// @Success 200 {object} models.User
// @Failure 400 {object} httputil.HTTPError
func GetUserById(context *gin.Context) {
	user, err := middlewares.GetUserById(context.Param("id"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err})
		context.Abort()
		return
	}

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
// @Failure 400 {object} httputil.HTTPError
// @Router /secured/user/{email} [get]
func GetUserByEmail(context *gin.Context) {
	user, err := middlewares.GetUserByEmail(context.Param("email"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err})
		context.Abort()
		return
	}

	context.JSON(http.StatusOK, gin.H{"user": user})
}

// GetManyUsers godoc
// @Summary Endpoint to load a list of users limited to 100
// @Schemes
// @Description Get a list of users
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} models.User
// @Failure 400 {object} httputil.HTTPError
// @Router /secured/user [get]
func GetManyUsers(context *gin.Context) {
	usersList, err := middlewares.GetManyUsers()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err})
		context.Abort()
		return
	}

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
// @Failure 204 {object} httputil.HTTPError
// @Router /secured/user/update [put]
func UpdateUser(context *gin.Context) {
	var user models.User

	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusNoContent, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	currUser, err := middlewares.GetUserByEmail(user.Email)
	currUser.UpdateValues(user)

	if err = middlewares.UpdateUser(currUser); err != nil {
		context.JSON(http.StatusNoContent, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	context.JSON(http.StatusOK, gin.H{"user": user})
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
// @Failure 400 {object} httputil.HTTPError
// @Router /secured/user/update/password [put]
func UpdateUserPassword(context *gin.Context) {
	user, err := middlewares.GetUserByEmail(context.Query("email"))
	if err != nil {
		context.JSON(http.StatusNoContent, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	if err := user.Password.ValidatePasswordRestrictions(context.Param("password")); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	newPassword := models.NewPassword(context.Query("password"))
	if err := middlewares.UpdateUserPassword(user, *newPassword); err != nil {
		context.JSON(http.StatusNoContent, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

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
// @Success 200 {string} user
// @Failure 400 {object} httputil.HTTPError
// @Router /secured/user/update/email [put]
func UpdateUserEmail(context *gin.Context) {
	user, err := middlewares.GetUserByEmail(context.Query("email"))
	if err != nil {
		context.JSON(http.StatusNoContent, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	if !models.IsEmailValid(context.Query("newemail")) {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email"})
		context.Abort()
		return
	}

	err = middlewares.UpdateUserEmail(user, context.Param("newemail"))
	if err != nil {
		context.JSON(http.StatusNoContent, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

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
// @Failure 400 {object} httputil.HTTPError
// @Router /secured/user/{id} [delete]
func DeleteUserById(context *gin.Context) {
	err := middlewares.DeleteUserById(context.Param("id"))
	if err != nil {
		context.JSON(http.StatusNoContent, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

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
// @Failure 400 {object} httputil.HTTPError
// @Router /user/register [delete]
func DeleteUserByEmail(context *gin.Context) {
	err := middlewares.DeleteUserByEmail(context.Param("email"))
	if err != nil {
		context.JSON(http.StatusNoContent, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	context.JSON(http.StatusOK, gin.H{})
}
