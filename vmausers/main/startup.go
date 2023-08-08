package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"sort"
	"vmausers/constants"
	"vmausers/controllers"
	"vmausers/database"
	"vmausers/helper"
	"vmausers/middlewares"

	docs "vmausers/docs"

	"github.com/gin-gonic/gin"
	"github.com/julienschmidt/httprouter"
	cors "github.com/rs/cors/wrapper/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	httpSwagger "github.com/swaggo/http-swagger"

	log "github.com/sirupsen/logrus"
)

func StartBase(configFile string) {
	file := flag.String("config", configFile, "a config file")
	flag.Parse()

	_, err := os.Stat(*file)
	if err != nil {
		fmt.Printf("File not found: %v", err)
		os.Exit(1)
	}

	helper.AppConfig, err = helper.LoadConfig(*file)
	if err != nil || len(helper.AppConfig.Mongodb.Serveruri) == 0 {
		fmt.Printf("Error opening the config file: %v", err)
		os.Exit(1)
	}

	sort.Slice(helper.AppConfig.RatingBoard.Ratings, func(i, j int) bool {
		return helper.AppConfig.RatingBoard.Ratings[i].MinAge <= helper.AppConfig.RatingBoard.Ratings[j].MinAge
	})

	helper.InitLogger()
	fmt.Println("Base configuration OK")
}

func CheckDatabase() {
	client, err := database.NewConnection(&helper.AppConfig)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	client.Disconnect(context.Background())
	fmt.Println("Database OK")
}

func StartRouter() *gin.Engine {
	router := gin.Default()
	router.StaticFS("/static", http.Dir("./Static"))
	router.Use(gin.Recovery())
	router.Use(cors.AllowAll())

	// command <vmausers>swag init -g ./main/main.go -o ./docs
	docs.SwaggerInfo.BasePath = constants.Api_base_path + "/" + constants.Api_version
	api := router.Group(constants.Api_base_path + "/" + constants.Api_version)
	{
		api.GET("/", func(c *gin.Context) {
			c.JSON(
				http.StatusOK,
				gin.H{"message": "API working good"},
			)
		})

		api.POST(constants.Api_token_path, controllers.GenerateToken)
		api.POST(constants.Api_user_path+"/register", controllers.RegisterUser)
		secured := api.Group(constants.Api_secured_path).Use(middlewares.Authorization())
		{
			secured.GET("/", func(c *gin.Context) {
				c.JSON(
					http.StatusOK,
					gin.H{"message": "This token is valid"},
				)
			})

			secured.GET(constants.Api_user_path, controllers.GetManyUsers)
			secured.GET(constants.Api_user_path+"/:email", controllers.GetUserByEmail)

			secured.PUT(constants.Api_user_path, controllers.UpdateUser)
			secured.PUT(constants.Api_user_path+"/update/email", controllers.UpdateUserEmail)
			secured.PUT(constants.Api_user_path+"/update/password", controllers.UpdateUserPassword)

			secured.DELETE(constants.Api_user_path+"/:email", controllers.DeleteUserByEmail)

			secured.GET(constants.Api_rating_path, controllers.GetRatingList)
			secured.GET(constants.Api_rating_path+"/byage/:age", controllers.GetAgeRating)
			secured.GET(constants.Api_rating_path+"/user/byemail/:email", controllers.GetUserRatingByEmail)
		}

		api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}

	fmt.Println("Router OK")
	return router
}

func SwaggerRouter() error {
	router := httprouter.New()

	router.ServeFiles("/api/doc/static/*filepath", http.Dir("api/swagger/static"))
	router.HandlerFunc(http.MethodGet, "/api/doc/index.html", swaggerHandler)
	// router.HandlerFunc(http.MethodGet, "/api/doc", swaggerHandler)

	fmt.Println("Server on port 8080")
	return http.ListenAndServe(":8080", router)
}

func swaggerHandler(res http.ResponseWriter, req *http.Request) {
	httpSwagger.WrapHandler(res, req)
}
