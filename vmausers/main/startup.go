package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
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
	"go.mongodb.org/mongo-driver/mongo"
)

func StartDatabase(configFile string) (*mongo.Client, error) {
	file := flag.String("config", configFile, "a config file")
	flag.Parse()

	_, err := os.Stat(*file)
	if err != nil {
		fmt.Printf("File not found: %v", err)
		return nil, err
	}

	helper.AppConfig, err = helper.LoadConfig(*file)
	if err != nil || len(helper.AppConfig.Mongodb.Serveruri) == 0 {
		fmt.Printf("Error opening the config file: %v", err)
		return nil, err
	}

	client, err := database.NewConnection(&helper.AppConfig)
	if err != nil {
		fmt.Printf("Error connecting to the database: %v", err)
		return nil, err
	}

	return client, nil
}

func StartRouter() *gin.Engine {
	router := gin.Default()
	router.StaticFS("/static", http.Dir("./Static"))
	router.Use(gin.Recovery())
	router.Use(cors.AllowAll())

	// command <vmausers>swag init -g ./main/main.go -o ./docs
	docs.SwaggerInfo.BasePath = "/api/" + constants.Api_version
	api := router.Group("/api/" + constants.Api_version)
	{
		api.GET("/", func(c *gin.Context) {
			c.JSON(
				http.StatusOK,
				gin.H{"message": "API working good"},
			)
		})

		api.POST("/token", controllers.GenerateToken)
		api.POST("/user/register", controllers.RegisterUser)
		secured := api.Group("/secured").Use(middlewares.Authorization())
		{
			secured.GET("/", func(c *gin.Context) {
				c.JSON(
					http.StatusOK,
					gin.H{"message": "This token is valid"},
				)
			})

			secured.GET("/users", controllers.GetManyUsers)
			secured.GET("/user/:email", controllers.GetUserByEmail)

			secured.PUT("/user", controllers.UpdateUser)
			secured.PUT("/user/update/email", controllers.UpdateUserEmail)
			secured.PUT("/user/update/password", controllers.UpdateUserPassword)

			secured.DELETE("/user/:email", controllers.DeleteUserByEmail)
		}

		api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}

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
