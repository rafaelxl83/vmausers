package main

import (
	"context"
	"os"

	_ "vmausers/docs"
)

// @title VMA APIs
// @version 1.0
// @description VMA Swagger APIs.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email rafael.xavier.lima@gmail.com
// @securityDefinitions.apiKey JWT
// @in header
// @name token
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /api/v1
// @schemes http
func main() {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	client, err := StartDatabase(path + "\\config.json")
	if err != nil {
		panic(err)
	}

	client.Disconnect(context.Background())

	router := StartRouter()
	router.Run(":8080")

}
