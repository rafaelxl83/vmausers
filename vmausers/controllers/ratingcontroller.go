package controllers

import (
	"net/http"
	"sort"
	"strconv"
	"vmausers/helper"
	"vmausers/middlewares"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func getRatingProperty(age int) int {
	index := sort.Search(len(helper.AppConfig.RatingBoard.Ratings)-1, func(i int) bool {
		return helper.AppConfig.RatingBoard.Ratings[i+1].MinAge > age
	})
	return index
}

// GetRatingList godoc
// @Summary Endpoint to get the rating list in use
// @Schemes
// @Description Get the rating list
// @Tags rating
// @Accept json
// @Produce json
// @Success 200 {array} helper.Config.RatingBoard.Ratings
// @Router /secured/rating [get]
func GetRatingList(context *gin.Context) {
	log.Infof("GetRatingList: [%v]", len(helper.AppConfig.RatingBoard.Ratings))
	context.JSON(http.StatusOK, gin.H{"rating": helper.AppConfig.RatingBoard.Ratings})
}

// GetAgeRating godoc
// @Summary Endpoint to get the rating
// @Schemes
// @Description Get an rating classification depending of the required age
// @Tags rating
// @Accept json
// @Produce json
// @Param age path string true "An Age"
// @Success 200 {object} helper.Config.RatingBoard.Ratings
// @Failure 204 {string} string "No Content"
// @Failure 400 {string} string "Bad request"
// @Router /secured/rating/byage/{age} [get]
func GetAgeRating(context *gin.Context) {
	age, err := strconv.ParseInt(context.Param("age"), 0, 16)
	if err != nil {
		log.Errorf("GetAgeRating: Bad Request [%v]", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": err})
		context.Abort()
		return
	}

	index := getRatingProperty(int(age))
	if index >= len(helper.AppConfig.RatingBoard.Ratings) {
		log.Errorf("GetAgeRating: No Content [%v]", age)
		context.JSON(http.StatusNoContent, gin.H{"error": err})
		context.Abort()
		return
	}

	log.Infof("GetAgeRating: [%v]", helper.AppConfig.RatingBoard.Ratings[index])
	context.JSON(http.StatusOK, gin.H{"rating": helper.AppConfig.RatingBoard.Ratings[index]})
}

// GetUserRatingByEmail godoc
// @Summary Endpoint to get the rating related to the user
// @Schemes
// @Description Get the user classification rating based on it's age
// @Tags rating
// @Accept json
// @Produce json
// @Param email path string true "User Email"
// @Success 200 {object} helper.Config.RatingBoard.Ratings
// @Failure 204 {string} string "No Content"
// @Failure 400 {string} string "Bad request"
// @Router /secured/user/{email} [get]
func GetUserRatingByEmail(context *gin.Context) {
	user, err := middlewares.GetUserByEmail(context.Param("email"))
	if err != nil {
		log.Errorf("GetUserRatingByEmail: Bad Request [%v]", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": err})
		context.Abort()
		return
	}

	index := getRatingProperty(user.Age)
	if index >= len(helper.AppConfig.RatingBoard.Ratings) {
		log.Errorf("GetUserRatingByEmail: No Content [%v]", user.Age)
		context.JSON(http.StatusNoContent, gin.H{"error": err})
		context.Abort()
		return
	}

	log.Infof("GetUserRatingByEmail: [%v]", helper.AppConfig.RatingBoard.Ratings[index])
	context.JSON(http.StatusOK, gin.H{"rating": helper.AppConfig.RatingBoard.Ratings[index]})
}
